package cfdi_test

import (
	"testing"

	"github.com/invopop/gobl.cfdi/test"
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
		assert.Equal(t, "400.40", doc.SubTotal)
		assert.Equal(t, "200.20", doc.Descuento)
		assert.Equal(t, "190.86", doc.Total)
		assert.Equal(t, "MXN", doc.Moneda)
		assert.Equal(t, "01", doc.Exportacion)
		assert.Equal(t, "PUE", doc.MetodoPago)
		assert.Equal(t, "03", doc.FormaPago)
		assert.Equal(t, "Pago a 30 días.", doc.CondicionesDePago)

		assert.Equal(t, 0, len(doc.Complementos))
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
