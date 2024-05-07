package cfdi_test

import (
	"testing"

	"github.com/invopop/gobl.cfdi/test"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regimes/mx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComprobanteIngreso(t *testing.T) {
	t.Run("should return a Document with the Comprobante data", func(t *testing.T) {
		doc, err := test.NewDocumentFrom("invoice.json")
		require.NoError(t, err)

		assert.Equal(t, "http://www.sat.gob.mx/cfd/4", doc.CFDINamespace)
		assert.Equal(t, "http://www.w3.org/2001/XMLSchema-instance", doc.XSINamespace)
		assert.Equal(t, "http://www.sat.gob.mx/cfd/4 http://www.sat.gob.mx/sitio_internet/cfd/4/cfdv40.xsd", doc.SchemaLocation)
		assert.Equal(t, "4.0", doc.Version)

		assert.Equal(t, "I", doc.TipoDeComprobante)
		assert.Equal(t, "LMC", doc.Serie)
		assert.Equal(t, "0010", doc.Folio)
		assert.Equal(t, "2023-05-29T00:00:00", doc.Fecha)
		assert.Equal(t, "26015", doc.LugarExpedicion)
		assert.Equal(t, "420.90", doc.SubTotal)
		assert.Equal(t, "200.20", doc.Descuento)
		assert.Equal(t, "211.36", doc.Total)
		assert.Equal(t, "MXN", doc.Moneda)
		assert.Equal(t, "01", doc.Exportacion)
		assert.Equal(t, "PUE", doc.MetodoPago)
		assert.Equal(t, "03", doc.FormaPago)
		assert.Equal(t, "Pago a 30 días.", doc.CondicionesDePago)

		assert.Nil(t, doc.Complemento)
	})

	t.Run("should return the proper MetodoPago and FormaPago", func(t *testing.T) {
		inv, _ := test.LoadTestInvoice("invoice.json")

		// No advances
		inv.Payment.Advances = nil
		doc, _ := test.GenerateCFDIFrom(inv)
		assert.Equal(t, "PPD", doc.MetodoPago)
		assert.Equal(t, "99", doc.FormaPago)

		// Partial settlement
		inv.Payment.Advances = append(inv.Payment.Advances, &pay.Advance{
			Percent:     num.NewPercentage(40, 2),
			Description: "First partial settlement",
			Key:         pay.MeansKeyCash,
		})
		doc, _ = test.GenerateCFDIFrom(inv)
		assert.Equal(t, "PPD", doc.MetodoPago)
		assert.Equal(t, "99", doc.FormaPago)

		// Full settlement
		inv.Payment.Advances = append(inv.Payment.Advances, &pay.Advance{
			Percent:     num.NewPercentage(60, 2),
			Description: "Second partial settlement",
			Key:         pay.MeansKeyOnline.With(mx.MeansKeyWallet),
		})
		doc, _ = test.GenerateCFDIFrom(inv)
		assert.Equal(t, "PUE", doc.MetodoPago)
		assert.Equal(t, "05", doc.FormaPago)
	})
}

func TestComprobanteEgreso(t *testing.T) {
	t.Run("should return a Document with the Comprobante data", func(t *testing.T) {
		doc, err := test.NewDocumentFrom("credit-note.json")
		require.NoError(t, err)

		assert.Equal(t, "E", doc.TipoDeComprobante)
		assert.Equal(t, "CN", doc.Serie)
		assert.Equal(t, "0003", doc.Folio)
	})
}
