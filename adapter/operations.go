package adapter

// Operation holds the informormation for list of operations
type Operation struct {
	Type       int32             `json:"type,string,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

// Operations hold a map of Operation objects
type Operations map[string]*Operation

// ListOperations lists the operations available
func (h *BaseAdapter) ListOperations() (Operations, error) {
	operations := make(Operations, 0)
	err := h.Config.Operations(&operations)
	if err != nil {
		return nil, err
	}
	return operations, nil
}
