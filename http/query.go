package http

import (
	"log"
	"encoding/json"
	"github.com/valyala/fasthttp"
)


/* 查询交易， 只允许查询自己的 */
func queryDeals(ctx *fasthttp.RequestCtx) {
	log.Println("query_deals")

	// POST 的数据
	content := ctx.PostBody()

	// 验签
	_, me, err := checkSign(content)
	if err!=nil {
		respError(ctx, 9000, err.Error())
		return
	}

	// 检查参数 - 无需参数

	// 只查询当前用户的交易
	respBytes, err := me.Query("deal", "_")
	if err!=nil {
		respError(ctx, 9001, err.Error())
		return
	}

	// 转换成map, 生成返回数据
	var respData []map[string]interface{}

	if err := json.Unmarshal(respBytes, &respData); err != nil {
		respError(ctx, 9004, err.Error())
		return
	}

	resp := map[string] interface{} {
		"deals" : respData,
	}

	respJson(ctx, &resp)
}


/* 查询授权， 只允许查询自己的 */
func queryAuths(ctx *fasthttp.RequestCtx) {
	log.Println("query_auths")

	// POST 的数据
	content := ctx.PostBody()

	// 验签
	_, me, err := checkSign(content)
	if err!=nil {
		respError(ctx, 9000, err.Error())
		return
	}

	// 检查参数 - 无需参数

	// 只查询当前用户的交易
	respBytes, err := me.Query("auth", "_")
	if err!=nil {
		respError(ctx, 9001, err.Error())
		return
	}

	// 转换成map, 生成返回数据
	var respData []map[string]interface{}

	if err := json.Unmarshal(respBytes, &respData); err != nil {
		respError(ctx, 9004, err.Error())
		return
	}

	resp := map[string] interface{} {
		"auths" : respData,
	}

	respJson(ctx, &resp)
}



/* 指定区块查询交易 */
func queryBlock(ctx *fasthttp.RequestCtx) {
	log.Println("query_block")

	// POST 的数据
	content := ctx.PostBody()

	// 验签
	reqData, me, err := checkSign(content)
	if err!=nil {
		respError(ctx, 9000, err.Error())
		return
	}

	// 检查参数
	userId, ok := (*reqData)["user_id"].(string)
	if !ok {
		respError(ctx, 9001, "need user_id")
		return
	}
	blockId, ok := (*reqData)["block_id"].(string)
	if !ok {
		respError(ctx, 9002, "need block_id")
		return
	}

	respBytes, err := me.QueryTx(userId, blockId)
	if err!=nil {
		respError(ctx, 9003, err.Error())
		return
	}

	// 转换成map, 生成返回数据
	var respData map[string]interface{}
	if len(respBytes)>0 {
		if err := json.Unmarshal(respBytes, &respData); err != nil {
			respError(ctx, 9004, err.Error())
			return
		}
	}
	resp := map[string] interface{} {
		"blcok" : respData,
	}

	respJson(ctx, &resp)
}


/* 指定区块查询交易 */
func queryRawBlock(ctx *fasthttp.RequestCtx) {
	log.Println("query_raw_block")

	// POST 的数据
	content := ctx.PostBody()

	// 验签
	reqData, me, err := checkSign(content)
	if err!=nil {
		respError(ctx, 9000, err.Error())
		return
	}

	// 检查参数
	userId, ok := (*reqData)["user_id"].(string)
	if !ok {
		respError(ctx, 9001, "need user_id")
		return
	}
	blockId, ok := (*reqData)["block_id"].(string)
	if !ok {
		respError(ctx, 9002, "need block_id")
		return
	}

	respBytes, err := me.QueryRawBlock(userId, blockId)
	if err!=nil {
		respError(ctx, 9003, err.Error())
		return
	}

	// 转换成map, 生成返回数据
	var respData map[string]interface{}
	if len(respBytes)>0 {
		if err := json.Unmarshal(respBytes, &respData); err != nil {
			respError(ctx, 9004, err.Error())
			return
		}
	}
	resp := map[string] interface{} {
		"blcok" : respData,
	}

	respJson(ctx, &resp)
}
