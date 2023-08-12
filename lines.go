package cfdi

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/mx"
)

// Default keys
const (
	DefaultClaveUnidad = "ZZ" // Mutuamente definida
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

	Impuestos *Impuestos `xml:"cfdi:Impuestos,omitempty"`
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
		ClaveProdServ: mapToClaveProdServ(line),
		Cantidad:      line.Quantity.String(),
		ClaveUnidad:   mapToClaveUnidad(line),
		Descripcion:   line.Item.Name, // nolint:misspell
		ValorUnitario: line.Item.Price.String(),
		Importe:       line.Sum.String(),
		Descuento:     formatOptionalAmount(totalLineDiscount(line)),
		ObjetoImp:     ObjetoImpSi,
		Impuestos:     newImpuestosFromLine(line),
	}

	return concepto
}

func mapToClaveUnidad(line *bill.Line) string {
	if line.Item.Unit == "" {
		return DefaultClaveUnidad
	}

	return string(line.Item.Unit.UNECE())
}

func mapToClaveProdServ(line *bill.Line) string {
	if line.Item == nil {
		return ""
	}

	id := org.IdentityForKey(line.Item.Identities, mx.IdentityKeyProdServ)
	if id != nil {
		return string(id.Code)
	}

	return ""
}

func totalLineDiscount(l *bill.Line) num.Amount {
	td := num.MakeAmount(0, l.Sum.Exp()) // discount's precision must match the "Importe" field's one
	for _, d := range l.Discounts {
		td = td.Add(d.Amount)
	}
	return td
}
