package chain


import (
	"xchain/types"

	tmtypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

// 一些参数
var (
	stateKey        = []byte("stateKey")
	blockLinkPrefixKey = []byte("blockLink:")
	userLinkPrefixKey = []byte("userLink:")

	ProtocolVersion uint64 = 0x1

	actionMessage = []string{
		"", 
		"Buy in", 
		"Sell", 
		"Change of ownership", 
		"Request authorization",
		"Respond to authorization",
	}

	cdc = types.AminoCdc
)

// 保存应用状态使用
type State struct {
	db      dbm.DB
	Size    int64  `json:"size"`
	Height  int64  `json:"height"`
	AppHash []byte `json:"app_hash"`
}

// 应用的结构
type App struct {
	tmtypes.BaseApplication

	state State
	RetainBlocks int64 // blocks to retain after commit (via ResponseCommit.RetainHeight)

	logger log.Logger
}
