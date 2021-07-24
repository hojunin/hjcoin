package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/hojunin/hjcoin/db"
	"github.com/hojunin/hjcoin/utils"
)

type Block struct{
	Data string `json:"data"` 
	Hash string `json:"hash"`
	PrevHash string `json:"prev_hash,omitempty"`
	Height int `json:"height"`
}
// BlockChain에 블록을 저장한다. 
func (b *Block) persist()  {
	db.SaveBlock(b.Hash, utils.Tobytes(b))
}

func (b *Block) restore(data []byte)  {
	utils.FromBytes(b, data)
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

func createBlock(data string, prevHash string, height int) *Block{
	block:=&Block{
		Data: data,
		Hash: "",
		PrevHash:prevHash,
		Height: height,
	}
	payload := block.Data + block.PrevHash + fmt.Sprint(block.Height)
	block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	block.persist()

	return block
}