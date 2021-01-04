package main

/*
	主程序：

	编译：
	go build main.go

	运行：
	./xcli
*/

import (
	"xchain/client"
	"xchain/http"

	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	me      *client.User
	rootCmd = &cobra.Command{
		Use:   "xcli",
		Short: "xchainge client",
		Long:  "xcli is a client tool for xchainge",
	}
	initCmd = &cobra.Command{ // 生成key
		Use:   "init",
		Short: "make user key",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("home")
			// 如果 path不存在，会创建密钥文件
			_, err := client.GenUserKey(path)
			if err!=nil {
				return err
			}
			return nil
		},
	}
	dealCmd = &cobra.Command{	// 交易上链操作
		Use:   "deal",
		Short: "make a deal",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("home")
			me, err := client.GetMe(path)
			if err!=nil {
				return err
			}
			if len(args) == 0 {
				return errors.New("need more parameters")
			}
			data := args[2]
			return me.Deal(data)
		},
	}
	authRequestCmd = &cobra.Command{	// 上链操作，请求授权
		Use:   "authReq",
		Short: "Request authorization",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("home")
			me, err := client.GetMe(path)
			if err!=nil {
				return err
			}
			if len(args) < 2 {
				return errors.New("need more parameters")
			}
			fromUserId := args[0]
			dealId := args[1] // 请求授权的 dealID
			return me.AuthRequest(fromUserId, dealId)
		},
	}
	authResponseCmd = &cobra.Command{	// 上链操作，响应授权
		Use:   "authResp",
		Short: "Respond to authorization",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("home")
			me, err := client.GetMe(path)
			if err!=nil {
				return err
			}
			if len(args) == 0 {
				return errors.New("need more parameters")
			}
			authId := args[0] // 响应授权的 authID
			return me.AuthResponse(authId)
		},
	}

	queryDealCmd = &cobra.Command{	// 查询 用户 交易历史
		Use:   "queryDeal",
		Short: "query deals' history of user",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("home")
			me, err := client.GetMe(path)
			if err!=nil {
				return err
			}
			respBytes, err := me.Query("deal", "_")
			if err==nil {
				fmt.Printf("Deal ==> %s\n", respBytes)
			}
			return err
		},
	}

	queryAuthCmd = &cobra.Command{	// 查询 请求授权 历史
		Use:   "queryAuth",
		Short: "query requests of authorization",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("home")
			me, err := client.GetMe(path)
			if err!=nil {
				return err
			}
			respBytes, err := me.Query("auth", "_")
			if err==nil {
				fmt.Printf("Auth ==> %s\n", respBytes)
			}
			return err
		},
	}

	queryTxCmd = &cobra.Command{	// 查询 指定交易
		Use:   "queryTx",
		Short: "query deals by DealID",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("home")
			me, err := client.GetMe(path)
			if err!=nil {
				return err
			}
			if len(args) < 2 {
				return errors.New("need more parameters")
			}
			userId := args[0]
			realId := args[1]
			respBytes, err := me.QueryTx(userId, realId)
			if err==nil {
				fmt.Printf("Tx ==> %s\n", respBytes)
			}
			return err
		},
	}

	queryRawCmd = &cobra.Command{	// 查询 指定raw block
		Use:   "queryRaw",
		Short: "query block raw data by DealID",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("home")
			me, err := client.GetMe(path)
			if err!=nil {
				return err
			}
			if len(args) < 2 {
				return errors.New("need more parameters")
			}
			userId := args[0]
			realId := args[1]
			respBytes, err := me.QueryRawBlock(userId, realId)
			if err==nil {
				fmt.Printf("Raw ==> %s\n", respBytes)
			}
			return err
		},
	}

	httpCmd = &cobra.Command{	// 启动http服务
		Use:   "http <port> <user_path>",
		Short: "start http service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("need port number")
			} 
			port := args[0]
			userPath := args[1]
			http.RunServer(port, userPath)
			// 不会返回
			return nil
		},
	}

)

func init() {
	initCmd.Flags().StringP("home", "", "", "密钥文件路径");
	dealCmd.Flags().StringP("home", "", "", "密钥文件路径");
	authRequestCmd.Flags().StringP("home", "", "", "密钥文件路径")
	authResponseCmd.Flags().StringP("home", "", "", "密钥文件路径")
	queryDealCmd.Flags().StringP("home", "", "", "密钥文件路径")
	queryAuthCmd.Flags().StringP("home", "", "", "密钥文件路径")
	queryTxCmd.Flags().StringP("home", "", "", "密钥文件路径")
	queryRawCmd.Flags().StringP("home", "", "", "密钥文件路径")

	initCmd.MarkFlagRequired("home")
	dealCmd.MarkFlagRequired("home")
	authRequestCmd.MarkFlagRequired("home")
	authResponseCmd.MarkFlagRequired("home")
	queryDealCmd.MarkFlagRequired("home")
	queryAuthCmd.MarkFlagRequired("home")
	queryTxCmd.MarkFlagRequired("home")
	queryRawCmd.MarkFlagRequired("home")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(dealCmd)
	rootCmd.AddCommand(authRequestCmd)
	rootCmd.AddCommand(authResponseCmd)
	rootCmd.AddCommand(queryDealCmd)
	rootCmd.AddCommand(queryAuthCmd)
	rootCmd.AddCommand(queryTxCmd)
	rootCmd.AddCommand(queryRawCmd)
	rootCmd.AddCommand(httpCmd)

}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
