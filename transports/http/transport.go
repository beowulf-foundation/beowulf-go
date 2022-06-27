package http

import (
	"beowulf-go/config"
	"beowulf-go/types"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"math"
	"net/http"
	"sync"
	"time"
)

type Transport struct {
	Url    string
	client http.Client

	requestID uint64
	reqMutex  sync.Mutex
}

func NewTransport(url string) (*Transport, error) {
	timeout := time.Duration(config.HTTP_CONNECTION_TIMEOUT_SECOND * time.Second)

	return &Transport{
		client: http.Client{
			Timeout: timeout,
		},
		Url: url,
	}, nil
}

func (caller *Transport) Call(method string, args []interface{}, reply interface{}, scid string) error {
	caller.reqMutex.Lock()
	defer caller.reqMutex.Unlock()

	// increase request id
	if caller.requestID == math.MaxUint64 {
		caller.requestID = 0
	}
	caller.requestID++
	var request = types.RPCRequest{}
	if len(scid) > 0 {
		methodStr := fmt.Sprintf("%v", args[0])
		request = types.RPCRequest{
			Method: methodStr,
			JSON:   "2.0",
			ID:     caller.requestID,
			Params: args[1],
		}
	} else {
		request = types.RPCRequest{
			Method: method,
			JSON:   "2.0",
			ID:     caller.requestID,
			Params: args,
		}
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		return err
	}

	//resp, err := caller.client.Post(caller.Url, "application/json", bytes.NewBuffer(reqBody))
	req, err := http.NewRequest("POST", caller.Url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if len(scid) > 0 {
		req.Header.Set("scid", scid)
	}
	resp, err := caller.client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read body")
	}

	var rpcResponse types.RPCResponse
	if err = json.Unmarshal(respBody, &rpcResponse); err != nil {
		return errors.Wrapf(err, "failed to unmarshal response: %+v", string(respBody))
	}

	if rpcResponse.Error != nil {
		return rpcResponse.Error
	}

	if rpcResponse.Result != nil {
		if err := json.Unmarshal(*rpcResponse.Result, reply); err != nil {
			return errors.Wrapf(err, "failed to unmarshal rpc result: %+v", string(*rpcResponse.Result))
		}
	}

	return nil
}

func (caller *Transport) SetCallback(api string, method string, notice func(args json.RawMessage)) error {
	panic("not supported")
}

func (caller *Transport) Close() error {
	return nil
}

func check(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.Errorf("Error. URL: %s STATUS: %s\n", url, resp.StatusCode)
	}
	return nil
}
