package cfdi

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
)

// Impuestos store the invoice tax totals
type Impuestos struct {
	TotalImpuestosTrasladados string     `xml:",attr,omitempty"`
	Traslados                 *Traslados `xml:"cfdi:Traslados"`
}

// Traslados list the applicable taxes of the invoice or a line
type Traslados struct {
	Traslado []*Traslado `xml:"cfdi:Traslado"`
}

// Traslado stores the tax data of the invoice or a line
type Traslado struct {
	Base       string `xml:",attr"`
	Importe    string `xml:",attr"`
	Impuesto   string `xml:",attr"`
	TasaOCuota string `xml:",attr"`
	TipoFactor string `xml:",attr"`
}

func newImpuestos(totals *bill.Totals) *Impuestos {
	impuestos := &Impuestos{
		TotalImpuestosTrasladados: totals.Tax.String(),
		Traslados:                 newTraslados(totals.Taxes),
	}

	return impuestos
}

func newImpuestosFromLine(line *bill.Line) *Impuestos {
	impuestos := &Impuestos{
		Traslados: newTrasladosFromLine(line),
	}

	return impuestos
}

func newTraslados(taxTotal *tax.Total) *Traslados {
	var traslados []*Traslado

	for _, cat := range taxTotal.Categories {
		if cat.Code != common.TaxCategoryVAT {
			continue
		}

		for _, rate := range cat.Rates {
			traslados = append(traslados, newTraslado(rate))
		}
	}

	return &Traslados{traslados}
}

func newTrasladosFromLine(line *bill.Line) *Traslados {
	var traslados []*Traslado

	for _, tax := range line.Taxes {
		if tax.Category != common.TaxCategoryVAT {
			continue
		}

		traslados = append(traslados, newTrasladoFromLineTax(line, tax))
	}

	return &Traslados{traslados}
}

func newTraslado(rate *tax.RateTotal) *Traslado {
	traslado := &Traslado{
		Base:       rate.Base.String(),
		Importe:    rate.Amount.String(),
		Impuesto:   ImpuestoIVA,
		TasaOCuota: formatTaxPercent(rate.Percent),
		TipoFactor: TipoFactorTasa,
	}

	return traslado
}

func newTrasladoFromLineTax(line *bill.Line, tax *tax.Combo) *Traslado {
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
