package cfdi_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/invopop/gobl.cfdi/test"
	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/xsd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComprobante(t *testing.T) {
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
		assert.Equal(t, "200.00", doc.SubTotal)
		assert.Equal(t, "232.00", doc.Total)
		assert.Equal(t, "MXN", doc.Moneda)
		assert.Equal(t, "01", doc.Exportacion)
		assert.Equal(t, "PUE", doc.MetodoPago)
		assert.Equal(t, "03", doc.FormaPago)
		assert.Equal(t, "Pago a 30 días.", doc.CondicionesDePago)
	})
}

func TestXMLGeneration(t *testing.T) {
	schemaPath := filepath.Join(test.GetTestPath(), "schema", "cfdv40.xsd")
	schema, err := xsd.ParseFromFile(schemaPath)
	defer schema.Free()
	require.NoError(t, err)

	tests := []string{
		"bare-minimum-invoice.json",
		"invoice.json",
	}

	for _, testFile := range tests {
		name := fmt.Sprintf("should generate a schema-valid XML from %s", testFile)
		t.Run(name, func(t *testing.T) {
			doc, err := test.NewDocumentFrom(testFile)
			require.NoError(t, err)

			data, err := doc.Bytes()
			require.NoError(t, err)

			xmlDoc, err := libxml2.ParseString(string(data))
			require.NoError(t, err)

			err = schema.Validate(xmlDoc)
			if err != nil {
				for _, e := range err.(xsd.SchemaValidationError).Errors() {
					require.NoError(t, e)
				}
			}
		})
	}
}
