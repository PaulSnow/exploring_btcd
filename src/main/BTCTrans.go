// Copyright (c) 2014 Conformal Systems LLC.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
   "time"
   	"github.com/conformal/btcwire"
   	"github.com/conformal/btcutil"
)


var LastHash  []byte
var SLastHash string

//  We are going to compute hash every minute.
//
func computeHashes() {
	for {
		LastHash = btcwire.DoubleSha256([]byte(time.Now().String())) 
		SLastHash = btcutil.Base58Encode(LastHash)
		time.Sleep(time.Minute)
	}
}

func main() {
	go computeHashes()
    three() 
}