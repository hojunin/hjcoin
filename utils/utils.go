package utils

import (
	"bytes"
	"encoding/gob"
	"log"
)

func HandleErr(err error)  {
	if err != nil{
		log.Panic((err))
	}
}

// 어떤 데이터타입이든 받아서 []byte 타입으로 바꿔주는 함수
func Tobytes(i interface{}) []byte  {
	var aBuffer bytes.Buffer
	encoder := gob.NewEncoder(&aBuffer)	
	HandleErr(encoder.Encode(i))
	return aBuffer.Bytes()	
}

func FromBytes(i interface{}, data []byte)  {
	encoder := gob.NewDecoder(bytes.NewReader(data))
	HandleErr(encoder.Decode(i))
}