package addendas

import "github.com/invopop/gobl/bill"

// For returns a set of addenda objects for the given invoice.
func For(inv *bill.Invoice) []any {
	list := make([]any, 0)

	if isMabe(inv) {
		list = append(list, newMabe(inv))
	}

	return list
}
