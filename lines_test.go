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
		doc, err := test.NewDocumentFrom("bare-minimum-invoice.json")
		require.NoError(t, err)

		c := doc.Conceptos.Concepto[0]

		assert.Equal(t, "01010101", c.ClaveProdServ)
		assert.Equal(t, "2", c.Cantidad)
		assert.Equal(t, "H87", c.ClaveUnidad)
		assert.Equal(t, "Cigarros", c.Descripcion)
		assert.Equal(t, "100.00", c.ValorUnitario)
		assert.Equal(t, "200.00", c.Importe)
		assert.Equal(t, "02", c.ObjetoImp)
	})
}
