package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type block struct{
	Data string
	Hash string
	PrevHash string
}

// Block의 포인터들의 Slice. 블록체인은 길어질 수 있기때문에 매번 복사할 수 없다. 그래서 주소만 저장
type blockchain struct{
	blocks []*block
}

func (b *block) calculateHash(){
	hash := sha256.Sum256([]byte(b.Data+b.PrevHash))
	b.Hash=fmt.Sprintf("%x", hash)
}

func getLastHash() string {
	totalBlockLen := len(GetBlockchain().blocks)
	if totalBlockLen ==0{
		return ""
	}
	return GetBlockchain().blocks[totalBlockLen-1].Hash
}

func createBlock(data string) *block {
	newBlock := block{data, "", getLastHash()}
	newBlock.calculateHash()

	return &newBlock
}

func (b *blockchain) AddBlock(data string)  {
	b.blocks = append(b.blocks, createBlock(data))
}

//SingleTon Patturn
var b *blockchain
//Goroutine같은 병렬처리에 관계없이 무조건 1회만 실행되도록하는 함수
var once sync.Once

func GetBlockchain() *blockchain{
	if b ==nil{
		once.Do(func ()  {
			b= &blockchain{}
			b.AddBlock("Genesis")
		})
	}

	return b
}

func (b *blockchain) AllBlock() []*block  {
	return b.blocks
}