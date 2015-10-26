package main

import (
	//"encoding/hex"
	"github.com/Lefkok/blkparser"
	"log"
	"os"
)

func main() {

	f, err := os.OpenFile("Block_chain.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	Magic := [4]byte{0xF9, 0xBE, 0xB4, 0xD9}

	Chain, _ := blkparser.NewBlockchain("/home/lefteris/hi/blocks", Magic)

	first_block := 76500
	last_block := 76546

	for i := 0; i < last_block; i++ {
		raw, err := Chain.FetchNextBlock()

		if raw == nil || err != nil {
			log.Println("End of Chain")
		}

		bl, err := blkparser.NewBlock(raw[:])

		if err != nil {
			println("Block inconsistent:", err.Error())
			break
		}

		// Read block till we reach start_block
		if i < first_block {
			continue
		}

		log.Printf("Current block height: %v", i)

		// Basic block info
		log.Printf("Block hash: %v", bl.Hash)
		//log.Printf("Block version: %v", bl.Version)
		log.Printf("Block parent: %v", bl.Parent)
		log.Printf("Block merkle root: %v", bl.MerkleRoot)
		//log.Printf("Block size: %v", len(bl.Raw))

		for _, tx := range bl.Txs {
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
		log.Println()
	}
}
