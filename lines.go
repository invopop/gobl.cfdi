package cfdi

import "github.com/invopop/gobl/bill"

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
	ObjetoImp     string `xml:",attr"`

	Impuestos *Impuestos `xml:"cfdi:Impuestos"`
}

// nolint:misspell
func newConceptos(lines []*bill.Line) *Conceptos {
	var conceptos []*Concepto

	for _, line := range lines {
		conceptos = append(conceptos, newConcepto(line))
	}

	return &Conceptos{conceptos}
}

func newConcepto(line *bill.Line) *Concepto {
	concepto := &Concepto{
		ClaveProdServ: ClaveProdServNoExiste,
		Cantidad:      line.Quantity.String(),
		ClaveUnidad:   ClaveUnidadMutuamenteDefinida,
		Descripcion:   line.Item.Name, // nolint:misspell
		ValorUnitario: line.Item.Price.String(),
		Importe:       line.Total.String(),
		ObjetoImp:     ObjetoImpSi,
		Impuestos:     newImpuestosFromLine(line),
	}

	return concepto
}
