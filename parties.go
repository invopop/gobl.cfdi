package cfdi

import "github.com/invopop/gobl/org"

type Emisor struct {
	Rfc           string `xml:",attr"`
	Nombre        string `xml:",attr"`
	RegimenFiscal string `xml:",attr"`
}

type Receptor struct {
	Rfc                     string `xml:",attr"`
	Nombre                  string `xml:",attr"`
	DomicilioFiscalReceptor string `xml:",attr"`
	RegimenFiscalReceptor   string `xml:",attr"`
	UsoCFDI                 string `xml:",attr"`
}

func NewEmisor(supplier *org.Party) *Emisor {
	emisor := &Emisor{
		Rfc:           supplier.TaxID.Code.String(),
		Nombre:        supplier.Name,
		RegimenFiscal: RegimenFiscalGeneral,
	}

	return emisor
}

func NewReceptor(customer *org.Party) *Receptor {
	receptor := &Receptor{
		Rfc:                     customer.TaxID.Code.String(),
		Nombre:                  customer.Name,
		DomicilioFiscalReceptor: customer.Addresses[0].Code,
		RegimenFiscalReceptor:   RegimenFiscalGeneral,
		UsoCFDI:                 UsoCFDIGastosGenerales,
	}

	return receptor
}