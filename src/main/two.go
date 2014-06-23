package main

import (
	"github.com/conformal/btcrpcclient"
	"github.com/conformal/btcutil"
	"github.com/conformal/btcwire"
	//"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"log"
	"path/filepath"
	"errors"
	"time"
	"fmt"
	"encoding/hex"
) 

// Enodes up to 30 bytes into a Bitcoin public key.  
// Returns the public key.  The format is as follows:
// 1      byte   02  (per Bitcoin spec)
// 1      byte   len (Number of bytes encoded, between 1 and 63)
// len    bytes  encoded data
// 30-len bytes  random data
// fudge  byte   changed to put the value on the eliptical curve
//
func Encode (client *btcrpcclient.Client, hash []byte ) ([]byte, error) {
    length := len(hash)
    if(length==0 || length>30){
		return nil, errors.New("Encode can only handle 1 to 63 bytes")
	}
	data := btcwire.DoubleSha256(hash);
	if(length<30){
		hash = append(hash, data[:30-length]...)
    }
    b := []byte {2, byte(length)}
    b = append(b,hash...)
    b = append(b,0)
	
	for i:= 0; i< 256; i++ {
		b[len(b)-1] = byte(i)
	    adr2  := hex.EncodeToString(b)	
		_,e := btcutil.DecodeAddress(adr2, activeNet.Params);
    	if e == nil {
    	   return b, nil
    	}
    }

	log.Print("Failure")
	return b, errors.New("Couldn't fix the address")
}

//
// Faithfully extracts upto 63 bytes encoded into the given bitcoin address
func Decode (addr []byte) []byte {
    length := int(addr[1])
    data := addr[2:length+2]
    return data
}

func two () {
   
	// Only override the handlers for notifications you care about.
	// Also note most of the handlers will only be called if you register
	// for notifications.  See the documentation of the btcrpcclient
	// NotificationHandlers type for more details about each handler.
	ntfnHandlers := btcrpcclient.NotificationHandlers{
		
        OnBlockConnected: func(hash *btcwire.ShaHash, height int32) {
           log.Printf("Block connected: %v (%d)", hash, height)
        },
		
		OnBlockDisconnected: func(hash *btcwire.ShaHash, height int32) {
			log.Printf("Block disconnected: %v (%d)", hash, height)
		},
		
        OnAccountBalance: func(account string, balance btcutil.Amount, confirmed bool) {
            sconfirmed := "False"
            if confirmed { sconfirmed = "True" }
            log.Printf("Update: %s Amount %d Confirmed %s", account, balance, sconfirmed)
            const layout = "Jan 2, 2006 at 3:04pm (MST)"
		    t := time.Now()
	        fmt.Println(t.Format(layout))
        },
	}

	// Connect to local btcwallet RPC server using websockets.
	certHomeDir := btcutil.AppDataDir("btcwallet", false)
	certs, err := ioutil.ReadFile(filepath.Join(certHomeDir, "rpc.cert"))
	if err != nil {
		log.Fatal(err)
	}
	
	connCfg := &btcrpcclient.ConnConfig{
		Host:         "localhost:18332",
		Endpoint:     "frontend",
		User:         "testuser",
		Pass:         "notarychain",
		Certificates: certs,
	}
	
	client, err := btcrpcclient.New(connCfg, &ntfnHandlers)
	
	if err != nil {
		log.Fatal(err)
		return
	}else{
	    defer shutdown(client)
    }
    
    client.NotifyBlocks()

	{ // test encode
		test := []byte("13534523501")
		
		log.Println("Encoding : ", string(test) )	
		v,err := Encode(client, test)
		if err == nil {
	    	log.Println("test: "+ btcutil.Base58Encode(v))
		    code := Decode(v)
			log.Println("back: "+btcutil.Base58Encode(code)+"\n\n")
	    }else{
	    	log.Println(err)
	    }
		
		if err != nil {
			log.Fatal(err)
			return
		}
	}
		
	/*	
	// Get the list of unspent transaction outputs (utxos) that the
	// connected wallet has at least one private key for.
	unspent, err := client.ListUnspent()
	log.Printf("Num unspent outputs (utxos): %d", len(unspent))
	
	if len(unspent) > 0 {
		log.Printf("First utxo:\n%v", spew.Sdump(unspent[0]))
	}
    */
    
    log.Println("We watch...")
	
	for {
  		time.Sleep(time.Minute)
		log.Print(" ")
	}	
   
}




