package BitCoSi

import (
	"log"
	"net"
)

type TrBlock struct {
	Magic      [4]byte
	BlockSize  uint32
	HeaderHash string
	Header
	TransactionList
}

type Header struct {
	LeaderId   net.IP
	PublicKey  string
	MerkleRoot string
	Parent     string
}

func NewTrBlock(transactions TransactionList, header Header) (tr TrBlock) {
	trb := new(TrBlock)
	trb.Magic = [4]byte{0xF9, 0xBE, 0xB4, 0xD9}
	trb.HeaderHash = hash(header)
	trb.TransactionList = transactions
	trb.BlockSize = 0
	trb.Header = header
	return *trb
}

func NewHeader(transactions TransactionList, parent string, IP net.IP, key string) (hd Header) {
	hdr := new(Header)
	hdr.LeaderId = IP
	hdr.PublicKey = key
	hdr.Parent = parent
	hdr.MerkleRoot = calculate_root(transactions)
	return *hdr
}

func calculate_root(transactions TransactionList) (s string) {
	return transactions.Txs[0].Hash
}

func hash(h Header) (s string) {
	return h.Parent + "1"

}

func (trb *TrBlock) Print() {
	log.Println("Header:")
	log.Printf("Leader %v", trb.LeaderId)
	log.Printf("Pkey %v", trb.PublicKey)
	log.Printf("Parent %v", trb.Parent)
	log.Printf("Merkle %v", trb.MerkleRoot)
	trb.TransactionList.Print()

	log.Println("Rest:")
	log.Printf("Hash %v", trb.HeaderHash)

	return
}
