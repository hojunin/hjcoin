package blockchain

import (
	"fmt"
	"sync"

	"github.com/hojunin/hjcoin/db"
	"github.com/hojunin/hjcoin/utils"
)

//SingleTon Patturn
var b *blockchain
//Goroutine같은 병렬처리에 관계없이 무조건 1회만 실행되도록하는 함수
var once sync.Once

const (
	defaultDifficulty int = 2
	difficultyInterval int = 5
	blockInterval int = 2
	allowRange int = 2
)

// Block의 포인터들의 Slice. 블록체인은 길어질 수 있기때문에 매번 복사할 수 없다. 그래서 주소만 저장
type blockchain struct{
	NewestHash string `json:"newestHash"`
	Height int `json:"height"`
	CurrentDifficulty int `json:"currentDifficulty"`
}

func (b *blockchain) restore(data []byte)  {
	utils.FromBytes(b, data)
}

func (b *blockchain) persist()  {
	db.SaveCheckPoint(utils.Tobytes(b))
}

func (b *blockchain) AddBlock(data string)  {
	block := createBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height=block.Height
	b.CurrentDifficulty=block.Difficulty
	b.persist()
}

func (b *blockchain) Blocks() []*Block {
	var blocks []*Block
	hashCursor := b.NewestHash
	for{
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		}else{
			break
		}
	}
	return blocks
}

func (b *blockchain) recalculateDifficulty() int {
	allBlocks := b.Blocks()
	newestBlock := allBlocks[0]
	lastCalculatedBlock := allBlocks[difficultyInterval-1]
	actualTime := (newestBlock.Timestamp/60) - (lastCalculatedBlock.Timestamp/60)
	expectedTime := difficultyInterval * blockInterval
	if actualTime <= (expectedTime - allowRange) {
		return b.CurrentDifficulty+1
	}else if actualTime >= (expectedTime + allowRange){
		return b.CurrentDifficulty -1
	}
	return b.CurrentDifficulty
}

func (b *blockchain) difficulty() int {
	if b.Height == 0 {
		return defaultDifficulty
	}else if b.Height%difficultyInterval == 0 {
		return b.recalculateDifficulty()
	}else {
		return b.CurrentDifficulty
	}
}

func Blockchain() *blockchain{
	if b ==nil{
		once.Do(func ()  {
			b= &blockchain{ Height: 0}
			checkpoint := db.Checkpoint()
			if checkpoint==nil{
				fmt.Println("Initializing...")
				b.AddBlock("Genesis")
			}else{
				fmt.Println("Restoring...")
				b.restore(checkpoint)
			}
		})
	}
	return b
}
