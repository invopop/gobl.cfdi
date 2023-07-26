package cfdi_test

import (
	"testing"

	"github.com/invopop/gobl.cfdi/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCfdiRelacionados(t *testing.T) {
	t.Run("should return a Document with the CfdiRelacionados data", func(t *testing.T) {
		doc, err := test.NewDocumentFrom("credit_note.json")
		require.NoError(t, err)

		assert.Equal(t, "01", doc.CfdiRelacionados.TipoRelacion)

		rel := doc.CfdiRelacionados.CfdiRelacionado[0]

		assert.Equal(t, "1fac4464-1111-0000-1111-cd37179db12e", rel.UUID)
	})
}
