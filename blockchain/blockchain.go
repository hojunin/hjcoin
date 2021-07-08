package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type block struct{
	data string
	hash string
	prevHash string
}

// Block의 포인터들의 Slice. 블록체인은 길어질 수 있기때문에 매번 복사할 수 없다. 그래서 주소만 저장
type blockchain struct{
	blocks []*block
}

func (b *block) calculateHash(){
	hash := sha256.Sum256([]byte(b.data+b.prevHash))
	b.hash=fmt.Sprintf("%x", hash)
}

func getLastHash() string {
	totalBlockLen := len(GetBlockchain().blocks)
	if totalBlockLen ==0{
		return ""
	}
	return GetBlockchain().blocks[totalBlockLen-1].hash
}

func createBlock(data string) *block {
	newBlock := block{data, "", getLastHash()}
	newBlock.calculateHash()

	return &newBlock
}

//SingleTon Patturn
var b *blockchain
//Goroutine같은 병렬처리에 관계없이 무조건 1회만 실행되도록하는 함수
var once sync.Once

func GetBlockchain() *blockchain{
	if b ==nil{
		once.Do(func ()  {
			b= &blockchain{}
		})
	}

	return b
}