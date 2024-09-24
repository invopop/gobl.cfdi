package cfdi

import (
	"github.com/invopop/gobl.cfdi/internal/format"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
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

// Map of tax categories to SAT tax types
var taxCategoryMap = map[cbc.Code]string{
	mx.TaxCategoryISR:   "001",
	tax.CategoryVAT:     "002",
	mx.TaxCategoryRVAT:  "002",
	mx.TaxCategoryIEPS:  "003",
	mx.TaxCategoryRIEPS: "003",
}

func newImpuestos(totals *bill.Totals, currency *currency.Code) *Impuestos {
	var traslados, retenciones []*Impuesto
	totalTraslados, totalRetenciones := currency.Def().Zero(), currency.Def().Zero()

	for _, cat := range totals.Taxes.Categories {
		for _, rate := range cat.Rates {
			imp := newImpuesto(cat.Code, rate, currency)
			if cat.Retained {
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

		// Set tax total only if there are non-exempt taxes
		for _, t := range traslados {
			if t.TipoFactor != TipoFactorExento {
				impuestos.TotalImpuestosTrasladados = totalTraslados.String()
				break
			}
		}

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

func newImpuesto(catCode cbc.Code, rate *tax.RateTotal, currency *currency.Code) *Impuesto {
	cu := currency.Def().Subunits // SAT expects tax total amounts with no more decimals than supported by the currency

	if rate.Percent == nil {
		return &Impuesto{
			Base:       rate.Base.Rescale(cu).String(),
			Impuesto:   taxCategoryMap[catCode],
			TipoFactor: TipoFactorExento,
		}
	}

	return &Impuesto{
		Base:       rate.Base.Rescale(cu).String(),
		Importe:    rate.Amount.Rescale(cu).String(),
		Impuesto:   taxCategoryMap[catCode],
		TasaOCuota: format.TaxPercent(rate.Percent),
		TipoFactor: TipoFactorTasa,
	}
}

func newConceptoImpuestos(line *bill.Line) *ConceptoImpuestos {
	var traslados, retenciones []*Impuesto

	for _, tax := range line.Taxes {
		catDef := regime.CategoryDef(tax.Category)
		imp := newConceptoImpuesto(line, tax)
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

func newConceptoImpuesto(line *bill.Line, tax *tax.Combo) *Impuesto {
	if tax.Percent == nil {
		return &Impuesto{
			Base:       line.Total.String(),
			Impuesto:   taxCategoryMap[tax.Category],
			TipoFactor: TipoFactorExento,
		}
	}

	// GOBL doesn't provide an amount at line level, so we calculate it
	taxAmount := tax.Percent.Of(line.Total)
	return &Impuesto{
		Base:       line.Total.String(),
		Importe:    taxAmount.String(),
		Impuesto:   taxCategoryMap[tax.Category],
		TasaOCuota: format.TaxPercent(tax.Percent),
		TipoFactor: TipoFactorTasa,
	}
}
