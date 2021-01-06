package http

import (
	"log"
	"encoding/json"
	"github.com/valyala/fasthttp"
)


/* 提交交易 */
func deal(ctx *fasthttp.RequestCtx) {
	log.Println("deal")

	// POST 的数据
	content := ctx.PostBody()

	// 验签
	reqData, me, err := checkSign(content)
	if err!=nil {
		respError(ctx, 9000, err.Error())
		return
	}

	// 检查参数
	data, ok := (*reqData)["data"].(string)
	if !ok {
		respError(ctx, 9003, "need data")
		return
	}

	// 提交交易
	respBytes, err := me.Deal(data)
	if err != nil {
		respError(ctx, 9004, err.Error())
		return
	}

	// 转换成map, 生成返回数据
	var respData map[string]interface{}

	if err := json.Unmarshal(respBytes, &respData); err != nil {
		respError(ctx, 9005, err.Error())
		return
	}

	resp := map[string] interface{} {
		"data" : respData,
	}

	respJson(ctx, &resp)
}
