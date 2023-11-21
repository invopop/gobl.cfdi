package internal

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/regimes/mx"
)

// Default keys
const (
	DefaultClaveUnidad = "ZZ" // Mutuamente definida
)

// ClaveUnidad determines the line item's "ClaveUnidad" value.
func ClaveUnidad(line *bill.Line) string {
	if line.Item.Unit == "" {
		return DefaultClaveUnidad
	}

	return string(line.Item.Unit.UNECE())
}

// ClaveProdServ determines the line's Product-Service code
func ClaveProdServ(line *bill.Line) string {
	if line.Item == nil {
		return ""
	}

	return string(line.Item.Ext[mx.ExtKeyCFDIProdServ])
}
