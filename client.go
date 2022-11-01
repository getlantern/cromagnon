package cromagnon

import (
	"crypto/x509"
	"fmt"
	"net"
	"strconv"

	"github.com/sagernet/cronet-go"
)

type ClientOptions struct {
	Addr               string            // Host/port to connect to
	Path               string            // Path on server to request
	PinnedCert         *x509.Certificate // Only accept this certificate
	InsecureSkipVerify bool              // If true, don't perform any certificate verification
	UseH3              bool              // Connect using HTTP/3 instead of HTTP/2 (default)
}

type Client struct {
	engine       cronet.Engine
	streamEngine cronet.StreamEngine
	url          string
}

// NewClient
// construct a new cromagnon client with the given options
func NewClient(config *ClientOptions) (*Client, error) {

	engine := cronet.NewEngine()
	params := cronet.NewEngineParams()

	if config.PinnedCert != nil || config.InsecureSkipVerify {
		verifier, err := cronet.CreatePinnedCertVerifier(config.PinnedCert, config.InsecureSkipVerify)
		if err != nil {
			return nil, err
		}
		engine.SetMockCertVerifierForTesting(verifier)
	}

	if config.UseH3 {
		params.SetEnableQuic(true)
		params.SetEnableHTTP2(false)
		hint := cronet.NewQuicHint()

		host, portStr, err := net.SplitHostPort(config.Addr)
		if err != nil {
			return nil, err
		}
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, err
		}
		hint.SetHost(host)
		hint.SetPort(int32(port))
		hint.SetAlternatePort(int32(port))
		params.AddQuicHint(hint)
	} else {
		params.SetEnableQuic(false)
		params.SetEnableHTTP2(true)
	}

	engine.StartWithParams(params)
	params.Destroy()
	streamEngine := engine.StreamEngine()

	c := &Client{
		url:          fmt.Sprintf("https://%s/%s/", config.Addr, config.Path),
		engine:       engine,
		streamEngine: streamEngine,
	}

	return c, nil
}

// Dial prepares a new net.Conn using the configuration of the client.
// The result may be a multiplexed connection using a previously
// established connection or may establish a new connection to
// the configured server.
func (c *Client) Dial() (net.Conn, error) {
	bidiConn := c.streamEngine.CreateConn(true, false)
	err := bidiConn.Start("GET", c.url, nil, 0, false)
	if err != nil {
		return nil, err
	}
	return bidiConn, nil
}

// Close closes all connections established with this Client.
// Dial should not be called after the client is closed.
func (c *Client) Close() error {
	c.engine.Shutdown()
	c.engine.Destroy()
	return nil
}
