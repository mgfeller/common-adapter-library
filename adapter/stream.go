package adapter

import (
	"github.com/layer5io/gokit/errors"
)

type Event struct {
	Operationid string `json:"operationid,omitempty"`
	EType       int32  `json:"type,string,omitempty"`
	Summary     string `json:"summary,omitempty"`
	Details     string `json:"details,omitempty"`
}

func (h *BaseAdapter) StreamErr(e *Event, err error) {
	h.Log.Err(errors.GetCode(err), err.Error())
	e.EType = 2
	*h.Channel <- e
}

func (h *BaseAdapter) StreamInfo(e *Event) {
	h.Log.Info("Sending event")
	e.EType = 0
	*h.Channel <- e
}
