package main

import (
	"github.com/conformal/btcrpcclient"
	"github.com/conformal/btcutil"
	"github.com/conformal/btcwire"
	"code.google.com/p/go.net/websocket"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"
)

func handlers() *btcrpcclient.NotificationHandlers {
    // Only override the handlers for notifications you care about.
	// Also note most of these handlers will only be called if you register
	// for notifications.  See the documentation of the btcrpcclient
	// NotificationHandlers type for more details about each handler.
	
   ntfnHandlers := btcrpcclient.NotificationHandlers{
		OnBlockConnected: func(hash *btcwire.ShaHash, height int32) {
			log.Printf("Block connected: %v (%d)", hash, height)
		},
		OnBlockDisconnected: func(hash *btcwire.ShaHash, height int32) {
			log.Printf("Block disconnected: %v (%d)", hash, height)
		},
	}
	return &ntfnHandlers
}

func shutdown(client *btcrpcclient.Client) {
	// For this example gracefully shutdown the client after 10 seconds.
	// Ordinarily when to shutdown the client is highly application
	// specific.
	log.Println("Client shutdown in 2 seconds...")
	time.AfterFunc(time.Second*2, func() {
		/* =============> */ log.Println("Going down...")
		client.Shutdown()
	})
	defer /* =============> */ log.Println("Shutdown done!")
	// Wait until the client either shuts down gracefully (or the user
	// terminates the process with Ctrl+C).
	client.WaitForShutdown()
}


func one () {
   
	// Connect to local btcd RPC server using websockets.
	btcdHomeDir := btcutil.AppDataDir("btcd", false)
	certs, err := ioutil.ReadFile(filepath.Join(btcdHomeDir, "rpc.cert"))
	if err != nil {
		log.Print(err)
		return
	}
	
	connCfg := &btcrpcclient.ConnConfig{
		Host:         "localhost:18334",
		Endpoint:     "ws",
		User:         "testuser",
		Pass:         "notarychain",
		Certificates: certs,
	}
	

	client, err := btcrpcclient.New(connCfg, handlers())
	if err != nil {
		log.Print(err)
		return
	}
	
	defer shutdown(client);
	
	origin := "http://localhost/"
	url := "ws://localhost:18332/frontend"
	_ , err = websocket.Dial(url, "", origin)
    if err != nil {
       log.Print(err)
       return
    }

	// Register for block connect and disconnect notifications.
	if err := client.NotifyBlocks(); err != nil {
		/* =============> */ log.Fatal(err)
	}
	log.Println("NotifyBlocks: Registration Complete")

	// Get the current block count.
	blockCount, err := client.GetBlockCount()
	if err != nil {
		log.Fatal(err)
	}	
	log.Printf("Block count: %d", blockCount)
}




