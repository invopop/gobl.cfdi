package cfdi

import (
	"github.com/invopop/gobl.cfdi/internal/format"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
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

// Map of tax categories to SAT tax types
var taxCategoryMap = map[cbc.Code]string{
	mx.TaxCategoryISR:   "001",
	tax.CategoryVAT:     "002",
	mx.TaxCategoryRVAT:  "002",
	mx.TaxCategoryIEPS:  "003",
	mx.TaxCategoryRIEPS: "003",
}

func newImpuestos(totals *bill.Totals, lines []*bill.Line, currency currency.Code) *Impuestos {
	if totals.Taxes == nil {
		return nil
	}
	var traslados, retenciones []*Impuesto

	for _, cat := range totals.Taxes.Categories {
		for _, rate := range cat.Rates {
			imp := newImpuesto(cat.Code, rate, currency)
			if cat.Retained {
				// Clear out fields not supported by retained totals
				imp.Base = ""
				imp.TasaOCuota = ""
				imp.TipoFactor = ""

				retenciones = append(retenciones, imp)
			} else {
				traslados = append(traslados, imp)
			}
		}
	}

	traslados = taxesAddLineCharges(traslados, lines, currency)

	impuestos := &Impuestos{}
	empty := true
	if len(traslados) > 0 {
		impuestos.Traslados = &Traslados{traslados}

		// Set tax total only for non-exempt taxes
		for _, t := range traslados {
			if t.TipoFactor != TipoFactorExento {
				impuestos.TotalImpuestosTrasladados = addStringAmounts(impuestos.TotalImpuestosTrasladados, t.Importe)
			}
		}

		empty = false
	}
	if len(retenciones) > 0 {
		impuestos.Retenciones = &Retenciones{retenciones}
		for _, r := range retenciones {
			impuestos.TotalImpuestosRetenidos = addStringAmounts(impuestos.TotalImpuestosRetenidos, r.Importe)
		}
		empty = false
	}
	if empty {
		return nil
	}

	return impuestos
}

func taxesAddLineCharges(traslados []*Impuesto, lines []*bill.Line, cur currency.Code) []*Impuesto {
	taxes := []*Impuesto{}
	// generate the lines that need to be there
	for _, line := range lines {
		for _, charge := range line.Charges {
			tl := newImpuestoFromLineCharge(line, charge)
			if tl != nil {
				// ensure bases are rounded as they may have come from quantities
				b := stringToAmount(tl.Base)
				tl.Base = b.Rescale(cur.Def().Subunits).String()
				taxes = append(taxes, tl)
			}
		}
	}
	if len(taxes) == 0 {
		return traslados
	}

	// see if any of the new lines are already present to add to the existing ones
	for _, tl := range taxes {
		found := false
		for _, tlt := range traslados {
			if tlt.Impuesto == tl.Impuesto &&
				tlt.TipoFactor == tl.TipoFactor &&
				tlt.TasaOCuota == tl.TasaOCuota {

				tlt.Base = addStringAmounts(tlt.Base, tl.Base)
				tlt.Importe = addStringAmounts(tlt.Importe, tl.Importe)
				found = true
			}
		}
		if !found {
			traslados = append(traslados, tl)
		}
	}

	return traslados
}

// addStringAmounts is used to add to amounts together when the source is a
// string instead of a number. Any errors occur during parsing, we'll just
// return the other amount.
func addStringAmounts(a, b string) string {
	return stringToAmount(a).Add(stringToAmount(b)).String()
}

func stringToAmount(a string) num.Amount {
	if a == "" {
		a = "0.00"
	}
	am, err := num.AmountFromString(a)
	if err != nil {
		return currency.MXN.Def().Zero()
	}
	return am
}

func newImpuesto(catCode cbc.Code, rate *tax.RateTotal, currency currency.Code) *Impuesto {
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
	if line.Total == nil {
		// nothing to do when no line total
		return nil
	}
	var traslados, retenciones []*Impuesto

	for _, tax := range line.Taxes {
		catDef := regime.CategoryDef(tax.Category)
		imp := newImpuestoFromCombo(line, tax)
		if catDef.Retained {
			retenciones = append(retenciones, imp)
		} else {
			traslados = append(traslados, imp)
		}
	}
	for _, charge := range line.Charges {
		imp := newImpuestoFromLineCharge(line, charge)
		if imp != nil {
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

func newImpuestoFromCombo(line *bill.Line, tax *tax.Combo) *Impuesto {
	if tax.Percent == nil {
		return &Impuesto{
			Base:       line.Total.String(),
			Impuesto:   taxCategoryMap[tax.Category],
			TipoFactor: TipoFactorExento,
		}
	}

	// GOBL doesn't provide an amount at line level, so we calculate it
	taxAmount := tax.Percent.Of(*line.Total)
	return &Impuesto{
		Base:       line.Total.String(),
		Importe:    taxAmount.String(),
		Impuesto:   taxCategoryMap[tax.Category],
		TasaOCuota: format.TaxPercent(tax.Percent),
		TipoFactor: TipoFactorTasa,
	}
}

func newImpuestoFromLineCharge(line *bill.Line, charge *bill.LineCharge) *Impuesto {
	if charge.Code != mx.TaxCategoryIEPS {
		// only handle IEPS at the moment
		return nil
	}

	i := &Impuesto{
		Impuesto: taxCategoryMap[charge.Code],
		Importe:  charge.Amount.String(),
	}
	if charge.Percent != nil {
		i.Base = line.Sum.String()
		i.TasaOCuota = format.TaxPercent(charge.Percent)
		i.TipoFactor = TipoFactorTasa
	} else if charge.Rate != nil {
		q := line.Quantity
		if charge.Quantity != nil {
			q = *charge.Quantity
		}
		i.Base = q.String()
		i.TasaOCuota = format.TaxRate(*charge.Rate)
		i.TipoFactor = TipoFactorCuota
	} else {
		// Not enough details to process, ignore.
		return nil
	}

	return i
}
