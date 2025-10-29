package internal

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
)

// TotalInvoiceDiscount calculates the total discount for the invoice.
func TotalInvoiceDiscount(i *bill.Invoice) *num.Amount {
	td := i.Currency.Def().Zero()
	for _, l := range i.Lines {
		ld := TotalLineDiscount(l)
		if ld != nil {
			td = td.MatchPrecision(*ld)
			td = td.Add(*ld)
		}
	}
	if td.IsZero() {
		return nil
	}
	return &td
}

// TotalLineDiscount calculates the total discount for the line.
func TotalLineDiscount(l *bill.Line) *num.Amount {
	// discount's precision must match the "Importe" field's one
	td := num.MakeAmount(0, l.Sum.Exp())
	for _, d := range l.Discounts {
		td = td.Add(d.Amount)
	}
	if td.IsZero() {
		return nil
	}
	return &td
}
