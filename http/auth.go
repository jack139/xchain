package http

import (
	"log"
	"github.com/valyala/fasthttp"
)


/* 请求授权 */
func authRequest(ctx *fasthttp.RequestCtx) {
	log.Println("auth_request")

	// POST 的数据
	content := ctx.PostBody()

	// 验签
	reqData, me, err := checkSign(content)
	if err!=nil {
		respError(ctx, 9000, err.Error())
		return
	}

	// 检查参数
	dealId, ok := (*reqData)["deal_id"].(string)
	if !ok {
		respError(ctx, 9001, "need deal_id")
		return
	}
	fromUserId, ok := (*reqData)["from_user_id"].(string)
	if !ok {
		respError(ctx, 9002, "need from_user_id")
		return
	}

	// 提交 授权请求
	err = me.AuthRequest(fromUserId, dealId)
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


/* 请求授权 */
func authResponse(ctx *fasthttp.RequestCtx) {
	log.Println("auth_response")

	// POST 的数据
	content := ctx.PostBody()

	// 验签
	reqData, me, err := checkSign(content)
	if err!=nil {
		respError(ctx, 9000, err.Error())
		return
	}

	// 检查参数
	authId, ok := (*reqData)["auth_id"].(string)
	if !ok {
		respError(ctx, 9001, "need auth_id")
		return
	}

	// 提交 授权响应
	err = me.AuthResponse(authId)
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
