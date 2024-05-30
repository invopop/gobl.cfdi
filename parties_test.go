package cfdi_test

import (
	"testing"

	"github.com/invopop/gobl.cfdi/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmisor(t *testing.T) {
	t.Run("should return a Document with the Emisor data", func(t *testing.T) {
		doc, err := test.NewDocumentFrom("invoice-b2b-bare.json")
		require.NoError(t, err)

		e := doc.Emisor

		assert.Equal(t, "EKU9003173C9", e.Rfc)
		assert.Equal(t, "ESCUELA KEMPER URGATE", e.Nombre)
		assert.Equal(t, "601", e.RegimenFiscal)
	})
}

func TestReceptor(t *testing.T) {
	t.Run("should return a Document with the Receptor data", func(t *testing.T) {
		doc, err := test.NewDocumentFrom("invoice-b2b-bare.json")
		require.NoError(t, err)

		r := doc.Receptor

		assert.Equal(t, "URE180429TM6", r.Rfc)
		assert.Equal(t, "UNIVERSIDAD ROBOTICA ESPAÑOLA", r.Nombre)
		assert.Equal(t, "65000", r.DomicilioFiscalReceptor)
		assert.Equal(t, "601", r.RegimenFiscalReceptor)
		assert.Equal(t, "G01", r.UsoCFDI)
	})
}
