package blockchain

import (
	"errors"
	"strings"
	"time"

	"github.com/hojunin/hjcoin/db"
	"github.com/hojunin/hjcoin/utils"
)

// hash값이 difficulty개의 0으로 시작하는 것을 찾는다.
type Block struct{
	Hash string `json:"hash"`
	PrevHash string `json:"prev_hash,omitempty"`
	Height int `json:"height"`
	Difficulty int `json:"difficulty"`
	Nonce int `json:"nonce"`
	Timestamp int `json:"timestamp"`
	Transactions []*Tx `json:"transactions"`
}
// BlockChain에 블록을 저장한다. 
func (b *Block) persist()  {
	db.SaveBlock(b.Hash, utils.Tobytes(b))
}

func (b *Block) restore(data []byte)  {
	utils.FromBytes(b, data)
}

func (b *Block) mine(){
	target := strings.Repeat("0", b.Difficulty)
	for{
		b.Timestamp = int(time.Now().Unix())
		hash := utils.Hash(b)
		if strings.HasPrefix(hash, target){
			b.Hash =hash
			break
		}else{
			b.Nonce++
		}
	}
}

var ErrNotFound = errors.New("Block Not Found")

func FindBlock(hash string) (*Block, error) {
	blockBytes:=db.Block(hash)
	if blockBytes==nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockBytes)
	return block,nil
}

func createBlock(prevHash string, height int) *Block{
	block:=&Block{
		Hash: "",
		PrevHash:prevHash,
		Height: height,
		Difficulty: Blockchain().difficulty(),
		Nonce: 0,
		Transactions: []*Tx{makeCoinbaseTx("hj")},
	}
	block.mine()
	block.persist()

	return block
}