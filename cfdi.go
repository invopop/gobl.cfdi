// Package cfdi implements the conversion from GOBL to CFDI XML
package cfdi

import (
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl.cfdi/addendas"
	"github.com/invopop/gobl.cfdi/internal"
	"github.com/invopop/gobl.cfdi/internal/format"
	"github.com/invopop/gobl/addons/mx/cfdi"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
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
	FakeNoCertificado   = "00000000000000000000"
	ExportacionNoAplica = "01"
	FormaPagoPorDefinir = "99"
	ImpuestoIVA         = "002"
)

// Generic supplier constants
const (
	NombreReceptorGenerico       = "PÚBLICO EN GENERAL"
	RegimenFiscalSinObligaciones = "616" // no tax obligations
	UsoCFDISinEfectos            = "S01" // no tax effects
)

// TipoFactor definitions.
const (
	TipoFactorTasa   = "Tasa"
	TipoFactorCuota  = "Cuota" // Not supported
	TipoFactorExento = "Exento"
)

// Subject to tax constants
const (
	ObjetoImpNo = "01" // not subject to tax
	ObjetoImpSi = "02" // subject to tax
)

// ErrNotSupported is returned when the conversion of the invoice is not supported
var ErrNotSupported = errors.New("not supported")

// regime global
var regime = tax.RegimeDefFor("MX")

// Document is a pseudo-model for containing the XML document being created
type Document struct {
	XMLName        xml.Name `xml:"cfdi:Comprobante"`
	CFDINamespace  string   `xml:"xmlns:cfdi,attr"`
	XSINamespace   string   `xml:"xmlns:xsi,attr"`
	ECCNamespace   string   `xml:"xmlns:ecc12,attr,omitempty"`
	VDNamespace    string   `xml:"xmlns:valesdedespensa,attr,omitempty"`
	SchemaLocation string   `xml:"xsi:schemaLocation,attr"`
	Version        string   `xml:"Version,attr"`

	TipoDeComprobante string      `xml:",attr"`
	Serie             string      `xml:",attr,omitempty"`
	Folio             string      `xml:",attr,omitempty"`
	Fecha             string      `xml:",attr"`
	LugarExpedicion   string      `xml:",attr"`
	SubTotal          num.Amount  `xml:",attr"`
	Descuento         *num.Amount `xml:",attr,omitempty"`
	Total             num.Amount  `xml:",attr"`
	Moneda            string      `xml:",attr"`
	TipoCambio        *num.Amount `xml:",attr,omitempty"`
	Exportacion       string      `xml:",attr"`
	MetodoPago        string      `xml:",attr,omitempty"`
	FormaPago         string      `xml:",attr,omitempty"`
	CondicionesDePago string      `xml:",attr,omitempty"`
	Sello             string      `xml:",attr"`
	NoCertificado     string      `xml:",attr"`
	Certificado       string      `xml:",attr"`

	Global           *GlobalInformation `xml:"cfdi:InformacionGlobal,omitempty"`
	CFDIRelacionados *CFDIRelacionados  `xml:"cfdi:CfdiRelacionados,omitempty"`
	Emisor           *Emisor            `xml:"cfdi:Emisor"`
	Receptor         *Receptor          `xml:"cfdi:Receptor"`
	Conceptos        *Conceptos         `xml:"cfdi:Conceptos"` //nolint:misspell
	Impuestos        *Impuestos         `xml:"cfdi:Impuestos,omitempty"`

	Complemento *internal.Nodes `xml:"cfdi:Complemento,omitempty"`
	Addenda     *internal.Nodes `xml:"cfdi:Addenda,omitempty"`
}

// GlobalInformation is used for invoices that contain a summary of B2C documents.
type GlobalInformation struct {
	Period string `xml:"Periodicidad,attr"`
	Month  string `xml:"Meses,attr"`
	Year   string `xml:"Año,attr"`
}

// NewDocument converts a GOBL envelope into a CFDI document
func NewDocument(env *gobl.Envelope) (*Document, error) {
	inv, ok := env.Extract().(*bill.Invoice)
	if !ok {
		return nil, fmt.Errorf("invalid type %T", env.Document)
	}

	if err := validateSupport(inv); err != nil {
		return nil, err
	}

	issuePlace := issuePlace(inv)

	doc := &Document{
		CFDINamespace:  CFDINamespace,
		XSINamespace:   XSINamespace,
		SchemaLocation: format.SchemaLocation(CFDINamespace, CFDISchemaLocation),
		Version:        CFDIVersion,

		TipoDeComprobante: lookupTipoDeComprobante(inv),
		Serie:             inv.Series.String(),
		Folio:             inv.Code.String(),
		Fecha:             formatIssueDateTime(inv),
		LugarExpedicion:   issuePlace,
		Descuento:         internal.TotalInvoiceDiscount(inv),
		Total:             inv.Totals.Payable,
		Moneda:            string(inv.Currency),
		TipoCambio:        tipoCambio(inv),
		Exportacion:       ExportacionNoAplica,
		MetodoPago:        metodoPago(inv),
		FormaPago:         formaPago(inv),
		CondicionesDePago: paymentTermsNotes(inv),

		NoCertificado: FakeNoCertificado,

		Global:           newGlobalInformation(inv),
		CFDIRelacionados: newCfdiRelacionados(inv),
		Emisor:           newEmisor(inv.Supplier),
		Receptor:         newReceptor(inv.Customer, issuePlace),
		Conceptos:        newConceptos(inv.Lines), // nolint:misspell
		Impuestos:        newImpuestos(inv.Totals, inv.Lines, inv.Currency),
	}

	// Determine the subtotal directly from the concepts, as there may be some
	// additional taxes that needed to be taken into account
	doc.SubTotal = inv.Currency.Def().Zero()
	for _, c := range doc.Conceptos.Concepto {
		doc.SubTotal = doc.SubTotal.Add(c.Importe)
	}

	// Recalculate the total so that we can avoid any rounding issues
	doc.Total = doc.SubTotal
	if doc.Descuento != nil {
		doc.Total = doc.Total.Subtract(*doc.Descuento)
	}
	if doc.Impuestos != nil {
		if tit := doc.Impuestos.TotalImpuestosTrasladados; tit != nil {
			doc.Total = doc.Total.Add(*tit)
		}
		if tir := doc.Impuestos.TotalImpuestosRetenidos; tir != nil {
			doc.Total = doc.Total.Subtract(*tir)
		}
	}

	if err := addComplementos(doc, inv.Complements); err != nil {
		return nil, err
	}

	if err := addAddendas(doc, inv); err != nil {
		return nil, err
	}

	return doc, nil
}

func newGlobalInformation(inv *bill.Invoice) *GlobalInformation {
	if inv.Tax == nil || !inv.Tax.Ext.Has(cfdi.ExtKeyGlobalPeriod) {
		return nil
	}
	return &GlobalInformation{
		Period: inv.Tax.Ext[cfdi.ExtKeyGlobalPeriod].String(),
		Month:  inv.Tax.Ext[cfdi.ExtKeyGlobalMonth].String(),
		Year:   inv.Tax.Ext[cfdi.ExtKeyGlobalYear].String(),
	}
}

func validateSupport(inv *bill.Invoice) error {
	errs := validation.Errors{}

	if len(inv.Charges) > 0 {
		errs["charges"] = ErrNotSupported
	}

	// Deprecation pending...
	if inv.HasTags(tax.TagSelfBilled) {
		errs["self-billed"] = ErrNotSupported
	}
	if inv.HasTags(tax.TagCustomerRates) {
		errs["customer-rates"] = ErrNotSupported
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func issuePlace(inv *bill.Invoice) string {
	if inv.Tax != nil && inv.Tax.Ext.Has(cfdi.ExtKeyIssuePlace) {
		return inv.Tax.Ext[cfdi.ExtKeyIssuePlace].String()
	}
	// Fallback
	return inv.Supplier.Ext[cfdi.ExtKeyIssuePlace].String()
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
		case *cfdi.FuelAccountBalance:
			addEstadoCuentaCombustible(doc, o)
		case *cfdi.FoodVouchers:
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

func formatIssueDateTime(inv *bill.Invoice) string {
	tn := inv.IssueTime
	if tn == nil {
		tn = new(cal.Time) // zero
	}
	return inv.IssueDate.WithTime(*tn).String()
}

func lookupTipoDeComprobante(inv *bill.Invoice) string {
	if inv.Tax == nil {
		return ""
	}

	return inv.Tax.Ext[cfdi.ExtKeyDocType].String()
}

func tipoCambio(inv *bill.Invoice) *num.Amount {
	r := currency.MatchExchangeRate(inv.ExchangeRates, inv.Currency, currency.MXN)
	if r == nil {
		return nil
	}
	a := r.Amount
	return &a
}

func metodoPago(inv *bill.Invoice) string {
	if inv.Tax != nil && inv.Tax.Ext.Has(cfdi.ExtKeyPaymentMethod) {
		return inv.Tax.Ext[cfdi.ExtKeyPaymentMethod].String()
	}
	// Fallback to the payment method based on the detected payment advances
	if isPrepaid(inv) {
		return cfdi.ExtCodePaymentMethodPUE.String()
	}
	return cfdi.ExtCodePaymentMethodPPD.String()
}

func formaPago(inv *bill.Invoice) string {
	adv := largestAdvance(inv)
	if !isPrepaid(inv) || adv == nil {
		return FormaPagoPorDefinir
	}
	return adv.Ext[cfdi.ExtKeyPaymentMeans].String()
}

func isPrepaid(inv *bill.Invoice) bool {
	return inv.Totals.Due != nil && inv.Totals.Due.IsZero()
}

func largestAdvance(inv *bill.Invoice) *pay.Advance {
	if inv.Payment == nil || len(inv.Payment.Advances) == 0 {
		return nil
	}

	la := inv.Payment.Advances[0]
	for _, a := range inv.Payment.Advances {
		if a.Amount.Compare(la.Amount) == 1 {
			la = a
		}
	}
	return la
}

func paymentTermsNotes(inv *bill.Invoice) string {
	if inv.Payment == nil || inv.Payment.Terms == nil {
		return ""
	}

	return inv.Payment.Terms.Notes
}
