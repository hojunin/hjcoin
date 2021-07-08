package main

import (
	"crypto/sha256"
	"fmt"
)

type block struct {
	data string
	hash string
	prevHash string
}

type blockchain struct {
	blocks []block
}

// 마지막 블록의 해시를 리턴합니다.
func (b *blockchain) getLastHash() string{
	if len(b.blocks) > 0 {
		return b.blocks[len(b.blocks)-1].hash
	}
	return ""
}

func (b *blockchain) addBlock(data string)  {
	newBlock:=block{data, "", b.getLastHash()}
	hash:=sha256.Sum256([]byte(newBlock.data+newBlock.hash))
	newBlock.hash = fmt.Sprintf("%x",hash)
	b.blocks=append(b.blocks, newBlock)
}

func (b blockchain) listBlocks()  {
	for _, block := range b.blocks{
		fmt.Printf("Data : %s\n", block.data)
		fmt.Printf("Hash : %s\n", block.hash)
		fmt.Printf("Prev Hash : %s\n", block.prevHash)
	}
}

func main()  {
	chain := blockchain{}
	chain.addBlock("First")
	chain.addBlock("Second")
	chain.addBlock("Third")
	chain.listBlocks()


}