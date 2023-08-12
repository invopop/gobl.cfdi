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
	var rf string
	rfID := org.IdentityForKey(supplier.Identities, mx.IdentityKeyFiscalRegime)
	if rfID != nil {
		rf = rfID.Code.String()
	}
	emisor := &Emisor{
		Rfc:           supplier.TaxID.Code.String(),
		Nombre:        supplier.Name,
		RegimenFiscal: rf,
	}
	return emisor
}

func newReceptor(customer *org.Party) *Receptor {
	var rf, usoCFDI string
	rfID := org.IdentityForKey(customer.Identities, mx.IdentityKeyFiscalRegime)
	if rfID != nil {
		rf = rfID.Code.String()
	}
	useID := org.IdentityForKey(customer.Identities, mx.IdentityKeyCFDIUse)
	if useID != nil {
		usoCFDI = useID.Code.String()
	}
	receptor := &Receptor{
		Rfc:                     customer.TaxID.Code.String(),
		Nombre:                  customer.Name,
		DomicilioFiscalReceptor: customer.TaxID.Zone.String(),
		RegimenFiscalReceptor:   rf,
		UsoCFDI:                 usoCFDI,
	}

	return receptor
}
