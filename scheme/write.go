package scheme

// WriteTransaction
type WriteTransaction struct {
	Context     WriteContext `json:"context"`
	Transaction string       `json:"transaction"`
}

// WriteContext
type WriteContext struct {
	Action string   `json:"action"`
	Raw    []string `json:"raw"`
}
