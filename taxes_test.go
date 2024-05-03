package cfdi_test

import (
	"testing"

	"github.com/invopop/gobl.cfdi/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImpuestos(t *testing.T) {
	t.Run("should return a Document with the Impuestos data", func(t *testing.T) {
		doc, err := test.NewDocumentFrom("invoice.json")
		require.NoError(t, err)

		assert.Equal(t, "32.03", doc.Impuestos.TotalImpuestosTrasladados)

		tr := doc.Impuestos.Traslados.Traslado[0]

		assert.Equal(t, "200.20", tr.Base)   // SAT expects 2 decimals only
		assert.Equal(t, "32.03", tr.Importe) // SAT expects 2 decimals only
		assert.Equal(t, "002", tr.Impuesto)
		assert.Equal(t, "Tasa", tr.TipoFactor)
		assert.Equal(t, "0.160000", tr.TasaOCuota)
	})

	t.Run("should return a Document with the Impuestos data of each Concepto", func(t *testing.T) {
		doc, err := test.NewDocumentFrom("invoice.json")
		require.NoError(t, err)

		tr := doc.Conceptos.Concepto[0].Impuestos.Traslados.Traslado[0]

		assert.Equal(t, "200.2020", tr.Base)
		assert.Equal(t, "32.0323", tr.Importe)
		assert.Equal(t, "002", tr.Impuesto)
		assert.Equal(t, "Tasa", tr.TipoFactor)
		assert.Equal(t, "0.160000", tr.TasaOCuota)

		assert.Nil(t, doc.Conceptos.Concepto[1].Impuestos)
	})
}
