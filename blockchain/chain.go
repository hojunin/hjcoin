package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sync"

	"github.com/hojunin/hjcoin/db"
	"github.com/hojunin/hjcoin/utils"
)

//SingleTon Patturn
var b *blockchain
//Goroutine같은 병렬처리에 관계없이 무조건 1회만 실행되도록하는 함수
var once sync.Once

// Block의 포인터들의 Slice. 블록체인은 길어질 수 있기때문에 매번 복사할 수 없다. 그래서 주소만 저장
type blockchain struct{
	NewestHash string `json:"newestHash"`
	Height int `json:"height"`
}

func (b *blockchain) restore(data []byte)  {
	decoder:= gob.NewDecoder(bytes.NewReader(data))
	decoder.Decode(b)
}

func (b *blockchain) persist()  {
	db.SaveBlockchain(utils.Tobytes(b))
}

func (b *blockchain) AddBlock(data string)  {
	block := createBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height=block.Height
}

func Blockchain() *blockchain{
	if b ==nil{
		once.Do(func ()  {
			b= &blockchain{"", 0}
			fmt.Printf("Newest Hash : %s\nHeight:%d\n", b.NewestHash, b.Height)
			checkpoint := db.Checkpoint()
			if checkpoint==nil{
				b.AddBlock("Genesis")
			}else{
				fmt.Println("Restoring")
				b.restore(checkpoint)
			}
		})
	}
	fmt.Printf("Newest Hash : %s\nHeight:%d\n", b.NewestHash, b.Height)
	return b
}
