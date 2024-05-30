package cfdi_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/invopop/gobl"
	cfdi "github.com/invopop/gobl.cfdi"
	"github.com/invopop/gobl.cfdi/test"
	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/xsd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var updateOut = flag.Bool("update", false, "Update the JSON and XML files in test/data and test/data/out")

func TestXMLGeneration(t *testing.T) {
	schema, err := loadSchema()
	require.NoError(t, err)

	examples, err := lookupExamples()
	require.NoError(t, err)

	for _, example := range examples {
		name := fmt.Sprintf("should convert %s example file successfully", example)

		t.Run(name, func(t *testing.T) {
			data, env, err := convertExample(example)
			require.NoError(t, err)

			outPath := filepath.Join(test.GetDataPath(), "out", strings.TrimSuffix(example, ".json")+".xml")

			if *updateOut {
				errs := validateDoc(schema, data)
				for _, e := range errs {
					assert.NoError(t, e)
				}
				if len(errs) > 0 {
					assert.Fail(t, "Invalid XML:\n"+string(data))
					return
				}

				// Write the XML file
				err = os.WriteFile(outPath, data, 0644)
				require.NoError(t, err)

				// Update the (just re-calculated) GOBL file as well
				jsData, _ := json.MarshalIndent(env, "", "\t")
				err = os.WriteFile(filepath.Join(test.GetDataPath(), example), append(jsData, '\n'), 0644)
				require.NoError(t, err)

				return
			}

			expected, err := os.ReadFile(outPath)

			require.False(t, os.IsNotExist(err), "output file %s missing, run tests with `--update` flag to create", filepath.Base(outPath))
			require.NoError(t, err)
			require.Equal(t, string(expected), string(data), "output file %s does not match, run tests with `--update` flag to update", filepath.Base(outPath))
		})
	}
}

func loadSchema() (*xsd.Schema, error) {
	schemaPath := filepath.Join(test.GetTestPath(), "schema", "schema.xsd")
	schema, err := xsd.ParseFromFile(schemaPath)
	if err != nil {
		return nil, err
	}

	return schema, nil
}

func lookupExamples() ([]string, error) {
	examples, err := filepath.Glob(filepath.Join(test.GetDataPath(), "*.json"))
	if err != nil {
		return nil, err
	}

	for i, example := range examples {
		examples[i] = filepath.Base(example)
	}

	return examples, nil
}

func convertExample(example string) ([]byte, *gobl.Envelope, error) {
	env, err := test.LoadTestEnvelope(example)
	if err != nil {
		return nil, nil, err
	}

	doc, err := cfdi.NewDocument(env)
	if err != nil {
		return nil, nil, err
	}

	data, err := doc.Bytes()
	if err != nil {
		return nil, nil, err
	}

	return append(data, '\n'), env, nil
}

func validateDoc(schema *xsd.Schema, doc []byte) []error {
	xmlDoc, err := libxml2.ParseString(string(doc))
	if err != nil {
		return []error{err}
	}

	err = schema.Validate(xmlDoc)
	if err != nil {
		return err.(xsd.SchemaValidationError).Errors()
	}

	return nil
}
