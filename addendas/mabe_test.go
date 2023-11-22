package addendas_test

import (
	"testing"

	"github.com/invopop/gobl.cfdi/addendas"
	"github.com/invopop/gobl.cfdi/test"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddendaMabeValidation(t *testing.T) {
	env, err := test.LoadTestEnvelope("bare-minimum-invoice.json")
	require.NoError(t, err)

	inv := env.Extract().(*bill.Invoice)

	// Prepare the invoice to be raise all Mabe validation errors.
	inv.Type = bill.InvoiceTypeProforma
	inv.Supplier.Identities = []*org.Identity{
		{
			Key:  addendas.MabeKeyIdentityProviderID,
			Code: "12345",
		},
	}

	// Check every validation and then fix it.
	assertValidationError(t, inv, "type: must be a valid value")
	inv.Type = bill.InvoiceTypeStandard

	assertValidationError(t, inv, "delivery: cannot be blank")
	inv.Delivery = &bill.Delivery{}

	assertValidationError(t, inv, "delivery: (receiver: cannot be blank")
	inv.Delivery.Receiver = &org.Party{
		Name: "Test Receiver",
	}

	assertValidationError(t, inv, "delivery: (receiver: (identities: missing key mx-mabe-plant-id")
	inv.Delivery.Receiver.Identities = []*org.Identity{
		{
			Key:  addendas.MabeKeyIdentityPlantID,
			Code: "S001",
		},
	}

	assertValidationError(t, inv, "lines: (0: (item: (identities: missing key mx-mabe-item-id")
	inv.Lines[0].Item.Identities = []*org.Identity{
		{
			Key:  addendas.MabeKeyIdentityItemID,
			Code: "12345",
		},
	}

	assertValidationError(t, inv, "ordering: cannot be blank")
	inv.Ordering = &bill.Ordering{}

	assertValidationError(t, inv, "ordering: (code: cannot be blank")
	inv.Ordering.Code = "12345"

	assertValidationError(t, inv, "ordering: (identities: missing key mx-mabe-ref1")
	inv.Ordering.Identities = []*org.Identity{
		{
			Key:  addendas.MabeKeyIdentityRef1,
			Code: "12345",
		},
	}

	// All validation errors must be fixed by now.
	_, err = addendas.For(inv)
	require.NoError(t, err)
}

func assertValidationError(t *testing.T, inv *bill.Invoice, expected string) {
	_, err := addendas.For(inv)
	require.Error(t, err)
	assert.Contains(t, err.Error(), expected)
}
