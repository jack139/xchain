package http

import (
	"fmt"
	"log"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"github.com/AubSs/fasthttplogger"
)


/* 入口 */
func RunServer(port string, userPath string) {
	// 装入用户appid和secret
	err := loadSecretKey(userPath)
	if err != nil {
		log.Fatal(err)
	}

	/* router */
	r := router.New()
	r.GET("/", index)
	r.POST("/api/test", doNonthing)
	r.POST("/api/deal", deal)
	r.POST("/api/auth_request", authRequest)
	r.POST("/api/auth_response", authResponse)
	r.POST("/api/query_deals", queryDeals)
	r.POST("/api/query_auths", queryAuths)
	r.POST("/api/query_block", queryBlock)
	r.POST("/api/query_raw_block", queryRawBlock)

	fmt.Printf("start HTTP server at 0.0.0.0:%s\n", port)

	/* 启动server */
	s := &fasthttp.Server{
		Handler: fasthttplogger.Combined(r.Handler),
		Name: "FastHttpLogger",
	}
	log.Fatal(s.ListenAndServe(":"+port))
}


/* 根返回 */
func index(ctx *fasthttp.RequestCtx) {
	log.Printf("%v", ctx.RemoteAddr())
	ctx.WriteString("Hello world.")
}

