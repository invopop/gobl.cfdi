// Package cfdi implements the conversion from GOBL to CFDI XML
package cfdi

import (
	"encoding/xml"
	"fmt"

	"cloud.google.com/go/civil"
	"github.com/invopop/gobl"
	"github.com/invopop/gobl.cfdi/addendas"
	"github.com/invopop/gobl.cfdi/internal"
	"github.com/invopop/gobl.cfdi/internal/format"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regimes/mx"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/tax"
)

// CFDI schema constants
const (
	CFDINamespace      = "http://www.sat.gob.mx/cfd/4"
	CFDISchemaLocation = "http://www.sat.gob.mx/sitio_internet/cfd/4/cfdv40.xsd"
	XSINamespace       = "http://www.w3.org/2001/XMLSchema-instance"
	CFDIVersion        = "4.0"
)

// Hard-coded values for (yet) unsupported mappings
const (
	FakeNoCertificado       = "00000000000000000000"
	ExportacionNoAplica     = "01"
	MetodoPagoUnaExhibicion = "PUE"
	MetodoPagoParcialidades = "PPD"
	FormaPagoPorDefinir     = "99"
	ObjetoImpSi             = "02"
	ImpuestoIVA             = "002"
	TipoFactorTasa          = "Tasa"
)

// Document is a pseudo-model for containing the XML document being created
type Document struct {
	XMLName        xml.Name `xml:"cfdi:Comprobante"`
	CFDINamespace  string   `xml:"xmlns:cfdi,attr"`
	XSINamespace   string   `xml:"xmlns:xsi,attr"`
	ECCNamespace   string   `xml:"xmlns:ecc12,attr,omitempty"`
	VDNamespace    string   `xml:"xmlns:valesdedespensa,attr,omitempty"`
	SchemaLocation string   `xml:"xsi:schemaLocation,attr"`
	Version        string   `xml:"Version,attr"`

	TipoDeComprobante string `xml:",attr"`
	Serie             string `xml:",attr,omitempty"`
	Folio             string `xml:",attr,omitempty"`
	Fecha             string `xml:",attr"`
	LugarExpedicion   string `xml:",attr"`
	SubTotal          string `xml:",attr"`
	Descuento         string `xml:",attr,omitempty"`
	Total             string `xml:",attr"`
	Moneda            string `xml:",attr"`
	Exportacion       string `xml:",attr"`
	MetodoPago        string `xml:",attr,omitempty"`
	FormaPago         string `xml:",attr,omitempty"`
	CondicionesDePago string `xml:",attr,omitempty"`
	Sello             string `xml:",attr"`
	NoCertificado     string `xml:",attr"`
	Certificado       string `xml:",attr"`

	CFDIRelacionados *CFDIRelacionados `xml:"cfdi:CfdiRelacionados,omitempty"`
	Emisor           *Emisor           `xml:"cfdi:Emisor"`
	Receptor         *Receptor         `xml:"cfdi:Receptor"`
	Conceptos        *Conceptos        `xml:"cfdi:Conceptos"` //nolint:misspell
	Impuestos        *Impuestos        `xml:"cfdi:Impuestos,omitempty"`

	Complemento *internal.Nodes `xml:"cfdi:Complemento,omitempty"`
	Addenda     *internal.Nodes `xml:"cfdi:Addenda,omitempty"`
}

// NewDocument converts a GOBL envelope into a CFDI document
func NewDocument(env *gobl.Envelope) (*Document, error) {
	inv, ok := env.Extract().(*bill.Invoice)
	if !ok {
		return nil, fmt.Errorf("invalid type %T", env.Document)
	}

	discount := internal.TotalInvoiceDiscount(inv)
	subtotal := inv.Totals.Total.Add(discount)

	document := &Document{
		CFDINamespace:  CFDINamespace,
		XSINamespace:   XSINamespace,
		SchemaLocation: format.SchemaLocation(CFDINamespace, CFDISchemaLocation),
		Version:        CFDIVersion,

		TipoDeComprobante: lookupTipoDeComprobante(inv),
		Serie:             inv.Series,
		Folio:             inv.Code,
		Fecha:             formatIssueDate(inv.IssueDate),
		LugarExpedicion:   inv.Supplier.Ext[mx.ExtKeyCFDIPostCode].String(),
		SubTotal:          subtotal.String(),
		Descuento:         formatOptionalAmount(discount),
		Total:             inv.Totals.TotalWithTax.String(),
		Moneda:            string(inv.Currency),
		Exportacion:       ExportacionNoAplica,
		MetodoPago:        metodoPago(inv),
		FormaPago:         formaPago(inv),
		CondicionesDePago: paymentTermsNotes(inv),

		NoCertificado: FakeNoCertificado,

		CFDIRelacionados: newCfdiRelacionados(inv),
		Emisor:           newEmisor(inv.Supplier),
		Receptor:         newReceptor(inv.Customer),
		Conceptos:        newConceptos(inv.Lines, inv.TaxRegime()), // nolint:misspell
		Impuestos:        newImpuestos(inv.Totals, &inv.Currency, inv.TaxRegime()),
	}

	if err := addComplementos(document, inv.Complements); err != nil {
		return nil, err
	}

	if err := addAddendas(document, inv); err != nil {
		return nil, err
	}

	return document, nil
}

// Bytes returns the XML representation of the document in bytes
func (d *Document) Bytes() ([]byte, error) {
	bytes, err := xml.MarshalIndent(d, "", "  ")
	if err != nil {
		return nil, err
	}

	return append([]byte(xml.Header), bytes...), nil
}

// AppendComplemento appends a complement to the document
func (d *Document) AppendComplemento(c interface{}) {
	// We keep it nil unless an element is added so that no empty node is marshalled to XML
	if d.Complemento == nil {
		d.Complemento = &internal.Nodes{}
	}

	d.Complemento.Nodes = append(d.Complemento.Nodes, c)
}

// AppendAddenda appends an addenda to the document
func (d *Document) AppendAddenda(c interface{}) {
	// We keep it nil unless an element is added so that no empty node is marshalled to XML
	if d.Addenda == nil {
		d.Addenda = &internal.Nodes{}
	}

	d.Addenda.Nodes = append(d.Addenda.Nodes, c)
}

func addComplementos(doc *Document, complements []*schema.Object) error {
	for _, c := range complements {
		switch o := c.Instance().(type) {
		case *mx.FuelAccountBalance:
			addEstadoCuentaCombustible(doc, o)
		case *mx.FoodVouchers:
			addValesDeDespensa(doc, o)
		default:
			return fmt.Errorf("unsupported complement %T", o)
		}
	}

	return nil
}

func addAddendas(doc *Document, inv *bill.Invoice) error {
	ads, err := addendas.For(inv)
	if err != nil {
		return err
	}

	for _, ad := range ads {
		doc.AppendAddenda(ad)
	}
	return nil
}

func formatIssueDate(date cal.Date) string {
	dateTime := civil.DateTime{Date: date.Date, Time: civil.Time{}}
	return dateTime.String()
}

func lookupTipoDeComprobante(inv *bill.Invoice) string {
	ss := inv.ScenarioSummary()
	if ss == nil {
		return ""
	}

	code := ss.Codes[mx.KeySATTipoDeComprobante]
	return code.String()
}

func metodoPago(inv *bill.Invoice) string {
	if isPrepaid(inv) {
		return MetodoPagoUnaExhibicion
	}

	return MetodoPagoParcialidades
}

func formaPago(inv *bill.Invoice) string {
	if !isPrepaid(inv) {
		return FormaPagoPorDefinir
	}

	r := inv.TaxRegime()
	if r == nil {
		return ""
	}

	keyDef := findKeyDef(r.PaymentMeansKeys, largestAdvance(inv).Key)
	if keyDef == nil {
		return ""
	}

	code := keyDef.Map[mx.KeySATFormaPago]
	return code.String()
}

func isPrepaid(inv *bill.Invoice) bool {
	return inv.Totals.Due != nil && inv.Totals.Due.IsZero()
}

func largestAdvance(inv *bill.Invoice) *pay.Advance {
	la := inv.Payment.Advances[0]

	for _, a := range inv.Payment.Advances {
		if a.Amount.Compare(la.Amount) == 1 {
			la = a
		}
	}

	return la
}

func findKeyDef(keyDefs []*tax.KeyDefinition, key cbc.Key) *tax.KeyDefinition {
	for _, keyDef := range keyDefs {
		if keyDef.Key == key {
			return keyDef
		}
	}

	return nil
}

func paymentTermsNotes(inv *bill.Invoice) string {
	if inv.Payment == nil || inv.Payment.Terms == nil {
		return ""
	}

	return inv.Payment.Terms.Notes
}

func formatOptionalAmount(a num.Amount) string {
	if a.IsZero() {
		return ""
	}

	return a.String()
}
