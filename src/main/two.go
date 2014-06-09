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
) 

// Hashes up to 20 bytes into a Bitcoin Address.  Any more than 20 bytes, and
// Encode returns an error.  Any less than 20 bytes, and Encode pads with 0's 
// to make 20 bytes. 
func Encode (hash []byte ) ([]byte, error) {
    // Format is 1 byte for a network and address class (i.e. P2PKH vs
	// P2SH), 20 bytes for a RIPEMD160 hash, and 4 bytes of checksum.
	if(len(hash)>20){
		return nil, errors.New("Encode can only handle 20 bytes or less")
	}
	if(len(hash)<20){
		var padlen int = 20-len(hash)
		var pad []byte
	    pad = make([]byte,padlen)
		hash = append(hash, pad...)
    }
	b := make([]byte, 0, 1+20+4)
	b = append(b, 111)
	b = append(b, hash...)
	cksum := btcwire.DoubleSha256(b)[:4]
	b = append(b, cksum...)
	return b, nil
}

//
// Faithfully extracts the 20 bytes encoded into the given bitcoin address
func Decode (addr []byte) []byte {
    data := addr[1:21]
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
		v,err := Encode(test)
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




