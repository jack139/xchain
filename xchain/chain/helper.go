package chain

/*
	cleveldb 相关操作 及 一些辅助函数
*/


import (
	//"fmt"
	"encoding/json"
	"strconv"
	"path"
	"github.com/tendermint/tendermint/types"
	rpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
	rpc "github.com/tendermint/tendermint/rpc/core"
	dbm "github.com/tendermint/tm-db"
)



// 区块链表前缀
func blockPrefixKey(linkType string, height int64) []byte {
	/* 要按类型区分，否则同一区块的不同链表的链接信息会互相覆盖 */
	t := append(blockLinkPrefixKey, []byte(linkType+":")...)
	return append(t, Int64ToByteArray(height)...)
}

func userPrefixKey(user []byte) []byte {
	return append(userLinkPrefixKey, user...)
}


// 初始化/链接db
func InitDB(rootDir string) dbm.DB {
	// 生成数据文件路径, 放在 --home 目录下的 data 下
	dbDir := path.Join(rootDir, "data")
	//fmt.Println("xchain.db path: ", dbDir)

	// 初始化数据库
	db, err := dbm.NewCLevelDB("xchain", dbDir)  
	if err != nil {
		panic(err)
	}

	return db
}

// 从db转入应用状态
func loadState(db dbm.DB) State {
	var state State
	state.db = db
	stateBytes, err := db.Get(stateKey)
	if err != nil {
		panic(err)
	}
	if len(stateBytes) == 0 {
		return state
	}
	err = json.Unmarshal(stateBytes, &state)
	if err != nil {
		panic(err)
	}
	return state
}

// 保存应用状态
func saveState(state State) {
	stateBytes, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	state.db.Set(stateKey, stateBytes)
}


// 获取数据: 未找到返回 nil
func FindKey(db dbm.DB, key []byte) []byte {
	value2, err := db.Get(key)
	if err != nil {
		panic(err)
	}

	return value2
}


// 添加key 成功返回 nil
func AddKV(db dbm.DB, key []byte, value []byte) error {
	if value==nil {
		value = []byte("")  // db.Set 传入的值不允许是 nil
	}
	err := db.Set(key, value)
	if err != nil {
		panic(err)
	}

	return nil
}


/*
	// int64 <---> []byte 

	i := int64(-123456789)

	fmt.Println(i)

	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))

	fmt.Println(b)

	i = int64(binary.LittleEndian.Uint64(b))
	fmt.Println(i)
*/


/*
	// string --> int64
	int64, err := strconv.ParseInt(string, 10, 64)

	// int64 --> string
	string:=strconv.FormatInt(int64,10)
*/
func Int64ToByteArray(a int64) []byte {
	return []byte(strconv.FormatInt(a,10))
}

func ByteArrayToInt64(b []byte) int64 {
	a, err := strconv.ParseInt(string(b), 10, 64)
	if err!=nil {
		panic(err)
	}
	return a
}


// 获取指定高度的区块内容
func GetBlock(height int64) *types.Block{
	var ctx rpctypes.Context

	// func Block(ctx *rpctypes.Context, heightPtr *int64) (*ctypes.ResultBlock, error) 
	re, err := rpc.Block(&ctx, &height)
	if err!=nil {
		panic(err)
	}

	return re.Block		
}


// 添加到链表
func AddToLink(db dbm.DB, linkType string, keyContent []byte, height int64) {
	var linkKey []byte

	// 生成资产key
	switch linkType {
	case "user":
		linkKey = userPrefixKey(keyContent)
	default:
		return
	}
	
	// 值
	linkValue := Int64ToByteArray(height)

	// 生成blcok链表key
	blockLinkKey := blockPrefixKey(linkType, height)		
	blockLinkValue := FindKey(db, linkKey)

	// 添加到 db
	AddKV(db, linkKey, linkValue) 
	AddKV(db, blockLinkKey, blockLinkValue) 
}