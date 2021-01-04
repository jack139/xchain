package http

import (
	"log"
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
	err = me.Deal(data)
	if err != nil {
		respError(ctx, 9004, err.Error())
		return
	}

	// 正常 返回空
	resp := map[string] interface{} {
		"data" : nil,
	}
	respJson(ctx, &resp)
}
