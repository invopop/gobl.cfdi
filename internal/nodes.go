package internal

// Nodes is an auxiliary struct to marshal a sequence of arbitrary XML nodes,
// like the ones inside `cfdi:Complemento` or `cfdi:Addenda`.
type Nodes struct {
	Nodes []interface{} `xml:",omitempty"`
}
