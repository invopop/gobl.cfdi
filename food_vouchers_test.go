package cfdi_test

import (
	"testing"

	cfdi "github.com/invopop/gobl.cfdi"
	"github.com/invopop/gobl.cfdi/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValesDeDespensa(t *testing.T) {
	t.Run("should return a Document with the ValesDeDespensa data", func(t *testing.T) {
		doc, err := test.NewDocumentFrom("food-vouchers.json")
		require.NoError(t, err)

		require.Equal(t, 1, len(doc.Complementos))

		vd := doc.Complementos[0].Content.(*cfdi.ValesDeDespensa)

		assert.Equal(t, "12345678901234567890", vd.RegistroPatronal)
		assert.Equal(t, "0123456789", vd.NumeroDeCuenta)
		assert.Equal(t, "30.52", vd.Total)

		require.Equal(t, 2, len(vd.Conceptos))

		c := vd.Conceptos[0]

		assert.Equal(t, "ABC1234", c.Identificador)
		assert.Equal(t, "2022-07-19T10:20:30", c.Fecha)
		assert.Equal(t, "JUFA7608212V6", c.Rfc)
		assert.Equal(t, "JUFA760821MDFRRR00", c.Curp)
		assert.Equal(t, "Adriana Juarez Fernández", c.Nombre)
		assert.Equal(t, "12345678901", c.NumSeguridadSocial)
		assert.Equal(t, "10.12", c.Importe)
	})
}
