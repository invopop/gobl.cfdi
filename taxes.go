package cfdi

import (
	"github.com/invopop/gobl.cfdi/internal/format"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/mx"
	"github.com/invopop/gobl/tax"
)

// Impuestos store the invoice tax totals
type Impuestos struct {
	TotalImpuestosTrasladados string       `xml:",attr,omitempty"`
	TotalImpuestosRetenidos   string       `xml:",attr,omitempty"`
	Retenciones               *Retenciones `xml:"cfdi:Retenciones,omitempty"`
	Traslados                 *Traslados   `xml:"cfdi:Traslados,omitempty"`
}

// ConceptoImpuestos store the line tax totals
type ConceptoImpuestos struct {
	Traslados   *Traslados   `xml:"cfdi:Traslados,omitempty"`
	Retenciones *Retenciones `xml:"cfdi:Retenciones,omitempty"`
}

// Traslados lists the non-retained taxes of a line or the invoice
type Traslados struct {
	Traslado []*Impuesto `xml:"cfdi:Traslado"`
}

// Retenciones lists the retained taxes of a line or the invoice
type Retenciones struct {
	Retencion []*Impuesto `xml:"cfdi:Retencion"`
}

// Impuesto stores the tax data of the invoice or a line
type Impuesto struct {
	Base       string `xml:",attr,omitempty"`
	Importe    string `xml:",attr,omitempty"`
	Impuesto   string `xml:",attr"`
	TasaOCuota string `xml:",attr,omitempty"`
	TipoFactor string `xml:",attr,omitempty"`
}

func newImpuestos(totals *bill.Totals, currency *currency.Code, regime *tax.Regime) *Impuestos {
	var traslados, retenciones []*Impuesto
	totalTraslados, totalRetenciones := currency.Def().Zero(), currency.Def().Zero()

	for _, cat := range totals.Taxes.Categories {
		catDef := regime.Category(cat.Code)

		for _, rate := range cat.Rates {
			if rate.Percent == nil {
				// skip exempt taxes
				continue
			}
			imp := newImpuesto(rate, currency, catDef)
			if catDef.Retained {
				// Clear out fields not supported by retained totals
				imp.Base = ""
				imp.TasaOCuota = ""
				imp.TipoFactor = ""

				retenciones = append(retenciones, imp)
				totalRetenciones = totalRetenciones.Add(rate.Amount)
			} else {
				traslados = append(traslados, imp)
				totalTraslados = totalTraslados.Add(rate.Amount)
			}
		}
	}

	impuestos := &Impuestos{}
	empty := true
	if len(traslados) > 0 {
		impuestos.Traslados = &Traslados{traslados}
		impuestos.TotalImpuestosTrasladados = totalTraslados.String()
		empty = false
	}
	if len(retenciones) > 0 {
		impuestos.Retenciones = &Retenciones{retenciones}
		impuestos.TotalImpuestosRetenidos = totalRetenciones.String()
		empty = false
	}
	if empty {
		return nil
	}

	return impuestos
}

func newImpuesto(rate *tax.RateTotal, currency *currency.Code, catDef *tax.Category) *Impuesto {
	cu := currency.Def().Subunits // SAT expects tax total amounts with no more decimals than supported by the currency

	ratePercent := rate.Percent
	imp := &Impuesto{
		Base:       rate.Base.Rescale(cu).String(),
		Importe:    rate.Amount.Rescale(cu).String(),
		Impuesto:   catDef.Map[mx.KeySATImpuesto].String(),
		TasaOCuota: format.TaxPercent(ratePercent),
		TipoFactor: TipoFactorTasa,
	}

	return imp
}

func newConceptoImpuestos(line *bill.Line, regime *tax.Regime) *ConceptoImpuestos {
	var traslados, retenciones []*Impuesto

	for _, tax := range line.Taxes {
		if tax.Percent == nil {
			// skip exempt taxes
			continue
		}
		catDef := regime.Category(tax.Category)
		imp := newConceptoImpuesto(line, tax, catDef)
		if catDef.Retained {
			retenciones = append(retenciones, imp)
		} else {
			traslados = append(traslados, imp)
		}
	}

	impuestos := &ConceptoImpuestos{}
	if len(traslados) > 0 {
		impuestos.Traslados = &Traslados{traslados}
	}
	if len(retenciones) > 0 {
		impuestos.Retenciones = &Retenciones{retenciones}
	}
	if impuestos.Traslados == nil && impuestos.Retenciones == nil {
		return nil
	}

	return impuestos
}

func newConceptoImpuesto(line *bill.Line, tax *tax.Combo, catDef *tax.Category) *Impuesto {
	// GOBL doesn't provide an amount at line level, so we calculate it
	taxPercent := tax.Percent
	if taxPercent == nil {
		taxPercent = num.NewPercentage(0, 3)
	}
	taxAmount := taxPercent.Of(line.Total)

	i := &Impuesto{
		Base:       line.Total.String(),
		Importe:    taxAmount.String(),
		Impuesto:   catDef.Map[mx.KeySATImpuesto].String(),
		TasaOCuota: format.TaxPercent(taxPercent),
		TipoFactor: TipoFactorTasa,
	}

	return i
}
