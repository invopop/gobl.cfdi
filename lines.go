package cfdi

import (
	"github.com/invopop/gobl.cfdi/internal"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/tax"
)

// Conceptos list invoice lines
// nolint:misspell
type Conceptos struct {
	Concepto []*Concepto `xml:"cfdi:Concepto"`
}

// Concepto stores an invoice line data
type Concepto struct {
	ClaveProdServ string `xml:",attr"`
	Cantidad      string `xml:",attr"`
	ClaveUnidad   string `xml:",attr"`
	Descripcion   string `xml:",attr"` // nolint:misspell
	ValorUnitario string `xml:",attr"`
	Importe       string `xml:",attr"`
	Descuento     string `xml:",attr,omitempty"`
	ObjetoImp     string `xml:",attr"`

	Impuestos *ConceptoImpuestos `xml:"cfdi:Impuestos,omitempty"`
}

// nolint:misspell
func newConceptos(lines []*bill.Line, regime *tax.RegimeDef) *Conceptos {
	var conceptos []*Concepto

	for _, line := range lines {
		conceptos = append(conceptos, newConcepto(line, regime))
	}

	return &Conceptos{conceptos}
}

func newConcepto(line *bill.Line, regime *tax.RegimeDef) *Concepto {
	concepto := &Concepto{
		ClaveProdServ: internal.ClaveProdServ(line).String(),
		Cantidad:      line.Quantity.String(),
		ClaveUnidad:   internal.ClaveUnidad(line).String(),
		Descripcion:   line.Item.Name, // nolint:misspell
		ValorUnitario: line.Item.Price.String(),
		Importe:       line.Sum.String(),
		Descuento:     formatOptionalAmount(internal.TotalLineDiscount(line)),
		ObjetoImp:     lineSubjectToTax(line),
		Impuestos:     newConceptoImpuestos(line, regime),
	}

	return concepto
}

func lineSubjectToTax(line *bill.Line) string {
	if len(line.Taxes) == 0 {
		return ObjetoImpNo
	}
	return ObjetoImpSi
}
