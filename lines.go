package cfdi

import (
	"github.com/invopop/gobl.cfdi/internal"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
)

// Conceptos list invoice lines
// nolint:misspell
type Conceptos struct {
	Concepto []*Concepto `xml:"cfdi:Concepto"`
}

// Concepto stores an invoice line data
type Concepto struct {
	ClaveProdServ string      `xml:",attr"`
	Ref           string      `xml:"NoIdentificacion,attr,omitempty"`
	Cantidad      string      `xml:",attr"`
	ClaveUnidad   string      `xml:",attr"`
	Desc          string      `xml:"Descripcion,attr"` // nolint:misspell
	ValorUnitario num.Amount  `xml:",attr"`
	Importe       num.Amount  `xml:",attr"`
	Descuento     *num.Amount `xml:",attr,omitempty"`
	ObjetoImp     string      `xml:",attr"`

	Impuestos *ConceptoImpuestos `xml:"cfdi:Impuestos,omitempty"`
}

// nolint:misspell
func newConceptos(lines []*bill.Line) *Conceptos {
	var conceptos []*Concepto

	for _, line := range lines {
		if line.Sum == nil {
			continue
		}
		conceptos = append(conceptos, newConcepto(line))
	}

	return &Conceptos{conceptos}
}

func newConcepto(line *bill.Line) *Concepto {
	concepto := &Concepto{
		ClaveProdServ: internal.ClaveProdServ(line).String(),
		Ref:           line.Item.Ref.String(),
		Cantidad:      line.Quantity.String(),
		ClaveUnidad:   internal.ClaveUnidad(line).String(),
		Desc:          line.Item.Name,
		ValorUnitario: *line.Item.Price,
		Importe:       *line.Sum,
		Descuento:     internal.TotalLineDiscount(line),
		ObjetoImp:     lineSubjectToTax(line),
		Impuestos:     newConceptoImpuestos(line),
	}

	return concepto
}

func lineSubjectToTax(line *bill.Line) string {
	if len(line.Taxes) == 0 {
		return ObjetoImpNo
	}
	return ObjetoImpSi
}
