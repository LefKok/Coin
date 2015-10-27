package main

import (
	"github.com/LefKok/Coin"
	"log"
)

func main() {

	Magic := [4]byte{0xF9, 0xBE, 0xB4, 0xD9}

	Parser, _ := Coin.NewParser("/home/lefteris/hi/blocks", Magic)

	transaction := Parser.parse(0, 10)

	for _, tx := range transactions {

		log.Printf("TxId: %v", tx.Hash)
		//log.Printf("Tx Size: %v", tx.Size)
		//log.Printf("Tx Lock time: %v", tx.LockTime)
		//log.Printf("Tx Version: %v", tx.Version)

		log.Println("TxIns:")
		if tx.TxInCnt == 1 && tx.TxIns[0].InputVout == 4294967295 {
			log.Printf("TxIn coinbase, newly generated coins")
		} else {
			for txin_index, txin := range tx.TxIns {
				log.Printf("TxIn index: %v", txin_index)
				log.Printf("TxIn Input_Hash: %v", txin.InputHash)
				log.Printf("TxIn Input_Index: %v", txin.InputVout)

				//		log.Printf("TxIn ScriptSig: %v", hex.EncodeToString(txin.ScriptSig))
				//	log.Printf("TxIn Sequence: %v", txin.Sequence)
			}
		}

		log.Println("TxOuts:")

		for txo_index, txout := range tx.TxOuts {
			log.Printf("TxOut index: %v", txo_index)
			log.Printf("TxOut value: %v", txout.Value)
			//	log.Printf("TxOut script: %s", hex.EncodeToString(txout.Pkscript))
			txout_addr := txout.Addr
			if txout_addr != "" {
				log.Printf("TxOut address: %v", txout_addr)
			} else {
				log.Printf("TxOut address: can't decode address")
			}
		}
	}
}
