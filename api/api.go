package api

import (
	"encoding/json"

	"github.com/beowulf-foundation/beowulf-go/transports"
)

//API plug-in structure
type API struct {
	caller transports.Caller
}

//NewAPI plug-in initialization
func NewAPI(caller transports.Caller) *API {
	return &API{caller}
}

func (api *API) call(apiID string, method string, params, resp interface{}, scid string) error {
	if len(apiID) > 0 {
		return api.caller.Call("call", []interface{}{apiID, method, params}, resp, scid)
	}
	return api.caller.Call("call", []interface{}{method, params}, resp, scid)
}

func (api *API) setCallback(apiID string, method string, callback func(raw json.RawMessage)) error {
	return api.caller.SetCallback(apiID, method, callback)
}
