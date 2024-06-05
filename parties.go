package cfdi

import (
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
		RegimenFiscal: supplier.Ext[mx.ExtKeyCFDIFiscalRegime].String(),
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

	if customer.TaxID.Country != l10n.MX {
		return &Receptor{
			Nombre:                  customer.Name,
			Rfc:                     mx.TaxIdentityCodeForeign.String(),
			RegimenFiscalReceptor:   RegimenFiscalSinObligaciones,
			UsoCFDI:                 UsoCFDISinEfectos,
			DomicilioFiscalReceptor: issuePlace,
			NumRegIdTrib:            customer.TaxID.Code.String(),
			ResidenciaFiscal:        customer.TaxID.Country.Alpha3(),
		}
	}

	return &Receptor{
		Nombre:                  customer.Name,
		Rfc:                     customer.TaxID.Code.String(),
		RegimenFiscalReceptor:   customer.Ext[mx.ExtKeyCFDIFiscalRegime].String(),
		UsoCFDI:                 customer.Ext[mx.ExtKeyCFDIUse].String(),
		DomicilioFiscalReceptor: customer.Ext[mx.ExtKeyCFDIPostCode].String(),
	}
}
