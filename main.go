package main

import (
	"fmt"

	"github.com/sagernet/cronet-go"
)

func main() {
	engine := cronet.NewEngine()
	fmt.Printf("libcronet %v", engine.Version())
	engine.Destroy()
	// params := cronet.NewEngineParams()

	// params.setEnableHTTP2(true)
	// params.setEnableQuic(false)

	// engine.StartWithParams(params)
	// params.Destroy()

	// streamEngine = engine.StreamEngine()

	// bidiConn := streamEngine.CreateConn(true, false)
	// bidiConn.Start("CONNECT", l.url, headers, 0, false)

	// engine.Shutdown()
	// engine.Destroy()
	// inbound.Close()
}
