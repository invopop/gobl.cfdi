// Package addendas adds additional functionality for "Addendas" to the CFDI documents.
package addendas

import "github.com/invopop/gobl/bill"

// For returns a set of addenda objects for the given invoice.
func For(inv *bill.Invoice) ([]any, error) {
	list := make([]any, 0)

	if isMabe(inv) {
		ad, err := newMabe(inv)
		if err != nil {
			return nil, err
		}
		list = append(list, ad)
	}

	return list, nil
}
