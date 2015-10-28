package main

import (
	"errors"
	"github.com/LefKok/Coin/BitCoSi"
	"github.com/LefKok/Coin/blkparser"
	"log"
	"net"
)

type Node struct {
	IP               net.IP
	PublicKey        string
	Last_Block       string
	transaction_pool []blkparser.Tx
}

func main() {
	Current := new(Node)
	Magic := [4]byte{0xF9, 0xBE, 0xB4, 0xD9}
	Current.IP = net.IPv4(0, 1, 2, 3)
	Current.PublicKey = "my_cool_key"
	Current.Last_Block = "0"
	Parser, _ := BitCoSi.NewParser("/home/lefteris/hi/blocks", Magic)
	Current.transaction_pool = Parser.Parse(0, 2)
	var err error
	var trblock BitCoSi.TrBlock

	for len(Current.transaction_pool) > 0 {
		trblock, err = getblock(Current, 2)
		if err != nil {
			log.Println(err)
		} else {
			trblock.Print()
		}
	}

	trblock, err = getblock(Current, 2)
	if err != nil {
		log.Println(err)
	} else {
		trblock.Print()
	}

	trblock, err = getblock(Current, 2)
	if err != nil {
		log.Println(err)
	} else {
		trblock.Print()
	}

}

func getblock(l *Node, n int) (_ BitCoSi.TrBlock, _ error) {
	if len(l.transaction_pool) > 0 {

		trlist := BitCoSi.NewTransactionList(l.transaction_pool, n)
		header := BitCoSi.NewHeader(trlist, l.Last_Block, l.IP, l.PublicKey)
		trblock := BitCoSi.NewTrBlock(trlist, header)
		l.transaction_pool = l.transaction_pool[trblock.TransactionList.TxCnt:]
		l.Last_Block = trblock.HeaderHash
		return trblock, nil
	} else {
		return *new(BitCoSi.TrBlock), errors.New("no transaction available")
	}

}
