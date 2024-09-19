package cfdi

import (
	addon "github.com/invopop/gobl/addons/mx/cfdi"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/mx"
)

// Emisor stores the invoice supplier data
type Emisor struct {
	Rfc           string `xml:",attr"`
	Nombre        string `xml:",attr"`
	RegimenFiscal string `xml:",attr"`
}

// Receptor stores the invoice customer data
type Receptor struct {
	Rfc                     string `xml:",attr"`
	Nombre                  string `xml:",attr"`
	DomicilioFiscalReceptor string `xml:",attr"`
	RegimenFiscalReceptor   string `xml:",attr"`
	UsoCFDI                 string `xml:",attr"`
	NumRegIdTrib            string `xml:",attr,omitempty"` //nolint:revive
	ResidenciaFiscal        string `xml:",attr,omitempty"`
}

func newEmisor(supplier *org.Party) *Emisor {
	emisor := &Emisor{
		Rfc:           supplier.TaxID.Code.String(),
		Nombre:        supplier.Name,
		RegimenFiscal: supplier.Ext[addon.ExtKeyFiscalRegime].String(),
	}
	return emisor
}

func newReceptor(customer *org.Party, issuePlace string) *Receptor {
	if customer == nil {
		return &Receptor{
			Nombre:                  NombreReceptorGenerico,
			Rfc:                     mx.TaxIdentityCodeGeneric.String(),
			RegimenFiscalReceptor:   RegimenFiscalSinObligaciones,
			UsoCFDI:                 UsoCFDISinEfectos,
			DomicilioFiscalReceptor: issuePlace,
		}
	}

	if customer.TaxID.Country != l10n.MX.Tax() {
		alpha3 := l10n.Countries().Code(customer.TaxID.Country.Code()).Alpha3

		return &Receptor{
			Nombre:                  customer.Name,
			Rfc:                     mx.TaxIdentityCodeForeign.String(),
			RegimenFiscalReceptor:   RegimenFiscalSinObligaciones,
			UsoCFDI:                 UsoCFDISinEfectos,
			DomicilioFiscalReceptor: issuePlace,
			NumRegIdTrib:            customer.TaxID.Code.String(),
			ResidenciaFiscal:        alpha3,
		}
	}

	return &Receptor{
		Nombre:                  customer.Name,
		Rfc:                     customer.TaxID.Code.String(),
		RegimenFiscalReceptor:   customer.Ext[addon.ExtKeyFiscalRegime].String(),
		UsoCFDI:                 customer.Ext[addon.ExtKeyUse].String(),
		DomicilioFiscalReceptor: customer.Ext[addon.ExtKeyPostCode].String(),
	}
}
