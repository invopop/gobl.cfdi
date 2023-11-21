package addendas

import (
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/invopop/gobl.cfdi/internal"
	"github.com/invopop/gobl.cfdi/internal/format"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
)

// Mabe schema constants
const (
	MabeVersion        = "1.0"
	MabeNamespace      = "https://recepcionfe.mabempresa.com/cfd/addenda/v1"
	MabeSchemaLocation = "https://recepcionfe.mabempresa.com/cfd/addenda/v1/mabev1.xsd"
	MabeNotApplicable  = "NA"
)

// TipoDocumento valid values
const (
	MabeTipoDocumentoFactura     = "FACTURA"
	MabeTipoDocumentoNotaCredito = "NOTA CREDITO"
	MabeTipoDocumentoNotaCargo   = "NOTA CARGO"
)

// Mabe specific identity codes.
const (
	MabeKeyIdentityProviderID = "mx-mabe-provider-id"
	MabeKeyIdentityRef1       = "mx-mabe-ref1"
	MabeKeyIdentityRef2       = "mx-mabe-ref2"
	MabeKeyIdentityPlantID    = "mx-mabe-plant-id"
	MabeKeyIdentityItemID     = "mx-mabe-item-id"
)

// MabeFactura is the root element of the Mabe addendum
type MabeFactura struct {
	XMLName        xml.Name `xml:"mabe:Factura"`
	Namespace      string   `xml:"xmlns:mabe,attr"`
	SchemaLocation string   `xml:"xsi:schemaLocation,attr"`

	Version       string `xml:"version,attr"`
	TipoDocumento string `xml:"tipoDocumento,attr"`
	Folio         string `xml:"folio,attr"`
	Fecha         string `xml:"fecha,attr"`
	OrdenCompra   string `xml:"ordenCompra,attr"`
	Referencia1   string `xml:"referencia1,attr"`
	Referencia2   string `xml:"referencia2,attr,omitempty"`

	Moneda    *MabeMoneda     `xml:"mabe:Moneda"`
	Proveedor *MabeProveedor  `xml:"mabe:Proveedor"`
	Entrega   *MabeEntrega    `xml:"mabe:Entrega"`
	Detalles  *[]*MabeDetalle `xml:"mabe:Detalles>mabe:Detalle"`

	Descuentos  *MabeDescuentos  `xml:"mabe:Descuentos,omitempty"`
	Subtotal    *MabeImporte     `xml:"mabe:Subtotal"`
	Traslados   *[]*MabeImpuesto `xml:"mabe:Traslados>mabe:Traslado"`
	Retenciones *[]*MabeImpuesto `xml:"mabe:Retenciones>mabe:Retencion"`
	Total       *MabeImporte     `xml:"mabe:Total"`
}

// MabeMoneda carries the data about the invoice's currency
type MabeMoneda struct {
	TipoMoneda      string `xml:"tipoMoneda,attr"`
	TipoCambio      string `xml:"tipoCambio,attr,omitempty"`      // Not implemented yet
	ImporteConLetra string `xml:"importeConLetra,attr,omitempty"` // Not implemented yet
}

// MabeProveedor carries the data about the invoice's supplier
type MabeProveedor struct {
	Codigo string `xml:"codigo,attr"`
}

// MabeEntrega carries the data about the invoice's delivery
type MabeEntrega struct {
	PlantaEntrega string `xml:"plantaEntrega,attr"`
	Calle         string `xml:"calle,attr,omitempty"`
	NoExterior    string `xml:"noExterior,attr,omitempty"`
	NoInterior    string `xml:"noInterior,attr,omitempty"`
	CodigoPostal  string `xml:"codigoPostal,attr,omitempty"`
}

// MabeDetalle carries the data about one invoice's line
type MabeDetalle struct {
	NoLineaArticulo int    `xml:"noLineaArticulo,attr"`
	CodigoArticulo  string `xml:"codigoArticulo,attr"`
	Descripcion     string `xml:"descripcion,attr"` //nolint:misspell
	Unidad          string `xml:"unidad,attr"`
	Cantidad        string `xml:"cantidad,attr"`
	PrecioSinIva    string `xml:"precioSinIva,attr"`
	ImporteSinIva   string `xml:"importeSinIva,attr"`
	PrecioConIva    string `xml:"precioConIva,attr,omitempty"`  // Not implemented yet
	ImporteConIva   string `xml:"importeConIva,attr,omitempty"` // Not implemented yet
}

// MabeImporte carries the data about an invoice's total
type MabeImporte struct {
	Importe string `xml:"importe,attr"`
}

// MabeImpuesto carries the data about an invoice's tax
type MabeImpuesto struct {
	Tipo    string `xml:"tipo,attr"`
	Tasa    string `xml:"tasa,attr"`
	Importe string `xml:"importe,attr"`
}

// MabeDescuentos carries the data about an invoice's discount
type MabeDescuentos struct {
	Tipo        string `xml:"tipo,attr"`
	Descripcion string `xml:"descripcion,attr"` //nolint:misspell
	Importe     string `xml:"importe,attr"`
}

func isMabe(inv *bill.Invoice) bool {
	if inv.Supplier == nil {
		return false
	}
	id := extractIdentity(inv.Supplier.Identities, MabeKeyIdentityProviderID)
	return id != cbc.CodeEmpty
}

// newMabe provides a new Mabe addenda.
func newMabe(inv *bill.Invoice) (*MabeFactura, error) {
	tipoDocumento, err := mapMabeTipoDocumento(inv)
	if err != nil {
		return nil, err
	}
	if inv.Ordering == nil {
		return nil, errors.New("missing ordering field")
	}

	f := &MabeFactura{
		Namespace:      MabeNamespace,
		SchemaLocation: format.SchemaLocation(MabeNamespace, MabeSchemaLocation),

		Version:       MabeVersion,
		TipoDocumento: tipoDocumento,
		Folio:         formatMabeFolio(inv),
		Fecha:         inv.IssueDate.String(),
		OrdenCompra:   inv.Ordering.Code,
		Referencia1:   extractIdentity(inv.Ordering.Identities, MabeKeyIdentityRef1).String(),
		Referencia2:   "NA",

		Moneda:     newMabeMoneda(inv),
		Proveedor:  newMabeProveedor(inv),
		Entrega:    newMabeEntrega(inv),
		Descuentos: newMabeDescuentos(inv),
		Detalles:   newMabeDetalles(inv),

		Subtotal: newMabeImporte(inv.Totals.Sum),
		Total:    newMabeImporte(inv.Totals.TotalWithTax),
	}

	setMabeTaxes(inv, f)

	return f, nil
}

func mapMabeTipoDocumento(inv *bill.Invoice) (string, error) {
	switch inv.Type {
	case bill.InvoiceTypeStandard:
		return MabeTipoDocumentoFactura, nil
	case bill.InvoiceTypeCreditNote:
		return MabeTipoDocumentoNotaCredito, nil
	case bill.InvoiceTypeDebitNote:
		return MabeTipoDocumentoNotaCargo, nil
	default:
		return "", fmt.Errorf("invalid invoice type: %s", inv.Type)
	}
}

func newMabeMoneda(inv *bill.Invoice) *MabeMoneda {
	return &MabeMoneda{TipoMoneda: string(inv.Currency)}
}

func newMabeProveedor(inv *bill.Invoice) *MabeProveedor {
	if inv.Supplier == nil {
		return nil
	}
	id := extractIdentity(inv.Supplier.Identities, MabeKeyIdentityProviderID)
	return &MabeProveedor{
		Codigo: id.String(),
	}
}

func newMabeEntrega(inv *bill.Invoice) *MabeEntrega {
	rec := inv.Delivery.Receiver
	id := extractIdentity(rec.Identities, MabeKeyIdentityPlantID)
	e := &MabeEntrega{
		PlantaEntrega: id.String(),
	}

	if len(rec.Addresses) > 0 {
		addr := rec.Addresses[0]

		e.Calle = addr.Street
		e.NoExterior = addr.Number
		e.NoInterior = MabeNotApplicable
		e.CodigoPostal = addr.Code
	}

	return e
}

func newMabeDescuentos(inv *bill.Invoice) *MabeDescuentos {
	d := internal.TotalInvoiceDiscount(inv)

	if d.IsZero() {
		return nil
	}

	return &MabeDescuentos{
		Tipo:        MabeNotApplicable,
		Descripcion: MabeNotApplicable, //nolint:misspell
		Importe:     d.String(),
	}
}

func newMabeDetalles(inv *bill.Invoice) *[]*MabeDetalle {
	var detalles []*MabeDetalle

	for _, line := range inv.Lines {
		id := extractIdentity(line.Item.Identities, MabeKeyIdentityItemID)
		d := &MabeDetalle{
			NoLineaArticulo: line.Index,
			CodigoArticulo:  id.String(),
			Descripcion:     line.Item.Name, //nolint:misspell
			Unidad:          internal.ClaveUnidad(line),
			Cantidad:        line.Quantity.String(),
			PrecioSinIva:    line.Item.Price.String(),
			ImporteSinIva:   line.Sum.String(),
		}

		detalles = append(detalles, d)
	}

	return &detalles
}

func newMabeImporte(amount num.Amount) *MabeImporte {
	return &MabeImporte{
		Importe: amount.String(),
	}
}

func setMabeTaxes(inv *bill.Invoice, mabe *MabeFactura) {
	var traslados, retenciones []*MabeImpuesto

	for _, cat := range inv.Totals.Taxes.Categories {
		catDef := inv.TaxRegime().Category(cat.Code)

		for _, rate := range cat.Rates {
			t := &MabeImpuesto{
				Tipo:    catDef.Name.In(i18n.ES),
				Tasa:    format.TaxPercent(rate.Percent),
				Importe: rate.Amount.String(),
			}

			if catDef.Retained {
				retenciones = append(retenciones, t)
			} else {
				traslados = append(traslados, t)
			}
		}
	}

	if len(traslados) > 0 {
		mabe.Traslados = &traslados
	}

	if len(retenciones) > 0 {
		mabe.Retenciones = &retenciones
	}
}

func formatMabeFolio(inv *bill.Invoice) string {
	return fmt.Sprintf("%s%s", inv.Series, inv.Code)
}

func extractIdentity(ids []*org.Identity, key cbc.Key) cbc.Code {
	if ids == nil {
		return ""
	}
	for _, id := range ids {
		if id.Key == key {
			return id.Code
		}
	}
	return ""
}
