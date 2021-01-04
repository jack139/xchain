package chain

/*
	交易检查
*/



import (
	"xchain/types"

	"fmt"
	tmtypes "github.com/tendermint/tendermint/abci/types"
)


// 检查交易
func (app *App) isValidDeal(deal *types.Deal) error {
	m := *deal

	// 检查参数
	if len(m.UserID)==0 || len(m.Data)==0 { 
		return fmt.Errorf("bad parameters") // 参数问题
	}

	return nil
}

// 检查授权
func (app *App) isValidAuth(auth *types.Auth) error {
	m := *auth

	// 检查参数
	if len(m.FromUserID)==0 || len(m.DealID)==0 || len(m.ToUserID)==0 { 
		return fmt.Errorf("bad parameters") // 参数问题
	}

	switch m.Action {
	case 0x04, 0x05:
		// 检查 dealID 是否存在
		fromUserIdStr, _ := cdc.MarshalJSON(m.FromUserID)
		respTx := queryTx(app, fromUserIdStr, []byte(m.DealID.String()))
		if respTx==nil {
			return fmt.Errorf("DealID not exist.")	
		}
	default:
		return fmt.Errorf("weird auth command")
	}

	return nil
}

/*
	检查交易
*/
func (app *App) CheckTx(req tmtypes.RequestCheckTx) (rsp tmtypes.ResponseCheckTx) {
	app.logger.Info("CheckTx()", "para", req.Tx)

	var err error
	var tx types.Transx

	err = cdc.UnmarshalJSON(req.Tx, &tx)
	if err != nil {
		rsp.Code = 1
		rsp.Log = "error occured in decoding when CheckTx"
		return
	}
	if !tx.Verify() {
		rsp.Code = 2
		rsp.Log = "CheckTx failed"
		return
	}

	// 业务校验
	deal, ok := tx.Payload.(*types.Deal)
	if ok {
		err = app.isValidDeal(deal) // 检查 交易 合法性
	} else {
		auth, ok := tx.Payload.(*types.Auth)
		if ok {
			err = app.isValidAuth(auth)  // 检查 授权 合法性
		} else {
			rsp.Code = 3
			rsp.Log = "CheckTx unknown type"
			app.logger.Info("CheckTx() fail", "unknown type", rsp.Log)
			return				
		}
	}

	if err!=nil {
		rsp.Log = err.Error()
		rsp.Code = 4
		app.logger.Info("CheckTx() fail", "reason", rsp.Log)
	} else {
		rsp.GasWanted = 1
	}

	return 
}


