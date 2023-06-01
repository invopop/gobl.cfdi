package cfdi

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

type Impuestos struct {
	TotalImpuestosTrasladados string     `xml:",attr,omitempty"`
	Traslados                 *Traslados `xml:"cfdi:Traslados"`
}

type Traslados struct {
	Traslado []*Traslado `xml:"cfdi:Traslado"`
}

type Traslado struct {
	Base       string `xml:",attr"`
	Importe    string `xml:",attr"`
	Impuesto   string `xml:",attr"`
	TasaOCuota string `xml:",attr"`
	TipoFactor string `xml:",attr"`
}

func NewImpuestos(totals *bill.Totals) *Impuestos {
	impuestos := &Impuestos{
		TotalImpuestosTrasladados: totals.Tax.String(),
		Traslados:                 NewTraslados(totals.Taxes),
	}

	return impuestos
}

func NewImpuestosFromLine(line *bill.Line) *Impuestos {
	impuestos := &Impuestos{
		Traslados: NewTrasladosFromLine(line),
	}

	return impuestos
}

func NewTraslados(taxTotal *tax.Total) *Traslados {
	var traslados []*Traslado

	for _, cat := range taxTotal.Categories {
		if cat.Code != common.TaxCategoryVAT {
			continue
		}

		for _, rate := range cat.Rates {
			traslados = append(traslados, NewTraslado(rate))
		}
	}

	return &Traslados{traslados}
}

func NewTrasladosFromLine(line *bill.Line) *Traslados {
	var traslados []*Traslado

	for _, tax := range line.Taxes {
		if tax.Category != common.TaxCategoryVAT {
			continue
		}

		traslados = append(traslados, NewTrasladoFromLineTax(line, tax))
	}

	return &Traslados{traslados}
}

func NewTraslado(rate *tax.RateTotal) *Traslado {
	traslado := &Traslado{
		Base:       rate.Base.String(),
		Importe:    rate.Amount.String(),
		Impuesto:   ImpuestoIVA,
		TasaOCuota: formatTaxPercent(rate.Percent),
		TipoFactor: TipoFactorTasa,
	}

	return traslado
}

func NewTrasladoFromLineTax(line *bill.Line, tax *tax.Combo) *Traslado {
	// GOBL doesn't provide an amount at line level, so we calculate it
	taxAmount := tax.Percent.Of(line.Total)

	traslado := &Traslado{
		Base:       line.Total.String(),
		Importe:    taxAmount.String(),
		Impuesto:   ImpuestoIVA,
		TasaOCuota: formatTaxPercent(tax.Percent),
		TipoFactor: TipoFactorTasa,
	}

	return traslado
}

func formatTaxPercent(percent *num.Percentage) string {
	return percent.Amount.Rescale(6).String()
}
