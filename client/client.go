package client

import (
	"beowulf-go/api"
	"beowulf-go/config"
	"beowulf-go/transports"
	"beowulf-go/transports/http"
	"beowulf-go/transports/websocket"
	"github.com/pkg/errors"
	"net/url"
)

var (
	ErrInitializeTransport = errors.New("Failed to initialize transport.")
)

// Client can be used to access BEOWULF remote APIs.
// There is a public field for every BEOWULF API available,
type Client struct {
	cc transports.CallCloser

	chainID string

	AsyncProtocol bool

	// Database represents database_api.
	API *api.API

	// Current private keys for operations
	CurrentKeys *Keys
}

// NewClient creates a new RPC client that use the given CallCloser internally.
// Initialize only server present API. Absent API initialized as nil value.
func NewClient(s string, isTestNet bool) (*Client, error) {
	// Parse URL
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}

	// Initializing Transport
	var call transports.CallCloser
	switch u.Scheme {
	case "wss", "ws":
		call, err = websocket.NewTransport(s)
		if err != nil {
			return nil, err
		}
	case "https", "http":
		call, err = http.NewTransport(s)
		if err != nil {
			return nil, err
		}
	default:
		return nil, ErrInitializeTransport
	}
	client := &Client{cc: call}

	client.AsyncProtocol = true

	client.API = api.NewAPI(client.cc)

	if isTestNet {
		client.chainID = config.CHAIN_ID_TESTNET
	} else {
		client.chainID = config.CHAIN_ID_MAINNET
	}
	return client, nil
}

// Close should be used to close the client when no longer needed.
// It simply calls Close() on the underlying CallCloser.
func (client *Client) Close() error {
	return client.cc.Close()
}
