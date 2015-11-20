package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/dedis/cothority/lib/app"
	"github.com/dedis/cothority/lib/bitcosi"
	"github.com/dedis/cothority/lib/bitcosi/blkparser"
	"github.com/dedis/cothority/lib/coconet"
	dbg "github.com/dedis/cothority/lib/debug_lvl"
	"github.com/dedis/cothority/lib/hashid"
	"github.com/dedis/cothority/lib/proof"
	"github.com/dedis/crypto/abstract"
	"net"
	"time"
)

var suite abstract.Suite

type Client struct {
	IP               net.IP
	PublicKey        string
	Last_Block       string
	transaction_pool []blkparser.Tx
	//mux              sync.Mutex
}

func main() {
	c := new(Client)
	Magic := [4]byte{0xF9, 0xBE, 0xB4, 0xD9}
	c.IP = net.IPv4(0, 1, 2, 3)
	c.PublicKey = "my_cool_key"
	c.Last_Block = "0"
	//c.mux = sync.Mutex{}
	Parser, _ := BitCoSi.NewParser("/home/lefteris/hi/blocks", Magic)
	server := "localhost:2011"

	//	suite = app.GetSuite("25519")

	dbg.Lvl2("Connecting to", server)
	conn := coconet.NewTCPConn(server)
	err := conn.Connect()
	if err != nil {
		dbg.Fatal("Error when getting the connection to the host:", err)
	}
	dbg.Lvl1("Connected to ", server)
	go c.wait_for_blocks()

	for i := 0; i < 1000; i++ {

		c.transaction_pool = Parser.Parse(i, 100+i)
		for len(c.transaction_pool) > 0 {
			msg := &BitCoSi.BitCoSiMessage{
				Type:  BitCoSi.TransactionAnnouncmentType,
				ReqNo: 0,
				Treq:  &BitCoSi.TransactionAnnouncment{Val: c.transaction_pool[0]}}

			err = conn.PutData(msg)
			c.transaction_pool = c.transaction_pool[1:]
			if err != nil {
				dbg.Fatal("Couldn't send hash-message to server: ", err)
			}
			time.Sleep(90 * time.Millisecond)

		}
	}
	// Asking to close the connection
	err = conn.PutData(&BitCoSi.BitCoSiMessage{
		ReqNo: 1,
		Type:  BitCoSi.BitCoSiClose,
	})

	conn.Close()
	dbg.Lvl2("Connection closed with server")

}

func (c *Client) wait_for_blocks() {

	server := "localhost:2011"
	suite = app.GetSuite("25519")

	dbg.Lvl2("Connecting to", server)
	conn := coconet.NewTCPConn(server)
	err := conn.Connect()
	if err != nil {
		dbg.Fatal("Error when getting the connection to the host:", err)
	}
	dbg.Lvl1("Connected to ", server)
	for i := 0; i < 1000; i++ {
		time.Sleep(1 * time.Second)
		msg := &BitCoSi.BitCoSiMessage{
			Type:  BitCoSi.BlockRequestType,
			ReqNo: 0,
		}

		err = conn.PutData(msg)
		if err != nil {
			dbg.Fatal("Couldn't send hash-message to server: ", err)
		}
		dbg.Lvl1("Sent signature request")
		// Wait for the signed message

		tsm := new(BitCoSi.BitCoSiMessage)
		tsm.Brep = &BitCoSi.BlockReply{}
		tsm.Brep.SuiteStr = suite.String()
		err = conn.GetData(tsm)
		if err != nil {
			dbg.Fatal("Error while receiving signature:", err)
		}
		dbg.Lvl1("Got signature response")
		//check block validit
		verified := c.verify_and_store(tsm.Brep.Block)
		if verified {
			tsm.Brep.Block.Print()
			dbg.Lvlf1("Signature %v", tsm.Brep.SigBroad)
			c.Last_Block = tsm.Brep.Block.HeaderHash
		} else {
			dbg.Lvlf1("Block is not valid")
		}
	}
	// Asking to close the connection
	err = conn.PutData(&BitCoSi.BitCoSiMessage{
		ReqNo: 1,
		Type:  BitCoSi.BitCoSiClose,
	})

	conn.Close()

}

func (c *Client) verify_and_store(block BitCoSi.TrBlock) bool {

	return block.Header.Parent == c.Last_Block && block.Header.MerkleRoot == c.calculate_root(block.TransactionList) && block.HeaderHash == c.hash(block.Header)

}

func (c *Client) calculate_root(transactions BitCoSi.TransactionList) (res string) {
	var hashes []hashid.HashId

	for _, t := range transactions.Txs {
		temp, _ := hex.DecodeString(t.Hash)
		hashes = append(hashes, temp)
	}
	out, _ := proof.ProofTree(sha256.New, hashes)
	res = hex.EncodeToString(out)
	return
}

func (c *Client) hash(h BitCoSi.Header) (res string) {
	//change it to be more portable
	data := fmt.Sprintf("%v", h)
	sha := sha256.New()
	sha.Write([]byte(data))
	hash := sha.Sum(nil)
	res = hex.EncodeToString(hash)
	return
}
