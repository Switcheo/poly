/*
 * Copyright (C) 2021 The poly network Authors
 * This file is part of The poly network library.
 *
 * The poly network is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The poly network is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with the poly network.  If not, see <http://www.gnu.org/licenses/>.
 */
package main

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/polynetwork/poly/cmd"
	"github.com/polynetwork/poly/cmd/abi"
	cmdsvr "github.com/polynetwork/poly/cmd/sigsvr"
	clisvrcom "github.com/polynetwork/poly/cmd/sigsvr/common"
	"github.com/polynetwork/poly/cmd/sigsvr/store"
	"github.com/polynetwork/poly/cmd/utils"
	"github.com/polynetwork/poly/common/config"
	"github.com/polynetwork/poly/common/log"
	"github.com/urfave/cli"
)

func setupSigSvr() *cli.App {
	app := cli.NewApp()
	app.Usage = "Ontology Sig server"
	app.Action = startSigSvr
	app.Version = config.Version
	app.Copyright = "Copyright in 2021 The Poly Network Authors"
	app.Flags = []cli.Flag{
		utils.LogLevelFlag,
		utils.CliWalletDirFlag,
		//cli setting
		utils.CliAddressFlag,
		utils.CliRpcPortFlag,
		utils.CliABIPathFlag,
	}
	app.Commands = []cli.Command{
		cmdsvr.ImportWalletCommand,
	}
	app.Before = func(context *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}
	return app
}

func startSigSvr(ctx *cli.Context) {
	logLevel := ctx.GlobalInt(utils.GetFlagName(utils.LogLevelFlag))
	log.InitLog(logLevel, log.PATH, log.Stdout)

	walletDirPath := ctx.String(utils.GetFlagName(utils.CliWalletDirFlag))
	if walletDirPath == "" {
		log.Errorf("Please using --walletdir flag to specific wallet saving path")
		return
	}

	walletStore, err := store.NewWalletStore(walletDirPath)
	if err != nil {
		log.Errorf("NewWalletStore error:%s", err)
		return
	}
	clisvrcom.DefWalletStore = walletStore

	accountNum, err := walletStore.GetAccountNumber()
	if err != nil {
		log.Errorf("GetAccountNumber error:%s", err)
		return
	}
	log.Infof("Load wallet data success. Account number:%d", accountNum)

	rpcAddress := ctx.String(utils.GetFlagName(utils.CliAddressFlag))
	rpcPort := ctx.Uint(utils.GetFlagName(utils.CliRpcPortFlag))
	if rpcPort == 0 {
		log.Errorf("Please using sig server port by --%s flag", utils.GetFlagName(utils.CliRpcPortFlag))
		return
	}
	go cmdsvr.DefCliRpcSvr.Start(rpcAddress, rpcPort)

	abiPath := ctx.GlobalString(utils.GetFlagName(utils.CliABIPathFlag))
	abi.DefAbiMgr.Init(abiPath)

	log.Infof("Sig server init success")
	log.Infof("Sig server listing on: %s:%d", rpcAddress, rpcPort)

	exit := make(chan bool, 0)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		for sig := range sc {
			log.Infof("Sig server received exit signal:%v.", sig.String())
			close(exit)
			break
		}
	}()
	<-exit
}

func main() {
	if err := setupSigSvr().Run(os.Args); err != nil {
		cmd.PrintErrorMsg(err.Error())
		os.Exit(1)
	}
}
