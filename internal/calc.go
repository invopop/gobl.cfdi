package internal

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
)

// TotalInvoiceDiscount calculates the total discount for the invoice.
func TotalInvoiceDiscount(i *bill.Invoice) num.Amount {
	td := i.Currency.Def().Zero() // currency's precision is required by the SAT
	for _, l := range i.Lines {
		td = td.Add(TotalLineDiscount(l))
	}
	return td
}

// TotalLineDiscount calculates the total discount for the line.
func TotalLineDiscount(l *bill.Line) num.Amount {
	td := num.MakeAmount(0, l.Sum.Exp()) // discount's precision must match the "Importe" field's one
	for _, d := range l.Discounts {
		td = td.Add(d.Amount)
	}
	return td
}
