// Copyright (c) 2014 Conformal Systems LLC.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"github.com/conformal/btcrpcclient"
	"github.com/conformal/btcutil"
	"github.com/conformal/btcwire"
	"github.com/conformal/btcjson"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"
	"sync"
)

    var client *btcrpcclient.Client
	var btcAdr string = "mttUhQD8TBw2TXNAoZhwyfCcyNzAgfxwvq" 

	func newBalance(account string, balance btcutil.Amount, confirmed bool){
		
		time.Sleep(time.Second)
		// Get the list of unspent transaction outputs (utxos) that the
		// connected wallet has at least one private key for.
		unspent, err := client.ListUnspent()
		if err != nil {
			log.Fatal(err)
		}
		
		// unspent outputs
		var outputs = make(map[string] []btcjson.ListUnspentResult)

		for _, input := range unspent {
			l , n := outputs[input.ScriptPubKey]
			if !n {
			    l = make([]btcjson.ListUnspentResult,1)
			    l[0]=input
				outputs[input.ScriptPubKey] = l
		    }else{
				outputs[input.ScriptPubKey] = append(l, input)
		    }
		}
		
		for index, unspentList := range outputs {
			// figure balance
		    b := float64(0)
			for i := range unspentList {
			 	b += unspentList[i].Amount
			}
			log.Print(index, " balance: ", b)
		} 
		
	
		sconf := "unconfirmed"
		if confirmed { sconf = "confirmed" }
        log.Printf("New %s balance for account %s: %v", sconf, account, balance)
	}
	
	type BlkCallCnt struct { 
    	m *sync.Mutex
		cnt int
	}
    
	func (bcc *BlkCallCnt) inc() int {
		bcc.m.Lock()
		defer bcc.m.Unlock()
		bcc.cnt++
		return bcc.cnt
	}
	
	func (bcc *BlkCallCnt) clr() {
		bcc.m.Lock()
		defer bcc.m.Unlock()
		bcc.cnt = 0
	}
	
    var bcc = &BlkCallCnt{
        m: &sync.Mutex{},
    }
     
	//
	// newBlock hashes the current LastHash into the Bitcoin Blockchain 5 minutes after
	// the previous block. (If a block is signed quicker than 5 minutes, then the second
	// block is ignored.)
	//
	func newBlock(hash *btcwire.ShaHash, height int32){
		
		log.Printf("Block connected: %v (%d)", hash, height)
		if bcc.inc() == 1 {
			//time.Sleep(2*time.Minute)
			log.Print("Waiting for 5 minutes...")
			log.Print("....................Register ", SLastHash)
			b2,_  := Encode(LastHash[:15])
			adr2 := btcutil.Base58Encode(b2)			
			b3,_ := Encode(LastHash[16:])
			adr3 := btcutil.Base58Encode(b3)
			realAddr := "mtmkwxRDToguqs6q9dX526wq8xxh8TCbGB"
			address1,_ := btcutil.DecodeAddress(realAddr, activeNet.Params);
			address2,_ := btcutil.DecodeAddress(adr2, activeNet.Params);
			address3,_ := btcutil.DecodeAddress(adr3, activeNet.Params);
			addressSlice := append(make([]btcutil.Address,0),address1,address2,address3)
			mSigResult, e := client.CreateMultisig(1, addressSlice)
			if e !=nil {
			   log.Print(mSigResult,"\n",e)
			   panic("Bad multisig")
			}
		log.Print("here")	
			tadr := append (Decode(btcutil.Base58Decode(adr2))[:16], Decode(btcutil.Base58Decode(adr2))...)
		log.Print("here")	
			log.Print("Decoded: ", btcutil.Base58Encode(tadr))
			bcc.clr()
		log.Print("here")	
		}
    }


func three() {

	// Only override the handlers for notifications you care about.
	// Also note most of the handlers will only be called if you register
	// for notifications.  See the documentation of the btcrpcclient
	// NotificationHandlers type for more details about each handler.
	ntfnHandlers := btcrpcclient.NotificationHandlers{
		OnAccountBalance: func(account string, balance btcutil.Amount, confirmed bool) {
		    go newBalance(account, balance, confirmed)
		},
		
		OnBlockConnected: func(hash *btcwire.ShaHash, height int32) {
            go newBlock(hash, height)
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
	
	client, err = btcrpcclient.New(connCfg, &ntfnHandlers)
	
	if err != nil {
		log.Fatal(err)
	}else{
        defer shutdown(client)
    }  		
    
    	newBlock(nil,0)
    
	log.Println("We watch...")
	
	for {
		log.Print(SLastHash)
  		time.Sleep(time.Minute)
	}
	
}