package main

import (
	"errors"
	"github.com/LefKok/Coin/BitCoSi"
	"github.com/LefKok/Coin/blkparser"
	"github.com/dedis/cothority/lib/app"
	"github.com/dedis/cothority/lib/coconet"
	"github.com/dedis/cothority/lib/conode"
	dbg "github.com/dedis/cothority/lib/debug_lvl"
	"github.com/dedis/cothority/proto/sign"
	"github.com/dedis/crypto/abstract"
	"log"
	"net"
)

var suite abstract.Suite

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
	Current.transaction_pool = Parser.Parse(1000, 1500)
	var err error
	var trblock BitCoSi.TrBlock

	for len(Current.transaction_pool) > 0 {
		trblock, err = getblock(Current, 15)
		suite = app.GetSuite("25519")
		sig := CoSi(trblock.HeaderHash)
		dbg.Lvlf1("Signature is: %+v", sig)
		if err != nil {
			log.Println(err)
		} else {
			trblock.Print()
		}
	}

}

func CoSi(myHash string) sign.SignatureBroadcastMessage {
	server := "profeda.org:2001"
	//server := "localhost:2011"

	dbg.Lvl2("Connecting to", server)
	conn := coconet.NewTCPConn(server)
	err := conn.Connect()
	if err != nil {
		dbg.Fatal("Error when getting the connection to the host:", err)
	}
	dbg.Lvl1("Connected to ", server)
	msg := &conode.TimeStampMessage{
		Type:  conode.StampRequestType,
		ReqNo: 0,
		Sreq:  &conode.StampRequest{Val: []byte(myHash)}}

	err = conn.PutData(msg)
	if err != nil {
		dbg.Fatal("Couldn't send hash-message to server: ", err)
	}
	dbg.Lvl1("Sent signature request")
	// Wait for the signed message
	tsm := &conode.TimeStampMessage{}
	tsm.Srep = &conode.StampReply{}
	tsm.Srep.SuiteStr = suite.String()
	err = conn.GetData(tsm)
	if err != nil {
		dbg.Fatal("Error while receiving signature:", err)
	}
	dbg.Lvl1("Got signature response")

	// Asking to close the connection
	err = conn.PutData(&conode.TimeStampMessage{
		ReqNo: 1,
		Type:  conode.StampClose,
	})
	conn.Close()
	dbg.Lvl2("Connection closed with server")
	// Verify if what we received is correct
	if !conode.VerifySignature(suite, tsm.Srep, tsm.Srep.SigBroad.X0_hat, []byte(myHash)) {
		dbg.Fatal("Verification of signature failed")
	} else {
		dbg.Lvl1("Verification OK")
	}

	return tsm.Srep.SigBroad
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
