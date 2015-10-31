package BitCoSi

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/dedis/cothority/lib/hashid"
	"github.com/dedis/cothority/lib/proof"
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

func calculate_root(transactions TransactionList) (res string) {
	var hashes []hashid.HashId

	for _, t := range transactions.Txs {
		temp, _ := hex.DecodeString(t.Hash)
		hashes = append(hashes, temp)
	}
	out, _ := proof.ProofTree(sha256.New, hashes)
	res = hex.EncodeToString(out)
	return
}

func hash(h Header) (res string) {
	//change it to be more portable
	data := fmt.Sprintf("%v", h)
	sha := sha256.New()
	sha.Write([]byte(data))
	hash := sha.Sum(nil)
	res = hex.EncodeToString(hash)
	return

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
