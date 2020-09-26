package adapter

type Operation struct {
	Type       int32             `json:"type,string,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

type Operations map[string]*Operation

func (h *BaseAdapter) ListOperations() (Operations, error) {
	operations := make(Operations)
	err := h.Config.Operations(&operations)
	if err != nil {
		return nil, err
	}
	return operations, nil
}
