package main

/*
	主程序：

	编译：
	go build -tags 'cleveldb' main.go

	运行：
	./xchain init --home n1
	./xchain node --home n1 --consensus.create_empty_blocks=false
*/

import (
	"xchain/xchain/chain"

	"github.com/tendermint/tendermint/cmd/tendermint/commands"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
)

func main() {
	root := commands.RootCmd
	root.AddCommand(commands.GenNodeKeyCmd)
	root.AddCommand(commands.GenValidatorCmd)
	root.AddCommand(commands.InitFilesCmd)
	root.AddCommand(commands.ResetAllCmd)
	root.AddCommand(commands.ShowNodeIDCmd)
	root.AddCommand(commands.TestnetFilesCmd)

	nodeProvider := makeNodeProvider()
	root.AddCommand(commands.NewRunNodeCmd(nodeProvider))

	exec := cli.PrepareBaseCmd(root, "wiz", ".")
	exec.Execute()
}

func makeNodeProvider() node.Provider {
	return func(config *cfg.Config, logger log.Logger) (*node.Node, error) {
		nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
		if err != nil {
			return nil, err
		}

		// instance app
		app := chain.NewApp(config.RootDir)

		// read private validator
		pv := privval.LoadFilePV(
			config.PrivValidatorKeyFile(),
			config.PrivValidatorStateFile(),
		)

		return node.NewNode(config,
			pv, //privval.LoadOrGenFilePV(config.PrivValidatorFile()),
			nodeKey,
			proxy.NewLocalClientCreator(app),
			node.DefaultGenesisDocProviderFunc(config),
			node.DefaultDBProvider,
			node.DefaultMetricsProvider(config.Instrumentation),
			logger,
		)
	}
}
