package db

import (
	"github.com/boltdb/bolt"
	"github.com/hojunin/hjcoin/utils"
)

// boltdbweb --db-name=blockchain.db로 조회
var db *bolt.DB

const (
	dbName = "blockchain.db"
	dataBucket = "data"
	blocksBucket = "blocks"
	checkpoint = "checkpoint"
)
// DB 연결 종료
func Close()  {
	DB().Close()
}
// DB 연결
func DB() *bolt.DB{
	if db == nil{
		dbPointer, err := bolt.Open(dbName,0600,nil)
		db=dbPointer
		utils.HandleErr(err)
		err = db.Update(func(t *bolt.Tx) error {
			_, err := t.CreateBucketIfNotExists([]byte(dataBucket))
			utils.HandleErr(err)
			_,err = t.CreateBucketIfNotExists([]byte(blocksBucket))
			return err
		})
		utils.HandleErr(err)
	}	
	return db
}

// DB에 Block 저장 key : Hash, data : []byte
func SaveBlock(hash string, data []byte)  {
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		err := bucket.Put([]byte(hash), data)
		return err
	})
	utils.HandleErr(err)
}

// DB에 BlockChain을 저장함
func SaveCheckPoint(data []byte)  {
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte(checkpoint), data)
		return err
	})
	utils.HandleErr(err)
}

// DB에서 특정 CheckPoint에 해당하는 데이터를 리턴
func Checkpoint() []byte {
	var data []byte
	DB().View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkpoint))
		return nil
	})
	return data
}

// DB에서 특정 Block을 리턴
func Block(hash string) []byte {
	var data []byte
	DB().View(func(t *bolt.Tx) error {
		bucket:=t.Bucket([]byte(blocksBucket))
		data = bucket.Get([]byte(hash))
		return nil
	})
	return data
}