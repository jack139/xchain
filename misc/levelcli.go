package main

/*
	cleveldb 查水表

	带 cleveldb 编译
	go build -tags 'cleveldb' levelcli.go

*/

import (
	//"encoding/json"
	"encoding/binary"
	"fmt"
	"time"

	dbm "github.com/tendermint/tm-db"
)


func ByteArrayToInt64(b []byte) int64 {
	return int64(binary.LittleEndian.Uint64(b))
}


func CountKeys(db dbm.DB, show int) int {
	// 循环获取
	itr, err := db.Iterator(nil, nil)
	if err != nil {
		panic(err)
	}

	count := 0
	for ; itr.Valid(); itr.Next() {
		if show==1 {
			fmt.Println(string(itr.Key()), "=", string(itr.Value()))
		}
		count += 1
	}

	return count	
}


func SearchKeys(db dbm.DB, start, end []byte) int {
	// 循环获取
	itr, err := db.Iterator(start, end)
	if err != nil {
		panic(err)
	}

	count := 0
	for ; itr.Valid(); itr.Next() {
		fmt.Println(string(itr.Key()), "=", string(itr.Value()))
		count += 1
	}

	return count
}

// 获取数据: 未找到返回 nil
func FindKey(db dbm.DB, key []byte) []byte {
	value2, err := db.Get(key)
	if err != nil {
		panic(err)
	}

	return value2
}

// 检查文件hash是否已存在
func FileHashExisted(db dbm.DB, fileHash string) bool {
	if FindKey(db, []byte(fileHash))!=nil {
		return true
	}

	return false
}


func main() {
	var db dbm.DB
	name := "xchain"
	dbDir := "n1/data"

	// 初始化数据库
	db, err := dbm.NewCLevelDB(name, dbDir)  
	if err != nil {
		panic(err)
	}

	start := time.Now()

	fmt.Println("count=", CountKeys(db, 1))

	//fmt.Println("count=", SearchKeys(db, []byte("abc|"), []byte("abc|\xff"))) // key可以包含汉字

	//fmt.Println("FindKey: ", FileHashExisted(db, "fileLink:1234"))

	fmt.Println("time elapsed: ", time.Now().Sub(start))
}


