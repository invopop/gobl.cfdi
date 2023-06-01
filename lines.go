package cfdi

import "github.com/invopop/gobl/bill"

type Conceptos struct {
	Concepto []*Concepto `xml:"cfdi:Concepto"`
}

type Concepto struct {
	ClaveProdServ string `xml:",attr"`
	Cantidad      string `xml:",attr"`
	ClaveUnidad   string `xml:",attr"`
	Descripcion   string `xml:",attr"`
	ValorUnitario string `xml:",attr"`
	Importe       string `xml:",attr"`
	ObjetoImp     string `xml:",attr"`

	Impuestos *Impuestos `xml:"cfdi:Impuestos"`
}

func NewConceptos(lines []*bill.Line) *Conceptos {
	var conceptos []*Concepto

	for _, line := range lines {
		conceptos = append(conceptos, NewConcepto(line))
	}

	return &Conceptos{conceptos}
}

func NewConcepto(line *bill.Line) *Concepto {
	concepto := &Concepto{
		ClaveProdServ: ClaveProdServNoExiste,
		Cantidad:      line.Quantity.String(),
		ClaveUnidad:   ClaveUnidadMutuamenteDefinida,
		Descripcion:   line.Item.Name,
		ValorUnitario: line.Item.Price.String(),
		Importe:       line.Total.String(),
		ObjetoImp:     ObjetoImpSi,
		Impuestos:     NewImpuestosFromLine(line),
	}

	return concepto
}
