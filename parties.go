package cfdi

import (
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
}

func newEmisor(supplier *org.Party) *Emisor {
	emisor := &Emisor{
		Rfc:           supplier.TaxID.Code.String(),
		Nombre:        supplier.Name,
		RegimenFiscal: supplier.Ext[mx.ExtKeyCFDIFiscalRegime].String(),
	}
	return emisor
}

func newReceptor(customer *org.Party) *Receptor {
	receptor := &Receptor{
		Rfc:                     customer.TaxID.Code.String(),
		Nombre:                  customer.Name,
		DomicilioFiscalReceptor: customer.Ext[mx.ExtKeyCFDIPostCode].String(),
		RegimenFiscalReceptor:   customer.Ext[mx.ExtKeyCFDIFiscalRegime].String(),
		UsoCFDI:                 customer.Ext[mx.ExtKeyCFDIUse].String(),
	}

	return receptor
}
