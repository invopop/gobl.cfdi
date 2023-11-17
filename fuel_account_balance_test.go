package cfdi_test

import (
	"testing"

	cfdi "github.com/invopop/gobl.cfdi"
	"github.com/invopop/gobl.cfdi/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEstadoDeCuentaCombustible(t *testing.T) {
	t.Run("should return a Document with the EstadoDeCuentaCombustible data", func(t *testing.T) {
		doc, err := test.NewDocumentFrom("fuel-account-balance.json")
		require.NoError(t, err)

		require.Equal(t, 1, len(doc.Complementos))

		ecc := doc.Complementos[0].Content.(*cfdi.EstadoDeCuentaCombustible)

		assert.Equal(t, "0123456789", ecc.NumeroDeCuenta)
		assert.Equal(t, "246.13", ecc.SubTotal)
		assert.Equal(t, "400.00", ecc.Total)

		require.Equal(t, 2, len(ecc.Conceptos))

		c := ecc.Conceptos[0]

		assert.Equal(t, "1234", c.Identificador)
		assert.Equal(t, "2022-07-19T10:20:30", c.Fecha)
		assert.Equal(t, "RWT860605OF5", c.Rfc)
		assert.Equal(t, "8171650", c.ClaveEstacion)
		assert.Equal(t, "9.661", c.Cantidad)
		assert.Equal(t, "3", c.TipoCombustible)
		assert.Equal(t, "Diesel", c.NombreCombustible)
		assert.Equal(t, "2794668", c.FolioOperacion)
		assert.Equal(t, "12.743", c.ValorUnitario)
		assert.Equal(t, "123.11", c.Importe)
		assert.Equal(t, "LTR", c.Unidad)

		require.Equal(t, 2, len(c.Traslados))

		ct := c.Traslados[0]

		assert.Equal(t, "IVA", ct.Impuesto)
		assert.Equal(t, "0.160000", ct.TasaOCuota)
		assert.Equal(t, "19.70", ct.Importe)
	})
}
