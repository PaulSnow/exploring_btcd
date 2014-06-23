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
	"encoding/hex"
	"strings"
)

    var client *btcrpcclient.Client
	var currentAddr *btcutil.Address
	var balance int64

	// Compute the balance for the currentAddr, and the list of its unspent
	// outputs
	func computeBalance() (cAmount btcutil.Amount, cList []btcjson.TransactionInput, err error) {
	    
	    // Get the list of unspent transaction outputs (utxos) that the
		// connected wallet has at least one private key for.
		
	    unspent, e := client.ListUnspent()
		if e != nil { err = e; return; }
	
		// This is going to be our map of addresses to all unspent outputs
		var outputs = make(map[string] []btcjson.ListUnspentResult)
		
		for _, input := range unspent{
			l , n := outputs[input.Address]		// Get the list of 
			if !n {
			    l = make([]btcjson.ListUnspentResult,1)
			    l[0]=input
				outputs[input.Address] = l
		    }else{
				outputs[input.Address] = append(l, input)
		    }
		}
		
		for index, unspentList := range outputs {
			if(strings.EqualFold(index, (*currentAddr).EncodeAddress())){
			    cAmount = btcutil.Amount(0)
				for i := range unspentList {
				 	cAmount +=  btcutil.Amount(unspentList[i].Amount*float64(100000000))
				}
			    cList   = make([]btcjson.TransactionInput,len(unspentList),len(unspentList))
			    for i,u := range unspentList {
			       v := new (btcjson.TransactionInput)
			       v.Txid = u.TxId
			       v.Vout = u.Vout
			       cList[i]= *v
			    }
			}
		} 		
		return 
	}
	
	func printBalance(){
	    // Get the list of unspent transaction outputs (utxos) that the
		// connected wallet has at least one private key for.

	    unspent, e := client.ListUnspent()
		if e != nil { return; }
		
		// This is going to be our map of addresses to all unspent outputs
		var outputs = make(map[string] []btcjson.ListUnspentResult)
		
		for _, input := range unspent{
			l , n := outputs[input.Address]		// Get the list of 
			if !n {
			    l = make([]btcjson.ListUnspentResult,1)
			    l[0]=input
				outputs[input.Address] = l
		    }else{
				outputs[input.Address] = append(l, input)
		    }
		}
		
		for index, unspentList := range outputs {
			// figure balance
		    b := btcutil.Amount(0)
			for i := range unspentList {
			 	b = b+ btcutil.Amount(unspentList[i].Amount*float64(100000000))
			}
			log.Print(index, " balance: ", b)
		} 		
	
	}

	func newBalance(account string, balance btcutil.Amount, confirmed bool){
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
     
    func recordHash(){
        log.Print("RecordHash() ", SLastHash)
	
		b0,_  := Encode(client, LastHash[:15])
		adr0  := hex.EncodeToString(b0)			
		b1,_  := Encode(client, LastHash[16:])
		adr1  := hex.EncodeToString(b1)			
			
		address   := make([]string,2,2)
		decodeAdr := make([]btcutil.Address,2,2)
		results   := make([]*btcjson.ValidateAddressResult,2,2   )
		
		address[0] = adr0  //external compressed public key
		address[1] = adr1  //external compressed public key
	
		for i := 0; i<len(address); i++ {
			var err error			
		
			decodeAdr[i],err = btcutil.DecodeAddress(address[i], activeNet.Params);
			if err != nil { 
			    log.Print("error decoding addr ",i, err); 
			    return
			}
					
		    results[i], err = client.ValidateAddress(decodeAdr[i]) 
			if err != nil { 
			
			     log.Print("error validating addr ",i,err); 
			     return;
			}		
			
		}	  

		addressSlice := append(make([]btcutil.Address,0,3),*currentAddr,decodeAdr[0],decodeAdr[1])
		log.Print(*currentAddr)
		log.Print(decodeAdr[0])
		log.Print(decodeAdr[1])
		
		multiAddr,e := client.AddMultisigAddress(1,addressSlice,"")
			
		if e !=nil { log.Print("Reported Error: ", multiAddr, e); return; }

        amount,unspent,err0 := computeBalance()
		fee,err1            := btcutil.NewAmount(.0005)
		change              := amount - 5430 - fee
		send,err2           := btcutil.NewAmount(.00005430)

	    if err0 != nil || err1 != nil || err2 != nil { log.Print("Reported Error: ", err1,err2 ); return; }
				
		log.Print("Amount at the address:  ", amount)
		log.Print("Change after the trans: ", change)
		log.Print("Amount to send:         ", send)
		log.Print("Send+Change+fee:        ", send+change+fee)
		log.Print("unspent: ",unspent)
				
        		
		adrs := make(map[btcutil.Address]btcutil.Amount)
		adrs[multiAddr]    = send
        // dest, _ := btcutil.DecodeAddress("mnyUYs1SJFQEKSLFZoGUsUTk8mZbbt37Ge", activeNet.Params);
		// adrs[dest] = send
        adrs[*currentAddr] = change			
	    rawTrans,err3 := client.CreateRawTransaction(unspent, adrs)
				
		if err3 != nil { log.Print ("raw trans create failed", err3); return; } 
				
		signedTrans,inputsSigned, err4 := client.SignRawTransaction(rawTrans)
				
        if err4 != nil { 
           log.Print ("Failed to sign transaction ", err4); 
           return; 
        }
		
		if !inputsSigned {
		    log.Print ("Inputs are not signed;  Is your wallet unlocked? "); 
		    return; 
		}
		
		txhash, err5 := client.SendRawTransaction(signedTrans,false)
	
		if err5 != nil { log.Print("Transaction submission failed", err5); return; }

		log.Print("WE HAVE DONE IT!")
		log.Print(txhash)
	} 
	 
	//
	// newBlock hashes the current LastHash into the Bitcoin Blockchain 5 minutes after
	// the previous block. (If a block is signed quicker than 5 minutes, then the second
	// block is ignored.)
	//
	func newBlock(hash *btcwire.ShaHash, height int32){
		
		log.Printf("Block connected: %v (%d)", hash, height)
		if bcc.inc() == 1 {
			defer bcc.clr()					// If we execute, then clear so we execute again 
											//   with the next block!
			time.Sleep(time.Minute*5)
			log.Printf("call to record hash...")
			recordHash();
			log.Printf("...hash recorded");										
	    }
    }
    
    

func three() {
	
	cadr, err := btcutil.DecodeAddress("mtmkwxRDToguqs6q9dX526wq8xxh8TCbGB", activeNet.Params);
	currentAddr = &cadr

	
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
		Endpoint:     "ws",
		User:         "testuser",
		Pass:         "notarychain",
		Certificates: certs,
	}
	
	client, err = btcrpcclient.New(connCfg, &ntfnHandlers)
	
	if err != nil {
		log.Fatal(err)
		return;
	}else{
        defer shutdown(client)
    }  		
    
    recordHash()
    
	log.Println("We watch...")
	
	for {
		log.Print(SLastHash)
  		time.Sleep(time.Minute)
	}
	
}