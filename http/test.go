package http

import (
	"fmt"
	"log"
	"github.com/valyala/fasthttp"
)


/* 空接口, 只进行签名校验 */
func doNonthing(ctx *fasthttp.RequestCtx) {
	log.Println("doNonthing")

	// POST 的数据
	content := ctx.PostBody()

	// 验签
	data, _, err := checkSign(content)
	if err!=nil {
		respError(ctx, 9000, err.Error())
		return
	}
	fmt.Printf("%v\n", *data)

	respJson(ctx, data)
}
