// Package format contains helps to help format output.
package format

import (
	"fmt"

	"github.com/invopop/gobl/num"
)

// OptionalAmount provides empty string for zero amounts.
func OptionalAmount(a num.Amount) string {
	if a.IsZero() {
		return ""
	}

	return a.String()
}

// SchemaLocation provides a string with the namespace and schema location.
func SchemaLocation(namespace, schemaLocation string) string {
	return fmt.Sprintf("%s %s", namespace, schemaLocation)
}

// TaxPercent provides a string with the tax percentage rescaled according to
// CFDI requirements.
func TaxPercent(percent *num.Percentage) string {
	return percent.Amount.Rescale(6).String()
}
