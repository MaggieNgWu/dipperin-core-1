// Copyright 2019, Keychain Foundation Ltd.
// This file is part of the dipperin-core library.
//
// The dipperin-core library is free software: you can redistribute
// it and/or modify it under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// The dipperin-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package commands

import (
	"encoding/json"
	"errors"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/third-party/log"
	"io/ioutil"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/core/accounts"
	"github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/rpc-interface"
	"github.com/dipperin/dipperin-core/tests/factory"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestInitRpcClient(t *testing.T) {
	assert.Panics(t, func() {
		InitRpcClient(12345)
	})
}

func TestInitAccountInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client = NewMockRpcClient(ctrl)

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(errors.New("test")).Times(2)
	InitAccountInfo(chain_config.NodeTypeOfVerifier, "", "", "")

	osExit = func(code int) {

	}

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(errors.New("test")).Times(2)

	InitAccountInfo(chain_config.NodeTypeOfNormal, "", "", "")
}

func TestCheckDownloaderSyncStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client = NewMockRpcClient(ctrl)
	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(errors.New("test"))

	assert.Panics(t, func() {
		CheckDownloaderSyncStatus()
	})

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, args ...interface{}) error {
		result = false
		return nil
	})

	assert.NotPanics(t, func() {
		CheckDownloaderSyncStatus()
	})

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, args ...interface{}) error {
		*result.(*bool) = true
		return nil
	}).Times(1)

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, args ...interface{}) error {
		*result.(*bool) = true
		return errors.New("test")
	}).Times(1)

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, args ...interface{}) error {
		*result.(*bool) = false
		return nil
	}).Times(1)

	CheckSyncStatusDuration = 1 * time.Millisecond

	assert.NotPanics(t, func() {
		CheckDownloaderSyncStatus()
	})

	client = nil
}

func TestRpcCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {

		assert.Panics(t, func() {
			RpcCall(c)
		})

		client = NewMockRpcClient(ctrl)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		RpcCall(c)
		//client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(),gomock.Any()).Return(nil).AnyTimes()

		RpcCall(c)

		//client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		SyncStatus.Store(false)
		RpcCall(c)
	}

	app.Run([]string{"xxx", "GetDefaultAccountBalance"})
	client = nil
}

func Test_getRpcParamFromString(t *testing.T) {
	assert.Equal(t, getRpcParamFromString(""), []string{})
	assert.Equal(t, getRpcParamFromString("test,test1"), []string{"test", "test1"})
}

func Test_getRpcMethodAndParam(t *testing.T) {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {

		c.Set("p", "test")
		mName, cParams, err := getRpcMethodAndParam(c)

		assert.Equal(t, mName, "test")
		assert.Equal(t, cParams, []string{"test"})
		assert.NoError(t, err)
	}

	app.Run([]string{"xxx", "test"})
}

func Test_checkSync(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client = NewMockRpcClient(ctrl)

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(errors.New("test"))
	SyncStatus.Store(false)

	assert.Equal(t, checkSync(), true)

	SyncStatus.Store(true)

	assert.Equal(t, checkSync(), false)

	client = nil
}

func Test_rpcCaller_GetDefaultAccountBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.GetDefaultAccountBalance(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		caller.GetDefaultAccountBalance(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, args ...interface{}) error {

			*result.(*rpc_interface.CurBalanceResp) = rpc_interface.CurBalanceResp{
				Balance: (*hexutil.Big)(big.NewInt(1)),
			}
			return nil
		})

		caller.GetDefaultAccountBalance(c)

	}

	app.Run([]string{"xxx"})
	client = nil
}

func Test_rpcCaller_CurrentBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		//caller.CurrentBalance(c)
		//
		//c.Set("p", "test, test, test")
		//caller.CurrentBalance(c)
		//
		//c.Set("p", "test")
		//caller.CurrentBalance(c)

		c.Set("p", common.HexToAddress("0x1234").Hex())
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.CurrentBalance(c)

		c.Set("p", "")

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), getDipperinRpcMethodByName("ListWallet")).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*[]accounts.WalletIdentifier) = []accounts.WalletIdentifier{
				{
					WalletType: 1,
					Path:       "",
					WalletName: "",
				},
			}
			return nil
		})
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), getDipperinRpcMethodByName("ListWalletAccount"), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*[]accounts.Account) = []accounts.Account{
				{
					Address: common.HexToAddress("0x1234"),
				},
			}
			return nil
		})
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		caller.CurrentBalance(c)

		c.Set("p", common.HexToAddress("0x1234").Hex())
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*rpc_interface.CurBalanceResp) = rpc_interface.CurBalanceResp{
				Balance: (*hexutil.Big)(big.NewInt(1)),
			}
			return nil
		})
		caller.CurrentBalance(c)
	}

	app.Run([]string{"xxx", "CurrentBalance"})
	client = nil
}

func Test_printBlockInfo(t *testing.T) {
	tx, _ := factory.CreateTestTx()
	m, _ := model.NewVoteMsgWithSign(uint64(1), uint64(1), common.HexToHash("0x1234"), model.PreVoteMessage, func(hash []byte) ([]byte, error) {
		return nil, nil
	}, common.HexToAddress("0x1234"))
	respBlock := rpc_interface.BlockResp{
		Body: model.Body{Txs: []*model.Transaction{tx}, Vers: []model.AbstractVerification{m}},
		Header: model.Header{
			Bloom: iblt.NewBloom(iblt.NewBloomConfig(8, 4)),
		},
	}

	printBlockInfo(respBlock)
}

func Test_printTransactionInfo(t *testing.T) {
	tx, _ := factory.CreateTestTx()
	respTx := rpc_interface.TransactionResp{
		Transaction: tx,
	}
	printTransactionInfo(respTx)
	printTransactionInfo(rpc_interface.TransactionResp{})
}

func Test_rpcCaller_CurrentBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		//caller.CurrentBlock(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(errors.New("test"))

		caller.CurrentBlock(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			tx, _ := factory.CreateTestTx()
			m, _ := model.NewVoteMsgWithSign(uint64(1), uint64(1), common.HexToHash("0x1234"), model.PreVoteMessage, func(hash []byte) ([]byte, error) {
				return nil, nil
			}, common.HexToAddress("0x1234"))
			*result.(*rpc_interface.BlockResp) = rpc_interface.BlockResp{
				Body: model.Body{Txs: []*model.Transaction{tx}, Vers: []model.AbstractVerification{m}},
				Header: model.Header{
					Bloom: iblt.NewBloom(iblt.NewBloomConfig(8, 4)),
				},
			}
			return nil
		})

		caller.CurrentBlock(c)
	}

	app.Run([]string{"xxx", "CurrentBlock"})
	client = nil
}

func Test_rpcCaller_GetGenesis(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		//caller.GetGenesis(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(errors.New("test"))

		caller.GetGenesis(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			tx, _ := factory.CreateTestTx()
			m, _ := model.NewVoteMsgWithSign(uint64(1), uint64(1), common.HexToHash("0x1234"), model.PreVoteMessage, func(hash []byte) ([]byte, error) {
				return nil, nil
			}, common.HexToAddress("0x1234"))
			*result.(*rpc_interface.BlockResp) = rpc_interface.BlockResp{
				Body: model.Body{Txs: []*model.Transaction{tx}, Vers: []model.AbstractVerification{m}},
				Header: model.Header{
					Bloom: iblt.NewBloom(iblt.NewBloomConfig(8, 4)),
				},
			}
			return nil
		})

		caller.GetGenesis(c)
	}

	app.Run([]string{"xxx", "GetGenesis"})
	client = nil
}

func Test_rpcCaller_GetBlockByNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "m", Usage: "operation"},
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}

		caller.GetBlockByNumber(c)

		c.Set("p", "")
		caller.GetBlockByNumber(c)

		c.Set("p", "s")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.GetBlockByNumber(c)

		c.Set("p", "1")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, args ...interface{}) error {
			tx, _ := factory.CreateTestTx()
			m, _ := model.NewVoteMsgWithSign(uint64(1), uint64(1), common.HexToHash("0x1234"), model.PreVoteMessage, func(hash []byte) ([]byte, error) {
				return nil, nil
			}, common.HexToAddress("0x1234"))
			*result.(*rpc_interface.BlockResp) = rpc_interface.BlockResp{
				Body: model.Body{Txs: []*model.Transaction{tx}, Vers: []model.AbstractVerification{m}},
				Header: model.Header{
					Bloom: iblt.NewBloom(iblt.NewBloomConfig(8, 4)),
				},
			}
			return nil
		})
		caller.GetBlockByNumber(c)
	}

	app.Run([]string{"xxx", "GetBlockByNumber"})
	client = nil
}

func Test_rpcCaller_GetBlockByHash(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "m", Usage: "operation"},
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}

		caller.GetBlockByHash(c)

		c.Set("p", "")
		caller.GetBlockByHash(c)

		c.Set("p", "s")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.GetBlockByHash(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, args ...interface{}) error {
			tx, _ := factory.CreateTestTx()
			m, _ := model.NewVoteMsgWithSign(uint64(1), uint64(1), common.HexToHash("0x1234"), model.PreVoteMessage, func(hash []byte) ([]byte, error) {
				return nil, nil
			}, common.HexToAddress("0x1234"))
			*result.(*rpc_interface.BlockResp) = rpc_interface.BlockResp{
				Body: model.Body{Txs: []*model.Transaction{tx}, Vers: []model.AbstractVerification{m}},
				Header: model.Header{
					Bloom: iblt.NewBloom(iblt.NewBloomConfig(8, 4)),
				},
			}
			return nil
		})
		caller.GetBlockByHash(c)
	}

	app.Run([]string{"xxx", "GetBlockByHash"})
	client = nil
}

func Test_rpcCaller_StartMine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		SyncStatus.Store(true)
		//caller.StartMine(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.StartMine(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(nil)
		caller.StartMine(c)
	}

	app.Run([]string{"xxx", "StartMine"})
	client = nil
}

func Test_rpcCaller_StopMine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		//caller.StopMine(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.StopMine(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(nil)
		caller.StopMine(c)
	}

	app.Run([]string{"xxx", "StopMine"})
	client = nil
}

func Test_rpcCaller_SetMineCoinBase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		//caller.SetMineCoinBase(c)
		//
		//c.Set("p", "")
		//
		//caller.SetMineCoinBase(c)
		//
		//c.Set("p", "ttt")
		//
		//caller.SetMineCoinBase(c)

		c.Set("p", common.HexToAddress("0x1234").Hex())

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.SetMineCoinBase(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		caller.SetMineCoinBase(c)
	}

	app.Run([]string{"xxx", "SetMineCoinBase"})
	client = nil
}

func Test_rpcCaller_SendTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		SyncStatus.Store(false)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(nil)
		caller.SendTx(c)

		SyncStatus.Store(true)
		caller.SendTx(c)

		c.Set("p", "test")
		caller.SendTx(c)

		c.Set("p", "test,test,test")
		caller.SendTx(c)

		c.Set("p", common.HexToAddress("0x1234").Hex()+",test,test")
		caller.SendTx(c)

		c.Set("p", common.HexToAddress("0x1234").Hex()+",10,test")
		caller.SendTx(c)

		c.Set("p", common.HexToAddress("0x1234").Hex()+",10,10,test")

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.SendTx(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		caller.SendTx(c)
	}

	app.Run([]string{"xxx", "SendTx"})
	client = nil
}

func Test_rpcCaller_SendTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		SyncStatus.Store(false)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		caller.SendTransaction(c)

		SyncStatus.Store(true)
		caller.SendTransaction(c)

		c.Set("p", "test")
		caller.SendTransaction(c)

		c.Set("p", "test,test,test,test")
		caller.SendTransaction(c)

		c.Set("p", common.HexToAddress("0x1234").Hex()+",test,test,test")
		caller.SendTransaction(c)

		c.Set("p", common.HexToAddress("0x1234").Hex()+","+common.HexToAddress("0x1234").Hex()+",test,test")
		caller.SendTransaction(c)

		c.Set("p", common.HexToAddress("0x1234").Hex()+","+common.HexToAddress("0x1234").Hex()+",10,test")
		caller.SendTransaction(c)

		c.Set("p", common.HexToAddress("0x1234").Hex()+","+common.HexToAddress("0x1234").Hex()+",10,10,test")

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test")).AnyTimes()
		caller.SendTransaction(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

		caller.SendTransaction(c)
	}

	app.Run([]string{"xxx", "SendTransaction"})
	client = nil
}

func TestRpcCaller_SendTransactionContract(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
		cli.StringFlag{Name: "abi", Usage: "abi path"},
		cli.StringFlag{Name: "wasm", Usage: "wasm path"},
		cli.StringFlag{Name: "input", Usage: "contract params"},
		cli.BoolFlag{Name: "is-create", Usage: "create contract or not"},
		cli.StringFlag{Name: "func-name", Usage: "call function name"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		SyncStatus.Store(true)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		//client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		caller.SendTransactionContract(c)

		SyncStatus.Store(true)
		//	caller.SendTransactionContract(c)

		c.Set("p", "0x000062be10f46b5d01Ecd9b502c4bA3d6131f6fc2e41,0x00144179D57e45Cb3b54D6FAEF69e746bf240E287978,11122,10")
		c.Set("abi", util.HomeDir()+"go/src/github.com/dipperin/dipperin-core/core/vm/test-data/event/example/example.cpp.abi.json")
		c.Set("input", "test,123,456")
		c.Set("func-name", "fake")
		c.Set("is-create", "false")
		//c.Set("abi", "Users/konggan/workspace/chain/dipperin/dipc/cmake-build-debug/example/example.cpp.abi.json")
		caller.SendTransactionContract(c)

	}

	app.Run([]string{"xxx", "SendTransactionContract"})
	client = nil
}

func TestRpcCaller_SendTransactionContract2(t *testing.T) {
	/*	app := cli.NewApp()

		app.Flags = []cli.Flag{
			cli.StringFlag{Name: "m", Usage: "operation"},
			cli.StringFlag{Name: "p", Usage: "parameters"},
			cli.StringFlag{Name: "abi", Usage:"abi path"},
			cli.StringFlag{Name: "wasm", Usage:"wasm path"},
			cli.StringFlag{Name: "input", Usage: "contract params"},
			cli.BoolFlag{Name:   "isCreate", Usage: "create contract or not"},
			cli.StringFlag{Name: "funcName", Usage: "call function name"},
		}

		app.Action = func(c *cli.Context) {
			caller := &rpcCaller{}

			c.Set("m", "SendTransactionContract")
			c.Set("p", "0x000062be10f46b5d01Ecd9b502c4bA3d6131f6fc2e41,10,11122,10")
			c.Set("abi", "Users/konggan/workspace/chain/dipperin/dipc/cmake-build-debug/example/example.cpp.abi.json")
			c.Set("wasm", "/Users/konggan/workspace/chain/dipperin/dipc/cmake-build-debug/example/example.wasm")
			c.Set("isCreate", "true")
			caller.SendTransactionContract(c)
		}*/
	//app := cli.NewApp()
	//
	//app.Flags = []cli.Flag{
	//	cli.StringFlag{Name: "m", Usage: "operation"},
	//	cli.StringFlag{Name: "p", Usage: "parameters"},
	//	cli.StringFlag{Name: "abi", Usage:"abi path"},
	//	cli.StringFlag{Name: "wasm", Usage:"wasm path"},
	//	cli.StringFlag{Name: "input", Usage: "contract params"},
	//	cli.BoolFlag{Name:   "isCreate", Usage: "create contract or not"},
	//	cli.StringFlag{Name: "funcName", Usage: "call function name"},
	//}
	//c := cli.NewContext(app, nil, nil)
	//c.Set("m", "SendTransactionContract")
	//c.Set("p", "0x000062be10f46b5d01Ecd9b502c4bA3d6131f6fc2e41,10,11122,10")
	//c.Set("abi", "Users/konggan/workspace/chain/dipperin/dipc/cmake-build-debug/example/example.cpp.abi.json")
	//c.Set("wasm", "/Users/konggan/workspace/chain/dipperin/dipc/cmake-build-debug/example/example.wasm")
	//c.Set("isCreate", "true")
	//
	//SendTransactionContractCreate(c)

}

/*func Test_AbiFile(t *testing.T)  {
	abiBytes, err := ioutil.ReadFile("/Users/konggan/workspace/chain/dipperin/dipc/cmake-build-debug/example/example.cpp.abi.json")
	assert.NoError(t, err)
	var wasmAbi utils.WasmAbi
	err = json.Unmarshal(abiBytes, &wasmAbi)
	assert.NoError(t, err)

}*/

func Test_generateTxData(t *testing.T) {
	log.InitLogger(log.LvlDebug)
	transactionJson := "{\"TxData\":{\"nonce\":\"0x0\",\"to\":\"0x00120000000000000000000000000000000000000000\",\"hashlock\":null,\"timelock\":\"0x0\",\"value\":\"0x2540be400\",\"fee\":\"0x69db9c0\",\"gasPrice\":\"0xa\",\"gas\":\"0x1027127dc00\",\"input\":\"0xf9027b823138b8eb0061736d01000000010d0360017f0060027f7f00600000021d0203656e76067072696e7473000003656e76087072696e74735f6c00010304030202000405017001010105030100020615037f01419088040b7f00419088040b7f004186080b073405066d656d6f727902000b5f5f686561705f6261736503010a5f5f646174615f656e64030204696e697400030568656c6c6f00040a450302000b02000b3d01017f230041106b220124004180081000200141203a000f2001410f6a41011001200010002001410a3a000e2001410e6a41011001200141106a24000b0b0d01004180080b0668656c6c6f00b901887b22616269417272223a5b0a202020207b0a2020202020202020226e616d65223a2022696e6974222c0a202020202020202022696e70757473223a205b5d2c0a2020202020202020226f757470757473223a205b5d2c0a202020202020202022636f6e7374616e74223a202266616c7365222c0a20202020202020202274797065223a202266756e6374696f6e220a202020207d2c0a202020207b0a2020202020202020226e616d65223a202268656c6c6f222c0a202020202020202022696e70757473223a205b0a2020202020202020202020207b0a20202020202020202020202020202020226e616d65223a20226e616d65222c0a202020202020202020202020202020202274797065223a2022737472696e67220a2020202020202020202020207d0a20202020202020205d2c0a2020202020202020226f757470757473223a205b5d2c0a202020202020202022636f6e7374616e74223a202274727565222c0a20202020202020202274797065223a202266756e6374696f6e220a202020207d0a5d7d0a\"},\"Wit\":{\"r\":\"0x30e173f7590a6e12bb4d51bbf6ae113ee668245d2e30a145d1845f55ae5a9f4a\",\"s\":\"0x7d9f36d62573ac09e1dd84d31650a8b5e20b5dffb34d3955dde224c61d299744\",\"v\":\"0x39\",\"hashkey\":\"0x\"}}"

	var tx model.Transaction
	err := json.Unmarshal([]byte(transactionJson), &tx)
	if err != nil {
		log.Info("TestAccountStateDB_ProcessContract", "err", err)
	}
	log.Info("processContract", "Tx", tx)

}

func Test_rpcCaller_Transaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		SyncStatus.Store(false)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(nil)
		caller.Transaction(c)

		SyncStatus.Store(true)
		caller.Transaction(c)

		c.Set("p", "")
		caller.Transaction(c)

		c.Set("p", "test")
		caller.Transaction(c)

		c.Set("p", "0x1234")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.Transaction(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, args ...interface{}) error {
			tx, _ := factory.CreateTestTx()
			*result.(*rpc_interface.TransactionResp) = rpc_interface.TransactionResp{
				Transaction: tx,
			}
			return nil
		})

		caller.Transaction(c)

	}

	app.Run([]string{"xxx", "Transaction"})
	client = nil
}

func Test_rpcCaller_ListWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}

		//caller.ListWallet(c)

		c.Set("p", "")

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.ListWallet(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, args ...interface{}) error {
			*result.(*[]accounts.WalletIdentifier) = []accounts.WalletIdentifier{
				{
					WalletType: 1,
					WalletName: "test",
					Path:       "",
				},
			}
			return nil
		})

		caller.ListWallet(c)

	}

	app.Run([]string{"xxx", "ListWallet"})
	client = nil
}

func Test_rpcCaller_ListWalletAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}

		//caller.ListWalletAccount(c)
		//
		//c.Set("p", "test")
		//
		//caller.ListWalletAccount(c)
		//
		//c.Set("p", "test, test")
		//caller.ListWalletAccount(c)

		c.Set("p", "")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.ListWalletAccount(c)

		c.Set("p", "SoftWallet, test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.ListWalletAccount(c)

		c.Set("p", "LedgerWallet, test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.ListWalletAccount(c)

		c.Set("p", "TrezorWallet, test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, args ...interface{}) error {
			*result.(*[]accounts.Account) = []accounts.Account{
				{
					Address: common.HexToAddress("0x1234"),
				},
			}
			return nil
		})
		caller.ListWalletAccount(c)

	}

	app.Run([]string{"xxx", "ListWalletAccount"})
	client = nil
}

func Test_rpcCaller_EstablishWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}

		caller.EstablishWallet(c)

		c.Set("p", "test")

		caller.EstablishWallet(c)

		c.Set("p", "test,test,test")
		caller.EstablishWallet(c)

		c.Set("p", "SoftWallet,test,test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.EstablishWallet(c)

		c.Set("p", "LedgerWallet,test,test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.EstablishWallet(c)

		c.Set("p", "TrezorWallet,test,test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, args ...interface{}) error {
			*result.(*string) = "test"
			return nil
		})
		caller.EstablishWallet(c)

	}

	app.Run([]string{"xxx", "EstablishWallet"})
	client = nil
}

func Test_rpcCaller_RestoreWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}

		caller.RestoreWallet(c)

		c.Set("p", "test")

		caller.RestoreWallet(c)

		c.Set("p", "test,test,test,test")
		caller.RestoreWallet(c)

		c.Set("p", "SoftWallet,test,test,test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.RestoreWallet(c)

		c.Set("p", "LedgerWallet,test,test,test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.RestoreWallet(c)

		c.Set("p", "TrezorWallet,test,test,test,test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		caller.RestoreWallet(c)

	}

	app.Run([]string{"xxx", "RestoreWallet"})
	client = nil
}

func Test_rpcCaller_OpenWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}

		caller.OpenWallet(c)

		c.Set("p", "")

		caller.OpenWallet(c)

		c.Set("p", "test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.OpenWallet(c)

		c.Set("p", "test,test,test")
		caller.OpenWallet(c)

		c.Set("p", "SoftWallet,test,test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.OpenWallet(c)

		c.Set("p", "LedgerWallet,test,test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.OpenWallet(c)

		c.Set("p", "TrezorWallet,test,test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		caller.OpenWallet(c)

	}

	app.Run([]string{"xxx", "OpenWallet"})
	client = nil
}

func Test_rpcCaller_CloseWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}

		//caller.CloseWallet(c)
		//
		//c.Set("p", "test")
		//
		//caller.CloseWallet(c)

		c.Set("p", "")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.CloseWallet(c)

		c.Set("p", "test,test")
		caller.CloseWallet(c)

		c.Set("p", "SoftWallet,test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.CloseWallet(c)

		c.Set("p", "LedgerWallet,test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.CloseWallet(c)

		c.Set("p", "TrezorWallet,test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		caller.CloseWallet(c)

	}

	app.Run([]string{"xxx", "CloseWallet"})
	client = nil
}

func Test_rpcCaller_AddAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}

		//caller.AddAccount(c)
		//
		//c.Set("p", "test")
		//
		//caller.AddAccount(c)

		c.Set("p", "")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.AddAccount(c)

		c.Set("p", "test,test")
		caller.AddAccount(c)

		c.Set("p", "SoftWallet,test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.AddAccount(c)

		c.Set("p", "LedgerWallet,test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.AddAccount(c)

		c.Set("p", "TrezorWallet,test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, args ...interface{}) error {
			*result.(*accounts.Account) = accounts.Account{
				Address: common.HexToAddress("0x1234"),
			}
			return nil
		})
		caller.AddAccount(c)

	}

	app.Run([]string{"xxx", "AddAccount"})
	client = nil
}

func Test_rpcCaller_SendRegisterTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		SyncStatus.Store(false)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(errors.New("test")).Times(1)
		caller.SendRegisterTx(c)

		SyncStatus.Store(true)
		caller.SendRegisterTx(c)

		c.Set("p", "test")
		caller.SendRegisterTx(c)

		c.Set("p", "test,test")
		caller.SendRegisterTx(c)

		c.Set("p", "10,test")
		caller.SendRegisterTx(c)

		c.Set("p", "10,10")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.SendRegisterTx(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Nil()).Return(nil)

		caller.SendRegisterTx(c)

	}

	app.Run([]string{"xxx", "SendRegisterTx"})
	client = nil
}

func Test_rpcCaller_SendRegisterTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		SyncStatus.Store(false)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(errors.New("test")).Times(1)
		caller.SendRegisterTransaction(c)

		SyncStatus.Store(true)
		caller.SendRegisterTransaction(c)

		c.Set("p", "test")
		caller.SendRegisterTransaction(c)

		c.Set("p", "test,test,test")
		caller.SendRegisterTransaction(c)

		c.Set("p", common.HexToAddress("0x1234").Hex()+",test,test")
		caller.SendRegisterTransaction(c)

		c.Set("p", common.HexToAddress("0x1234").Hex()+",10,test")
		caller.SendRegisterTransaction(c)

		c.Set("p", common.HexToAddress("0x1234").Hex()+",10,10")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.SendRegisterTransaction(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Nil()).Return(nil)

		caller.SendRegisterTransaction(c)

	}

	app.Run([]string{"xxx", "SendRegisterTransaction"})
	client = nil
}

func Test_rpcCaller_SendUnStakeTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		SyncStatus.Store(false)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(errors.New("test")).Times(1)
		caller.SendUnStakeTx(c)

		SyncStatus.Store(true)
		caller.SendUnStakeTx(c)

		c.Set("p", "")
		caller.SendUnStakeTx(c)

		c.Set("p", "test")
		caller.SendUnStakeTx(c)

		c.Set("p", "10")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.SendUnStakeTx(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		caller.SendUnStakeTx(c)

	}

	app.Run([]string{"xxx", "SendUnStakeTx"})
	client = nil
}

func Test_rpcCaller_SendUnStakeTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		SyncStatus.Store(false)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(errors.New("test")).Times(1)
		caller.SendUnStakeTransaction(c)

		SyncStatus.Store(true)
		caller.SendUnStakeTransaction(c)

		c.Set("p", "")
		caller.SendUnStakeTransaction(c)

		c.Set("p", "test,test")
		caller.SendUnStakeTransaction(c)

		c.Set("p", common.HexToAddress("0x1234").Hex()+",test")
		caller.SendUnStakeTransaction(c)

		c.Set("p", common.HexToAddress("0x1234").Hex()+",10")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.SendUnStakeTransaction(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		caller.SendUnStakeTransaction(c)

	}

	app.Run([]string{"xxx", "SendUnStakeTransaction"})
	client = nil
}

func Test_rpcCaller_SendCancelTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		SyncStatus.Store(false)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(errors.New("test")).Times(1)
		caller.SendCancelTx(c)

		SyncStatus.Store(true)
		caller.SendCancelTx(c)

		c.Set("p", "")
		caller.SendCancelTx(c)

		c.Set("p", "test")
		caller.SendCancelTx(c)

		c.Set("p", "10")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.SendCancelTx(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		caller.SendCancelTx(c)

	}

	app.Run([]string{"xxx", "SendCancelTx"})
	client = nil
}

func Test_rpcCaller_SendCancelTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		SyncStatus.Store(false)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(errors.New("test")).Times(1)
		caller.SendCancelTransaction(c)

		SyncStatus.Store(true)
		caller.SendCancelTransaction(c)

		c.Set("p", "")
		caller.SendCancelTransaction(c)

		c.Set("p", "test,test")
		caller.SendCancelTransaction(c)

		c.Set("p", common.HexToAddress("0x1234").Hex()+",test")
		caller.SendCancelTransaction(c)

		c.Set("p", common.HexToAddress("0x1234").Hex()+",10")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.SendCancelTransaction(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		caller.SendCancelTransaction(c)

	}

	app.Run([]string{"xxx", "SendCancelTransaction"})
	client = nil
}

func Test_rpcCaller_GetVerifiersBySlot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		SyncStatus.Store(false)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(errors.New("test")).Times(1)
		caller.GetVerifiersBySlot(c)

		SyncStatus.Store(true)
		caller.GetVerifiersBySlot(c)

		c.Set("p", "")
		caller.GetVerifiersBySlot(c)

		c.Set("p", "test")
		caller.GetVerifiersBySlot(c)

		c.Set("p", "1")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.GetVerifiersBySlot(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*[]common.Address) = []common.Address{common.HexToAddress("0x1234")}
			return nil
		})

		caller.GetVerifiersBySlot(c)

	}

	app.Run([]string{"xxx", "GetVerifiersBySlot"})
	client = nil
}

func Test_rpcCaller_VerifierStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		//caller.VerifierStatus(c)
		//
		//c.Set("p", "test,test")
		//caller.VerifierStatus(c)
		//
		//c.Set("p", "1234")
		//caller.VerifierStatus(c)

		c.Set("p", "")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test")).Times(2)
		caller.VerifierStatus(c)

		c.Set("p", common.HexToAddress("0x1234").Hex())

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*rpc_interface.VerifierStatus) = rpc_interface.VerifierStatus{}
			return nil
		})

		caller.VerifierStatus(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*rpc_interface.VerifierStatus) = rpc_interface.VerifierStatus{
				Balance: (*hexutil.Big)(big.NewInt(1)),
				Status:  VerifierStatusNoRegistered,
			}
			return nil
		})

		caller.VerifierStatus(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*rpc_interface.VerifierStatus) = rpc_interface.VerifierStatus{
				Balance: (*hexutil.Big)(big.NewInt(1)),
				Status:  VerifiedStatusUnstaked,
			}
			return nil
		})

		caller.VerifierStatus(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*rpc_interface.VerifierStatus) = rpc_interface.VerifierStatus{
				Balance: (*hexutil.Big)(big.NewInt(1)),
				Status:  VerifierStatusRegistered,
			}
			return nil
		})

		caller.VerifierStatus(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*rpc_interface.VerifierStatus) = rpc_interface.VerifierStatus{
				Balance: (*hexutil.Big)(big.NewInt(1)),
				Stake:   (*hexutil.Big)(big.NewInt(1)),
				Status:  VerifierStatusRegistered,
			}
			return nil
		})

		caller.VerifierStatus(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*rpc_interface.VerifierStatus) = rpc_interface.VerifierStatus{
				Balance: (*hexutil.Big)(big.NewInt(1)),
				Stake:   (*hexutil.Big)(big.NewInt(1)),
				Status:  VerifiedStatusCanceled,
			}
			return nil
		})

		caller.VerifierStatus(c)

	}

	app.Run([]string{"xxx", "VerifierStatus"})
	client = nil
}

func Test_rpcCaller_SetBftSigner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test")).AnyTimes()
		caller.SetBftSigner(c)

		c.Set("p", "")
		caller.SetBftSigner(c)

		c.Set("p", "test")
		caller.SetBftSigner(c)

		c.Set("p", common.HexToAddress("0x1234").Hex())
		caller.SetBftSigner(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		caller.SetBftSigner(c)
	}

	app.Run([]string{"xxx", "SetBftSigner"})
	client = nil
}

func Test_rpcCaller_GetDefaultAccountStake(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))

		caller.GetDefaultAccountStake(c)
	}

	app.Run([]string{"xxx", "GetDefaultAccountStake"})
	client = nil
}

func Test_rpcCaller_CurrentStake(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		//caller.CurrentStake(c)
		//
		//c.Set("p", "test,test")
		//caller.CurrentStake(c)
		//
		//c.Set("p", "test")
		//caller.CurrentStake(c)

		c.Set("p", "")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test")).Times(2)
		caller.CurrentStake(c)

		c.Set("p", common.HexToAddress("0x1234").Hex())
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*rpc_interface.CurBalanceResp) = rpc_interface.CurBalanceResp{}
			return nil
		})
		caller.CurrentStake(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*rpc_interface.CurBalanceResp) = rpc_interface.CurBalanceResp{
				Balance: (*hexutil.Big)(big.NewInt(1)),
			}
			return nil
		})
		caller.CurrentStake(c)
	}

	app.Run([]string{"xxx", "CurrentStake"})
	client = nil
}

func Test_rpcCaller_CurrentReputation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		//caller.CurrentReputation(c)
		//
		//c.Set("p", "test,test")
		//caller.CurrentReputation(c)
		//
		//c.Set("p", "test")
		//caller.CurrentReputation(c)

		c.Set("p", "")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test")).Times(2)
		caller.CurrentReputation(c)

		c.Set("p", common.HexToAddress("0x1234").Hex())
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		caller.CurrentReputation(c)
	}

	app.Run([]string{"xxx", "CurrentReputation"})
	client = nil
}

func Test_rpcCaller_GetCurVerifiers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		SyncStatus.Store(false)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test")).AnyTimes()
		caller.GetCurVerifiers(c)

		SyncStatus.Store(true)
		caller.GetCurVerifiers(c)

		//client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.GetCurVerifiers(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*[]common.Address) = []common.Address{common.HexToAddress("0x1234")}
			return nil
		}).AnyTimes()
		caller.GetCurVerifiers(c)
	}

	app.Run([]string{"xxx", "GetCurVerifiers"})
	client = nil
}

func Test_inDefaultVs(t *testing.T) {
	address := common.HexToAddress("0x1234")
	assert.Equal(t, inDefaultVs(address), false)
	address = chain_config.LocalVerifierAddress[0]
	assert.Equal(t, inDefaultVs(address), true)
}

func Test_rpcCaller_GetNextVerifiers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		SyncStatus.Store(false)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test")).MaxTimes(3)
		caller.GetNextVerifiers(c)

		SyncStatus.Store(true)
		caller.GetNextVerifiers(c)

		//client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.GetNextVerifiers(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*[]common.Address) = []common.Address{common.HexToAddress("0x1234")}
			return nil
		})
		caller.GetNextVerifiers(c)
	}

	app.Run([]string{"xxx", "GetNextVerifiers"})
	client = nil
}

func Test_getNonceInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		//nonce, err := getNonceInfo(c)
		//assert.Error(t, err)
		//assert.Equal(t, nonce, uint64(0))
		//
		//c.Set("p", "test,test")
		//
		//nonce, err = getNonceInfo(c)
		//assert.Error(t, err)
		//assert.Equal(t, nonce, uint64(0))

		c.Set("p", "")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test")).Times(2)
		nonce, err := getNonceInfo(c)
		assert.Error(t, err)
		assert.Equal(t, nonce, uint64(0))

		c.Set("p", "test")
		nonce, err = getNonceInfo(c)
		assert.Error(t, err)
		assert.Equal(t, nonce, uint64(0))

		c.Set("p", common.HexToAddress("0x1234").Hex())
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		nonce, err = getNonceInfo(c)
		assert.NoError(t, err)
		assert.Equal(t, nonce, uint64(0))
	}

	app.Run([]string{"xxx", "getNonceInfo"})
	client = nil
}

func Test_rpcCaller_GetTransactionNonce(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		SyncStatus.Store(false)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test")).AnyTimes()
		caller.GetTransactionNonce(c)

		SyncStatus.Store(true)
		caller.GetTransactionNonce(c)

		c.Set("p", common.HexToAddress("0x1234").Hex())
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		caller.GetTransactionNonce(c)
	}

	app.Run([]string{"xxx", "GetTransactionNonce"})
	client = nil
}

func Test_rpcCaller_GetAddressNonceFromWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
	}

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		//caller.GetAddressNonceFromWallet(c)

		c.Set("p", common.HexToAddress("0x1234").Hex())
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		caller.GetAddressNonceFromWallet(c)
	}

	app.Run([]string{"xxx", "GetAddressNonceFromWallet"})
	client = nil
}

func testTempJSONFile(t *testing.T) (string, func()) {
	t.Helper()
	tf, err := ioutil.TempFile("", "*.json")
	if err != nil {
		t.Fatalf(":err: %s", err)
	}

	tf.Close()
	return tf.Name(), func() { os.Remove(tf.Name()) }
}

func Test_initWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tf, tfClean := testTempJSONFile(t)
	defer tfClean()
	client = NewMockRpcClient(ctrl)

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
	err := initWallet(tf, "test", "test")
	assert.Error(t, err)

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	err = initWallet("test", "test", "test")

	client = nil
}

func Test_getDefaultAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client = NewMockRpcClient(ctrl)

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
		*result.(*[]accounts.WalletIdentifier) = []accounts.WalletIdentifier{{}}
		return nil
	})

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))

	address := getDefaultAccount()

	assert.Equal(t, address, common.Address{})
}

func Test_getDefaultWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client = NewMockRpcClient(ctrl)

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
		*result.(*[]accounts.WalletIdentifier) = []accounts.WalletIdentifier{{}}
		return nil
	})

	wallet := getDefaultWallet()

	assert.Equal(t, wallet, accounts.WalletIdentifier{})
}
