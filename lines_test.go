package cfdi_test

import (
	"testing"

	"github.com/invopop/gobl.cfdi/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConceptos(t *testing.T) {
	// nolint:misspell
	t.Run("should return a Document with the Conceptos data", func(t *testing.T) {
		doc, err := test.NewDocumentFrom("invoice.json")
		require.NoError(t, err)

		c := doc.Conceptos.Concepto[0]

		assert.Equal(t, "50211502", c.ClaveProdServ)
		assert.Equal(t, "2", c.Cantidad)
		assert.Equal(t, "H87", c.ClaveUnidad)
		assert.Equal(t, "Cigarros", c.Descripcion)
		assert.Equal(t, "200.2020", c.ValorUnitario)
		assert.Equal(t, "200.2020", c.Descuento)
		assert.Equal(t, "400.4040", c.Importe)
		assert.Equal(t, "02", c.ObjetoImp)
		assert.Equal(t, "H87", c.ClaveUnidad)
		assert.Len(t, c.Impuestos.Traslados.Traslado, 1)

		c = doc.Conceptos.Concepto[1]
		assert.Equal(t, "01", c.ObjetoImp)
		assert.Nil(t, c.Impuestos)
	})

	t.Run("should return the default ClaveUnidad when no unit is given", func(t *testing.T) {
		doc, err := test.NewDocumentFrom("bare-minimum-invoice.json")
		require.NoError(t, err)

		c := doc.Conceptos.Concepto[0]

		assert.Equal(t, "ZZ", c.ClaveUnidad)
	})
}
