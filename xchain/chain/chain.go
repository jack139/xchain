package chain

/*
	区块链主要定义
*/



import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/version"
)


func NewApp(rootDir string) *App {
	state := loadState(InitDB(rootDir))
	return &App{
		state:  state,
		logger: log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
	}
}

func (app *App) SetLogger(l log.Logger) {
	app.logger = l
}

func (app *App) Info(req types.RequestInfo) (resInfo types.ResponseInfo) {
	return types.ResponseInfo{
		Data:             fmt.Sprintf("{\"size\":%v}", app.state.Size),
		Version:          version.ABCIVersion,
		AppVersion:       ProtocolVersion,
		LastBlockHeight:  app.state.Height,
		LastBlockAppHash: app.state.AppHash,
	}
}

func (app *App) Commit() (rsp types.ResponseCommit) {
	app.logger.Info("Commit()")

	// Using a db - just return the big endian size of the db
	appHash := make([]byte, 8)
	binary.PutVarint(appHash, app.state.Size)
	app.state.AppHash = appHash
	app.state.Height++
	saveState(app.state)

	resp := types.ResponseCommit{Data: appHash}
	//if app.RetainBlocks > 0 && app.state.Height >= app.RetainBlocks {
	//	resp.RetainHeight = app.state.Height - app.RetainBlocks + 1
	//}
	return resp
}


