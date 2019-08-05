package vm

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/vm/common/utils"
	model2 "github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/stretchr/testify/assert"
	"github.com/vntchain/go-vnt/rlp"
	"math/big"
	"testing"
)

func TestNewVMContext(t *testing.T) {
	block := model.CreateBlock(0, common.Hash{}, 1)
	tx := block.GetTransactions()[0]
	from, err := tx.Sender(nil)
	assert.NoError(t, err)
	header := block.Header()

	context := NewVMContext(tx, header, getTestHashFunc())
	assert.Equal(t, tx.GetGasLimit(), context.GetGasLimit())
	assert.Equal(t, header.GetNumber(), context.GetBlockNumber().Uint64())
	assert.Equal(t, tx.GetGasPrice(), context.GetGasPrice())
	assert.Equal(t, common.Hash{}, context.GetBlockHash(0))
	assert.Equal(t, tx.CalTxId(), context.GetTxHash())
	assert.Equal(t, header.CoinBaseAddress(), context.GetCoinBase())
	assert.Equal(t, from, context.GetOrigin())
	assert.Equal(t, header.GetTimeStamp(), context.GetTime())

	db := NewFakeStateDB()
	db.CreateAccount(aliceAddr)
	db.CreateAccount(contractAddr)
	db.AddBalance(aliceAddr, big.NewInt(500))
	result := context.CanTransfer(db, aliceAddr, big.NewInt(100))
	assert.Equal(t, true, result)

	context.Transfer(db, aliceAddr, contractAddr, big.NewInt(100))
	assert.Equal(t, big.NewInt(400), db.GetBalance(aliceAddr))
	assert.Equal(t, big.NewInt(100), db.GetBalance(contractAddr))
}

func Test_Run(t *testing.T) {
	vm := getTestVm()
	WASMPath := g_testData.GetWASMPath("event", g_testData.CoreVmTestData)
	AbiPath := g_testData.GetAbiPath("event", g_testData.CoreVmTestData)
	contract := getContract(WASMPath, AbiPath, nil)
	_, err := run(vm, contract, true)
	assert.NoError(t, err)

	ch := make(chan int, 1)
	go func() {
		vm.Cancel()
		ch <- 0
	}()
	_, err = run(vm, contract, true)
	assert.NoError(t, err)
	<-ch
}

func TestVM_CreateAndCall(t *testing.T) {
	vm := getTestVm()
	ref := AccountRef(aliceAddr)
	gasLimit := g_testData.TestGasLimit * 100
	value := g_testData.TestValue
	WASMPath := g_testData.GetWASMPath("event", g_testData.CoreVmTestData)
	AbiPath := g_testData.GetAbiPath("event", g_testData.CoreVmTestData)
	code, _ := g_testData.GetCodeAbi(WASMPath, AbiPath)
	data, err := g_testData.GetCreateExtraData(WASMPath, AbiPath, "")
	assert.NoError(t, err)

	// vm.Create
	vm.GetStateDB().CreateAccount(ref.Address())
	vm.GetStateDB().AddBalance(ref.Address(), big.NewInt(10000))
	resp, addr, gasLimit, err := vm.Create(ref, data, gasLimit, value)
	expectAddr := cs_crypto.CreateContractAddress(ref.Address(), uint64(0))
	assert.Equal(t, code, resp)
	assert.Equal(t, expectAddr, addr)
	assert.NoError(t, err)

	// vm.Call
	data, err = g_testData.GetCallExtraData("returnString", "param")
	assert.NoError(t, err)

	resp, _, err = vm.Call(ref, addr, data, gasLimit, value)
	expectResp := utils.Align32BytesConverter(resp, "string")
	assert.Equal(t, "param", expectResp)
	assert.Equal(t, expectAddr, addr)
	assert.NoError(t, err)

	// vm.DelegateCall
	fmt.Println("-----------------------------------------")
	parentContract := getContract(WASMPath, AbiPath, data)
	resp, _, err = vm.DelegateCall(parentContract, addr, data, gasLimit)
	expectResp = utils.Align32BytesConverter(resp, "string")
	assert.Equal(t, "param", expectResp)
	assert.Equal(t, expectAddr, addr)
	assert.NoError(t, err)
}

func TestVM_CreateAndCallWithdraw(t *testing.T) {
	/*	vm := getTestVm()
		aliceRef := AccountRef(aliceAddr)
		gasLimit := g_testData.TestGasLimit * 100
		value := g_testData.TestValue
		WASMPath := g_testData.GetWASMPath("token-payable", g_testData.DIPCTestContract)
		AbiPath := g_testData.GetAbiPath("token-payable", g_testData.DIPCTestContract)
		code, abi := g_testData.GetCodeAbi(WASMPath, AbiPath)

		// create alice and bob account
		vm.GetStateDB().CreateAccount(aliceAddr)
		vm.GetStateDB().AddBalance(aliceAddr, big.NewInt(500))
		vm.GetStateDB().CreateAccount(bobAddr)
		vm.GetStateDB().AddBalance(bobAddr, big.NewInt(1000))
		expectAddr := cs_crypto.CreateContractAddress(aliceAddr, uint64(0))

		tokenName := []byte("tokenName")
		symbolName := []byte("symbolName")
		supply := utils.Uint64ToBytes(500)
		data, err := rlp.EncodeToBytes([]interface{}{code, abi, tokenName, symbolName, supply})
		assert.NoError(t, err)

		// init() is not payable function, couldn't transfer DIP
		resp, addr, gasLimit, err := vm.Create(aliceRef, data, gasLimit, value)
		assert.Equal(t, code, resp)
		assert.Equal(t, expectAddr, addr)
		assert.Equal(t, value, vm.GetStateDB().GetBalance(addr))
		assert.Error(t, err)

		resp, addr, gasLimit, err = vm.Create(aliceRef, data, gasLimit, big.NewInt(0))
		assert.Equal(t, code, resp)
		assert.Equal(t, expectAddr, addr)
		assert.Equal(t, value, vm.GetStateDB().GetBalance(addr))
		assert.Error(t, err)

		fmt.Println("-----------------------------------------------------")

		// transfer token to bob and transfer DIP to contract
		funcName := []byte("transfer")
		to := []byte(bobAddr.Hex())
		amount := utils.Uint64ToBytes(100)
		data, err = rlp.EncodeToBytes([]interface{}{funcName, to, amount})
		assert.NoError(t, err)
		resp, _, err = vm.Call(aliceRef, addr, data, gasLimit, value)
		assert.Equal(t, make([]byte, utils.ALIGN_LENGTH), resp)
		assert.NoError(t, err)

		// call withdraw
		funcName = []byte("withdraw")
		data, err = rlp.EncodeToBytes([]interface{}{funcName})
		assert.NoError(t, err)

		// bob can't withdraw alice's contract
		resp, _, err = vm.Call(AccountRef(bobAddr), addr, data, gasLimit, value)
		assert.Equal(t, []byte(nil), resp)
		assert.Equal(t, value, vm.GetStateDB().GetBalance(addr))
		assert.Equal(t, "VM execute fail: abort", err.Error())

		// alice withdraw the contract
		resp, _, err = vm.Call(aliceRef, addr, data, gasLimit, value)
		assert.Equal(t, make([]byte, utils.ALIGN_LENGTH), resp)
		assert.Equal(t, uint64(0), vm.GetStateDB().GetBalance(addr).Uint64())
		assert.NoError(t, err)*/
}

func TestVM_CreateAndCallToken(t *testing.T) {
	vm := getTestVm()
	ref := AccountRef(aliceAddr)
	gasLimit := g_testData.TestGasLimit * 100
	value := g_testData.TestValue
	WASMPath := g_testData.GetWASMPath("token-const", g_testData.CoreVmTestData)
	AbiPath := g_testData.GetAbiPath("token-const", g_testData.CoreVmTestData)
	code, abi := g_testData.GetCodeAbi(WASMPath, AbiPath)

	tokenName := []byte("tokenName")
	symbolName := []byte("symbolName")
	supply := utils.Uint64ToBytes(500)
	data, err := rlp.EncodeToBytes([]interface{}{code, abi, tokenName, symbolName, supply})
	assert.NoError(t, err)

	vm.GetStateDB().CreateAccount(ref.Address())
	vm.GetStateDB().AddBalance(ref.Address(), big.NewInt(10000))
	expectAddr := cs_crypto.CreateContractAddress(ref.Address(), uint64(0))
	resp, addr, gasLimit, err := vm.Create(ref, data, gasLimit, value)
	assert.Equal(t, code, resp)
	assert.Equal(t, expectAddr, addr)
	assert.NoError(t, err)

	// vm.Call transfer
	funcName := []byte("transfer")
	to := []byte(bobAddr.Hex())
	amount := utils.Uint64ToBytes(100)
	data, err = rlp.EncodeToBytes([]interface{}{funcName, to, amount})
	assert.NoError(t, err)

	resp, _, err = vm.Call(ref, addr, data, gasLimit, value)
	assert.Equal(t, make([]byte, utils.ALIGN_LENGTH), resp)
	assert.NoError(t, err)

	// vm.Call bob getBalance
	funcName = []byte("getBalance")
	data, err = rlp.EncodeToBytes([]interface{}{funcName, to})
	assert.NoError(t, err)

	resp, _, err = vm.Call(ref, addr, data, gasLimit, value)
	expectResp := utils.Align32BytesConverter(resp, "uint64")
	assert.Equal(t, uint64(100), expectResp)
	assert.NoError(t, err)

	// vm.Call alice getBalance
	from := []byte(aliceAddr.Hex())
	data, err = rlp.EncodeToBytes([]interface{}{funcName, from})
	assert.NoError(t, err)

	resp, _, err = vm.Call(ref, addr, data, gasLimit, value)
	expectResp = utils.Align32BytesConverter(resp, "uint64")
	assert.Equal(t, uint64(400), expectResp)
	assert.NoError(t, err)
}

func TestVM_Call_Error(t *testing.T) {
	vm := getTestVm()
	ref := AccountRef(aliceAddr)
	gasLimit := g_testData.TestGasLimit
	value := g_testData.TestValue

	vm.GetStateDB().CreateAccount(ref.Address())
	vm.GetStateDB().AddBalance(ref.Address(), big.NewInt(10000))

	vm.GetStateDB().CreateAccount(contractAddr)
	vm.GetStateDB().SetCode(contractAddr, []byte{123})
	vm.GetStateDB().SetAbi(contractAddr, []byte{123})
	_, _, err := vm.Call(ref, contractAddr, nil, gasLimit, value)
	assert.Equal(t, "unexpected EOF", err.Error())

	vm = getTestVm()
	vm.GetStateDB().CreateAccount(ref.Address())
	_, _, err = vm.Call(ref, contractAddr, nil, gasLimit, value)
	assert.Equal(t, g_error.ErrInsufficientBalance, err)

	vm.GetStateDB().AddBalance(ref.Address(), big.NewInt(10000))
	_, _, err = vm.Call(ref, contractAddr, nil, gasLimit, value)
	assert.NoError(t, err)

	vm.depth = int(model2.CallCreateDepth + 1)
	_, _, err = vm.Call(ref, contractAddr, nil, gasLimit, value)
	assert.Equal(t, g_error.ErrDepth, err)

	vm.depth = 1
	vm.vmConfig.NoRecursion = true
	_, _, err = vm.Call(ref, contractAddr, nil, gasLimit, value)
	assert.NoError(t, err)
}

func TestVM_DelegateCall_Error(t *testing.T) {
	vm := getTestVm()
	caller := AccountRef(aliceAddr)
	self := AccountRef(contractAddr)
	value := g_testData.TestValue
	gasLimit := g_testData.TestGasLimit
	contract := NewContract(caller, self, value, gasLimit, nil)

	vm.GetStateDB().CreateAccount(contractAddr)
	vm.GetStateDB().SetCode(contractAddr, []byte{123})
	vm.GetStateDB().SetAbi(contractAddr, []byte{123})

	_, _, err := vm.DelegateCall(contract, contractAddr, nil, gasLimit)
	assert.Equal(t, "unexpected EOF", err.Error())

	vm.depth = int(model2.CallCreateDepth + 1)
	_, _, err = vm.DelegateCall(contract, contractAddr, nil, gasLimit)
	assert.Equal(t, g_error.ErrDepth, err)

	vm.depth = 1
	vm.vmConfig.NoRecursion = true
	_, _, err = vm.DelegateCall(contract, contractAddr, nil, gasLimit)
	assert.NoError(t, err)
}

func TestVM_Create_Error(t *testing.T) {
	vm := getTestVm()
	caller := AccountRef(aliceAddr)
	value := g_testData.TestValue
	gasLimit := g_testData.TestGasLimit

	_, _, _, err := vm.Create(caller, nil, gasLimit, value)
	assert.Equal(t, "empty account", err.Error())

	vm.GetStateDB().CreateAccount(caller.Address())
	_, _, _, err = vm.Create(caller, nil, gasLimit, value)
	assert.Equal(t, g_error.ErrInsufficientBalance, err)

	vm.GetStateDB().AddBalance(caller.Address(), big.NewInt(10000))
	_, _, _, err = vm.Create(caller, nil, gasLimit, value)
	assert.Equal(t, errEmptyInput, err)

	input, _ := rlp.EncodeToBytes([]interface{}{"code", "abi"})
	_, _, _, err = vm.Create(caller, input, gasLimit, value)
	assert.Equal(t, "wasm: Invalid magic number", err.Error())

	vm.depth = 1
	vm.vmConfig.NoRecursion = true
	_, _, _, err = vm.Create(caller, input, gasLimit, value)
	assert.NoError(t, err)

	nonce, _ := vm.GetStateDB().GetNonce(caller.Address())
	expectAddr := cs_crypto.CreateContractAddress(caller.Address(), nonce)
	vm.GetStateDB().CreateAccount(expectAddr)
	vm.GetStateDB().AddNonce(expectAddr, uint64(1))
	_, _, _, err = vm.Create(caller, input, gasLimit, value)
	assert.Equal(t, g_error.ErrContractAddressCollision, err)

	vm.depth = int(model2.CallCreateDepth + 1)
	_, _, _, err = vm.Create(caller, nil, gasLimit, value)
	assert.Equal(t, g_error.ErrDepth, err)
}
