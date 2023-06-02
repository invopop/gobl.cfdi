// Package cfdi implements the conversion from GOBL to CFDI XML
package cfdi

import (
	"encoding/xml"
	"fmt"

	"cloud.google.com/go/civil"
	"github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
)

// CFDI schema constants
const (
	CFDINamespace  = "http://www.sat.gob.mx/cfd/4"
	XSINamespace   = "http://www.w3.org/2001/XMLSchema-instance"
	SchemaLocation = "http://www.sat.gob.mx/cfd/4 http://www.sat.gob.mx/sitio_internet/cfd/4/cfdv40.xsd"
	CFDIVersion    = "4.0"
)

// Hard-coded values for (yet) unsupported mappings
const (
	FakeNoCertificado             = "00000000000000000000"
	TipoDeComprobanteIngreso      = "I"
	ExportacionNoAplica           = "01"
	MetodoPagoUnaExhibicion       = "PUE"
	FormaPagoPorDefinir           = "99"
	ClaveProdServNoExiste         = "01010101"
	ClaveUnidadMutuamenteDefinida = "H87"
	ObjetoImpSi                   = "02"
	ImpuestoIVA                   = "002"
	TipoFactorTasa                = "Tasa"
	UsoCFDIGastosGenerales        = "G03"
	RegimenFiscalGeneral          = "601"
)

// Document is a pseudo-model for containing the XML document being created
type Document struct {
	XMLName        xml.Name `xml:"cfdi:Comprobante"`
	CFDINamespace  string   `xml:"xmlns:cfdi,attr"`
	XSINamespace   string   `xml:"xmlns:xsi,attr"`
	SchemaLocation string   `xml:"xsi:schemaLocation,attr"`
	Version        string   `xml:"Version,attr"`

	TipoDeComprobante string `xml:",attr"`
	Serie             string `xml:",attr"`
	Folio             string `xml:",attr"`
	Fecha             string `xml:",attr"`
	LugarExpedicion   string `xml:",attr"`
	SubTotal          string `xml:",attr"`
	Total             string `xml:",attr"`
	Moneda            string `xml:",attr"`
	Exportacion       string `xml:",attr"`
	MetodoPago        string `xml:",attr"`
	FormaPago         string `xml:",attr"`
	Sello             string `xml:",attr"`
	NoCertificado     string `xml:",attr"`
	Certificado       string `xml:",attr"`

	Emisor    *Emisor    `xml:"cfdi:Emisor"`
	Receptor  *Receptor  `xml:"cfdi:Receptor"`
	Conceptos *Conceptos `xml:"cfdi:Conceptos"` //nolint:misspell
	Impuestos *Impuestos `xml:"cfdi:Impuestos"`
}

// NewDocument converts a GOBL envelope into a CFDI document
func NewDocument(env *gobl.Envelope) (*Document, error) {
	inv, ok := env.Extract().(*bill.Invoice)
	if !ok {
		return nil, fmt.Errorf("invalid type %T", env.Document)
	}

	document := &Document{
		CFDINamespace:  CFDINamespace,
		XSINamespace:   XSINamespace,
		SchemaLocation: SchemaLocation,
		Version:        CFDIVersion,

		TipoDeComprobante: TipoDeComprobanteIngreso,
		Serie:             inv.Series,
		Folio:             inv.Code,
		Fecha:             formatIssueDate(inv.IssueDate),
		LugarExpedicion:   inv.Supplier.Addresses[0].Code,
		SubTotal:          inv.Totals.Total.String(),
		Total:             inv.Totals.TotalWithTax.String(),
		Moneda:            string(inv.Currency),
		Exportacion:       ExportacionNoAplica,
		MetodoPago:        MetodoPagoUnaExhibicion,
		FormaPago:         FormaPagoPorDefinir,

		NoCertificado: FakeNoCertificado,

		Emisor:    newEmisor(inv.Supplier),
		Receptor:  newReceptor(inv.Customer),
		Conceptos: newConceptos(inv.Lines), // nolint:misspell
		Impuestos: newImpuestos(inv.Totals),
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

func formatIssueDate(date cal.Date) string {
	dateTime := civil.DateTime{Date: date.Date, Time: civil.Time{}}
	return dateTime.String()
}
