package cromagnon

import (
	"crypto/x509"
	"fmt"
	"net"
	"strconv"

	"github.com/sagernet/cronet-go"
)

type ClientOptions struct {
	Addr               string
	Path               string
	PinnedCert         *x509.Certificate
	InsecureSkipVerify bool
	UseH3              bool
}

type Client struct {
	engine       cronet.Engine
	streamEngine cronet.StreamEngine
	url          string
}

func NewClient(config *ClientOptions) (*Client, error) {

	engine := cronet.NewEngine()
	params := cronet.NewEngineParams()

	verifier, err := cronet.CreatePinnedCertVerifier(config.PinnedCert, config.InsecureSkipVerify)
	if err != nil {
		return nil, err
	}
	engine.SetMockCertVerifierForTesting(verifier)

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

func (c *Client) Dial() (net.Conn, error) {
	bidiConn := c.streamEngine.CreateConn(true, false)
	err := bidiConn.Start("GET", c.url, nil, 0, false)
	if err != nil {
		return nil, err
	}
	return bidiConn, nil
}

func (c *Client) Close() error {
	c.engine.Shutdown()
	c.engine.Destroy()
	return nil
}
