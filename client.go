package main

import (
	"github.com/dedis/cothority/lib/app"
	"github.com/dedis/cothority/lib/bitcosi"
	"github.com/dedis/cothority/lib/bitcosi/blkparser"
	"github.com/dedis/cothority/lib/coconet"
	"github.com/dedis/cothority/lib/dbg"
	"github.com/dedis/crypto/abstract"
	"net"
	"time"
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
	server := "localhost:2011"
	//	suite = app.GetSuite("25519")

	dbg.Lvl2("Connecting to", server)
	conn := coconet.NewTCPConn(server)
	err := conn.Connect()
	if err != nil {
		dbg.Fatal("Error when getting the connection to the host:", err)
	}
	dbg.Lvl1("Connected to ", server)
	go wait_for_blocks()

	for i := 0; i < 10; i++ {

		Current.transaction_pool = Parser.Parse(i, 10+i)
		for len(Current.transaction_pool) > 0 {
			msg := &BitCoSi.BitCoSiMessage{
				Type:  BitCoSi.TransactionAnnouncmentType,
				ReqNo: 0,
				Treq:  &BitCoSi.TransactionAnnouncment{Val: Current.transaction_pool[0]}}

			err = conn.PutData(msg)
			Current.transaction_pool = Current.transaction_pool[1:]
			if err != nil {
				dbg.Fatal("Couldn't send hash-message to server: ", err)
			}
			time.Sleep(10 * time.Millisecond)

		}
	}

	wait_for_Key_blocks()
	time.Sleep(900000 * time.Millisecond)

	// Asking to close the connection
	err = conn.PutData(&BitCoSi.BitCoSiMessage{
		ReqNo: 1,
		Type:  BitCoSi.BitCoSiClose,
	})

	conn.Close()
	dbg.Lvl2("Connection closed with server")

}

func wait_for_blocks() {

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
		//dbg.Lvlf1("Got signature response %+v", tsm.Brep)

		T := new(BitCoSi.TrBlock)
		T.Block = tsm.Brep.Block
		T.Print()
		dbg.Lvlf1("Response %v ", tsm.Brep.Response)
	}
	// Asking to close the connection
	err = conn.PutData(&BitCoSi.BitCoSiMessage{
		ReqNo: 1,
		Type:  BitCoSi.BitCoSiClose,
	})

	conn.Close()

}

func wait_for_Key_blocks() {

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
			Type:  BitCoSi.KeyBlockRequestType,
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
		//dbg.Lvlf1("Got signature response %+v", tsm.Brep)
		K := new(BitCoSi.KeyBlock)
		K.Block = tsm.Brep.Block
		K.Print()
		dbg.Lvlf1("Response %v ", tsm.Brep.Response)
	}
	// Asking to close the connection
	err = conn.PutData(&BitCoSi.BitCoSiMessage{
		ReqNo: 1,
		Type:  BitCoSi.BitCoSiClose,
	})

	conn.Close()

}
