package state_processor

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/accounts"
	"github.com/dipperin/dipperin-core/core/accounts/soft-wallet"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/vm/common/utils"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/big"
	"testing"
)

func TestAccountStateDB_ProcessContract(t *testing.T) {
	ownAddress := common.HexToAddress("0x000062be10f46b5d01Ecd9b502c4bA3d6131f6fc2e41")
	log.InitLogger(log.LvlDebug)
	transactionJson := "{\"TxData\":{\"nonce\":\"0x0\",\"to\":\"0x00120000000000000000000000000000000000000000\",\"hashlock\":null,\"timelock\":\"0x0\",\"value\":\"0x2540be400\",\"fee\":\"0x69db9c0\",\"gasPrice\":\"0xa\",\"gas\":\"0x1027127dc00\",\"input\":\"0xf9026db8eb0061736d01000000010d0360017f0060027f7f00600000021d0203656e76067072696e7473000003656e76087072696e74735f6c00010304030202000405017001010105030100020615037f01419088040b7f00419088040b7f004186080b073405066d656d6f727902000b5f5f686561705f6261736503010a5f5f646174615f656e64030204696e697400030568656c6c6f00040a450302000b02000b3d01017f230041106b220124004180081000200141203a000f2001410f6a41011001200010002001410a3a000e2001410e6a41011001200141106a24000b0b0d01004180080b0668656c6c6f00b9017d5b0a202020207b0a2020202020202020226e616d65223a2022696e6974222c0a202020202020202022696e70757473223a205b5d2c0a2020202020202020226f757470757473223a205b5d2c0a202020202020202022636f6e7374616e74223a202266616c7365222c0a20202020202020202274797065223a202266756e6374696f6e220a202020207d2c0a202020207b0a2020202020202020226e616d65223a202268656c6c6f222c0a202020202020202022696e70757473223a205b0a2020202020202020202020207b0a20202020202020202020202020202020226e616d65223a20226e616d65222c0a202020202020202020202020202020202274797065223a2022737472696e67220a2020202020202020202020207d0a20202020202020205d2c0a2020202020202020226f757470757473223a205b5d2c0a202020202020202022636f6e7374616e74223a202274727565222c0a20202020202020202274797065223a202266756e6374696f6e220a202020207d0a5d0a\"},\"Wit\":{\"r\":\"0xa1509f3efb1e632643c9972b9183234445c539a1b483ad0ea4b36a4edabf8d04\",\"s\":\"0xa7a16d72b826aea44e8f56247abbad367cf7e300d564949e66ac97098b9f234\",\"v\":\"0x39\",\"hashkey\":\"0x\"}}"

	var tx model.Transaction
	err := json.Unmarshal([]byte(transactionJson), &tx)
	if err != nil {
		log.Info("TestAccountStateDB_ProcessContract", "err", err)
	}
	log.Info("processContract", "Tx", tx)

	tx.PaddingTxIndex(0)
	gasLimit := gasLimit * 10000000000
	block := CreateBlock(1, common.Hash{}, []*model.Transaction{&tx}, gasLimit)

	db, root := CreateTestStateDB()
	processor, err := NewAccountStateDB(root, NewStateStorageWithCache(db))
	assert.NoError(t, err)

	processor.NewAccountState(ownAddress)
	err = processor.AddNonce(ownAddress, 0)
	processor.AddBalance(ownAddress, new(big.Int).SetInt64(int64(1000000000000000000)))

	assert.NoError(t, err)
	balance, err := processor.GetBalance(ownAddress)
	nonce, err := processor.GetNonce(ownAddress)
	log.Info("balance", "balance", balance.String())
	//log.Info("nonce", "nonce", nonce, "tx.nonce", tx.Nonce())

	log.Info("gasLimit", "gasLimit", gasLimit)

	gasUsed := uint64(0)
	txConfigCreate := &TxProcessConfig{
		Tx:       &tx,
		Header:   block.Header().(*model.Header),
		GetHash:  fakeGetBlockHash,
		GasLimit: &gasLimit,
		GasUsed:  &gasUsed,
	}

	receipt, err := processor.ProcessContract(txConfigCreate, true)
	assert.NoError(t, err)
	log.Info("result", "receipt", receipt)
	assert.Equal(t, false, receipt.HandlerResult)
	tx.PaddingReceipt(receipt)
	receiptResult, err := tx.GetReceipt()
	assert.NoError(t, err)
	contractNonce, err := processor.GetNonce(receiptResult.ContractAddress)
	log.Info("TestAccountStateDB_ProcessContract", "contractNonce", contractNonce, "receiptResult", receiptResult)
	code, err := processor.GetCode(receiptResult.ContractAddress)
	abi, err := processor.GetAbi(receiptResult.ContractAddress)
	log.Info("TestAccountStateDB_ProcessContract", "code  get from state", code)
	assert.NoError(t, err)
	//assert.Equal(t, code, tx.ExtraData())

	sw, err := soft_wallet.NewSoftWallet()
	err = sw.Open(util.HomeDir()+"/go/src/github.com/dipperin/dipperin-core/core/vm/event/CSWallet", "CSWallet", "123")
	assert.NoError(t, err)

	callTx, err := newContractCallTx(nil, &receiptResult.ContractAddress, new(big.Int).SetUint64(1), uint64(1500000), "hello", "name", nonce+1, abi)
	account := accounts.Account{ownAddress}
	signCallTx, err := sw.SignTx(account, callTx, nil)

	assert.NoError(t, err)
	callTx.PaddingTxIndex(0)
	block2 := CreateBlock(2, common.Hash{}, []*model.Transaction{signCallTx}, gasLimit)
	log.Info("callTx info", "callTx", callTx)

	gasUsed2 := uint64(0)
	txConfig := &TxProcessConfig{
		Tx:       signCallTx,
		Header:   block2.Header().(*model.Header),
		GetHash:  fakeGetBlockHash,
		GasLimit: &gasLimit,
		GasUsed:  &gasUsed2,
	}

	callRecipt, err := processor.ProcessContract(txConfig, false)
	//assert.NoError(t, err)
	log.Info("TestAccountStateDB_ProcessContract++", "callRecipt", callRecipt, "err", err)

}

func newContractCallTx(from *common.Address, to *common.Address, gasPrice *big.Int, gasLimit uint64, funcName string, input string, nonce uint64, code []byte) (tx *model.Transaction, err error) {
	// RLP([funcName][params])
	inputRlp, err := rlp.EncodeToBytes([]interface{}{
		funcName, input,
	})
	if err != nil {
		log.Error("input rlp err")
		return
	}

	extraData, err := utils.ParseCallContractData(code, inputRlp)

	if err != nil {
		log.Error("ParseCallContractData  inputRlp", "err", err)
		return
	}

	tx = model.NewTransactionSc(nonce, to, nil, gasPrice, gasLimit, extraData)
	return tx, nil
}

func TestAccountStateDB_ProcessContract2(t *testing.T) {
	var testPath = "../../vm/event"
	tx1 := createContractTx(t, testPath+"/event.wasm", testPath+"/event.cpp.abi.json")
	contractAddr := cs_crypto.CreateContractAddress(aliceAddr, 0)
	name := []byte("ProcessContract")
	num := utils.Int64ToBytes(456)
	param := [][]byte{name, num}
	tx2 := callContractTx(t, &contractAddr, "hello", param, 1)

	db, root := CreateTestStateDB()
	tdb := NewStateStorageWithCache(db)
	processor, err := NewAccountStateDB(root, tdb)
	assert.NoError(t, err)

	block := CreateBlock(1, common.Hash{}, []*model.Transaction{tx1, tx2}, 5*gasLimit)
	tmpGasLimit := block.GasLimit()
	gasUsed := block.GasUsed()
	config := &TxProcessConfig{
		Tx:       tx1,
		Header:   block.Header(),
		GetHash:  getTestHashFunc(),
		GasLimit: &tmpGasLimit,
		GasUsed:  &gasUsed,
		TxFee:    big.NewInt(0),
	}
	err = processor.ProcessTxNew(config)
	assert.NoError(t, err)

	receipt1, err := tx1.GetReceipt()
	assert.NoError(t, err)
	assert.Equal(t, tx1.CalTxId(), receipt1.TxHash)
	assert.Equal(t, cs_crypto.CreateContractAddress(aliceAddr, 0), receipt1.ContractAddress)
	assert.Len(t, receipt1.Logs, 0)

	fmt.Println("---------------------------")

	config.Tx = tx2
	err = processor.ProcessTxNew(config)
	assert.NoError(t, err)

	receipt2, err := tx2.GetReceipt()
	assert.NoError(t, err)
	assert.Equal(t, tx2.CalTxId(), receipt2.TxHash)
	assert.Equal(t, receipt1.ContractAddress, receipt2.ContractAddress)
	assert.Len(t, receipt2.Logs, 1)

	log1 := receipt2.Logs[0]
	assert.Equal(t, tx2.CalTxId(), log1.TxHash)
	assert.Equal(t, common.Hash{}, log1.BlockHash)
	assert.Equal(t, receipt2.ContractAddress, log1.Address)
	assert.Equal(t, uint64(1), log1.BlockNumber)
}

func TestAccountStateDB_ProcessContractToken(t *testing.T) {

	singer := model.NewMercurySigner(new(big.Int).SetInt64(int64(1)))

	ownSK, _ := crypto.GenerateKey()
	ownPk := ownSK.PublicKey
	ownAddress := cs_crypto.GetNormalAddress(ownPk)

	aliceSK, _ := crypto.GenerateKey()
	alicePk := aliceSK.PublicKey
	aliceAddress := cs_crypto.GetNormalAddress(alicePk)

	brotherSK, _ := crypto.GenerateKey()
	brotherPk := brotherSK.PublicKey
	brotherAddress := cs_crypto.GetNormalAddress(brotherPk)

	addressSlice := []common.Address{
		ownAddress,
		aliceAddress,
		brotherAddress,
	}

	abiPath := "../../vm/event/token/token.cpp.abi.json"
	wasmPath := "../../vm/event/token/token-wh.wasm"
	err, data := utils.GetExtraData(abiPath, wasmPath, []string{"dipp", "DIPP", "100000000"})
	assert.NoError(t, err)

	addr := common.HexToAddress(common.AddressContractCreate)

	tx := model.NewTransactionSc(0, &addr, new(big.Int).SetUint64(uint64(10)), new(big.Int).SetUint64(uint64(1)), 26427000, data)

	signCreateTx := getSignedTx(t, ownSK, tx, singer)

	signCreateTx.PaddingTxIndex(0)

	gasLimit := gasLimit * 10000000000
	block := CreateBlock(1, common.Hash{}, []*model.Transaction{signCreateTx}, gasLimit)

	processor, err := CreateProcessorAndInitAccount(err, t, addressSlice)

	tmpGasLimit := block.GasLimit()
	gasUsed := block.GasUsed()
	config := &TxProcessConfig{
		Tx:       tx,
		Header:   block.Header(),
		GetHash:  getTestHashFunc(),
		GasLimit: &tmpGasLimit,
		GasUsed:  &gasUsed,
		TxFee:    big.NewInt(0),
	}

	err = processor.ProcessTxNew(config)
	assert.NoError(t, err)

	receipt, err := tx.GetReceipt()

	assert.NoError(t, err)
	//assert.Equal(t, receipt.Status)

	//byteBalance := []byte{7, 98, 97, 108, 97, 110, 99, 101}
	//baData := processor.GetData(receipt.ContractAddress, string(byteBalance))
	//fmt.Println("&&&&&", receipt.ContractAddress, baData, processor.smartContractData)

	contractNonce, err := processor.GetNonce(receipt.ContractAddress)
	log.Info("TestAccountStateDB_ProcessContract", "contractNonce", contractNonce, "receiptResult", receipt)
	code, err := processor.GetCode(receipt.ContractAddress)
	abi,err := processor.GetAbi(receipt.ContractAddress)
	log.Info("TestAccountStateDB_ProcessContract", "code  get from state", code)
	assert.NoError(t, err)
	//assert.Equal(t, code, tx.ExtraData())
	processor.Commit()

	accountOwn := accounts.Account{ownAddress}
	//  合约调用getBalance方法  获取合约原始账户balance
	ownTransferNonce, err := processor.GetNonce(ownAddress)
	assert.NoError(t, err)
	err = processContractCall(t, receipt.ContractAddress, abi, ownSK,  processor, accountOwn, ownTransferNonce, "getBalance", ownAddress.Hex(), 2, singer)
	assert.NoError(t, err)

	gasUsed2 := uint64(0)
	//  合约调用  transfer方法 转账给alice
	ownTransferNonce++
	err = processContractCall(t, receipt.ContractAddress, abi, ownSK,  processor, accountOwn, ownTransferNonce, "transfer", aliceAddress.Hex()+",20", 3, singer)
	assert.NoError(t, err)

	//  合约调用getBalance方法  获取alice账户balance
	ownTransferNonce++
	err = processContractCall(t, receipt.ContractAddress, abi, ownSK,  processor, accountOwn, ownTransferNonce, "getBalance", aliceAddress.Hex(), 4, singer)
	assert.NoError(t, err)

	//  合约调用approve方法
	log.Info("==========================================")
	ownTransferNonce++
	callTxApprove, err := newContractCallTx(nil, &receipt.ContractAddress, new(big.Int).SetUint64(1), uint64(1500000), "approve", brotherAddress.Hex()+",50", ownTransferNonce, abi)
	//accountAlice := accounts.Account{aliceAddress}
	signCallTxApprove, err := callTxApprove.SignTx(ownSK, singer)

	assert.NoError(t, err)
	signCallTxApprove.PaddingTxIndex(0)
	block5 := CreateBlock(5, common.Hash{}, []*model.Transaction{signCallTxApprove}, gasLimit)
	log.Info("signCallTxApprove info", "signCallTxApprove", signCallTxApprove)

	txConfig5 := &TxProcessConfig{
		Tx:       signCallTxApprove,
		Header:   block5.Header().(*model.Header),
		GetHash:  fakeGetBlockHash,
		GasLimit: &gasLimit,
		GasUsed:  &gasUsed2,
	}

	err = processor.ProcessTxNew(txConfig5)
	assert.NoError(t, err)
	processor.Commit()


	//  合约调用getApproveBalance方法  获取own授权给brother账户balance
	/*err = processContractCall(t, receipt.ContractAddress, abi, ownSK,  processor, accountOwn, 5, "getApproveBalance", "0x000062be10f46b5d01Ecd9b502c4bA3d6131f6fc2e41,0x00004179D57e45Cb3b54D6FAEF69e746bf240E287978", 6)
	assert.NoError(t, err)*/




	//  合约调用transferFrom方法
	log.Info("==========================================")
	callTxTransferFrom, err := newContractCallTx(nil, &receipt.ContractAddress, new(big.Int).SetUint64(1), uint64(1500000), "transferFrom", ownAddress.Hex()+","+aliceAddress.Hex() + ",5", 0, abi)
	assert.NoError(t, err)
	accountBrother := accounts.Account{Address: brotherAddress}
	assert.NoError(t, err)

	signCallTxTransferFrom, err := callTxTransferFrom.SignTx(brotherSK, singer)
	assert.NoError(t, err)
	signCallTxTransferFrom.PaddingTxIndex(0)
	block7 := CreateBlock(7, common.Hash{}, []*model.Transaction{signCallTxTransferFrom}, gasLimit)
	log.Info("signCallTxTransferFrom info", "signCallTxTransferFrom", signCallTxTransferFrom)

	txConfig7 := &TxProcessConfig{
		Tx:       signCallTxTransferFrom,
		Header:   block7.Header().(*model.Header),
		GetHash:  fakeGetBlockHash,
		GasLimit: &gasLimit,
		GasUsed:  &gasUsed2,
	}

	err = processor.ProcessTxNew(txConfig7)
	assert.NoError(t, err)
	processor.Commit()


	//  合约调用getBalance方法  获取alice账户获得转账授权后的balance
	ownTransferNonce++
	err = processContractCall(t, receipt.ContractAddress, abi, ownSK,  processor, accountOwn, ownTransferNonce, "getBalance", aliceAddress.Hex(), 8, singer)
	assert.NoError(t, err)

	//  合约调用getBalance方法  获取own账户最终的balance
	ownTransferNonce++
	err = processContractCall(t, receipt.ContractAddress, abi, ownSK,  processor, accountOwn, ownTransferNonce, "getBalance", ownAddress.Hex(), 9, singer)
	assert.NoError(t, err)

	// 合约调用  transfer方法  转账给brother
	ownTransferNonce++
	err = processContractCall(t, receipt.ContractAddress, abi, ownSK,  processor, accountOwn, ownTransferNonce, "transfer", brotherAddress.Hex()+",28", 10, singer)
	assert.NoError(t, err)

	//  合约调用getBalance方法  获取own账户最终的balance
	ownTransferNonce++
	err = processContractCall(t, receipt.ContractAddress, abi, ownSK,  processor, accountOwn, ownTransferNonce, "getBalance", ownAddress.Hex(), 11, singer)
	assert.NoError(t, err)

    // 合约调用burn方法,将账户余额返还给own
	err = processContractCall(t, receipt.ContractAddress, abi, brotherSK,  processor, accountBrother, 1, "burn", "15", 12, singer)
	assert.NoError(t, err)



	// 合约调用getBalance方法,获取own的余额
	ownTransferNonce++
	err = processContractCall(t, receipt.ContractAddress, abi, ownSK,  processor, accountOwn, ownTransferNonce, "getBalance", ownAddress.Hex(), 13, singer)
	assert.NoError(t, err)


	// 合约调用setName方法，设置合约名
	ownTransferNonce++
	err = processContractCall(t, receipt.ContractAddress, abi, ownSK,  processor, accountOwn, ownTransferNonce, "setName", "wujinhai", 14, singer)
	assert.NoError(t, err)



	log.Info("TestAccountStateDB_ProcessContract++", "callRecipt", "", "err", err)
}


//  合约调用getBalance方法
func processContractCall(t *testing.T,contractAddress common.Address, code []byte, priKey *ecdsa.PrivateKey, processor *AccountStateDB, accountOwn accounts.Account, nonce uint64, funcName string, params string, blockNum uint64, singer model.Signer) ( error) {
	gasUsed2 := uint64(0)
	gasLimit := gasLimit * 10000000000
	log.Info("processContractCall=================================================")
	callTx, err := newContractCallTx(nil, &contractAddress, new(big.Int).SetUint64(1), uint64(1500000), funcName, params, nonce, code)
	assert.NoError(t,err)
	signCallTx, err := callTx.SignTx(priKey, singer)
	//sw.SignTx(accountOwn, callTx, nil)
	assert.NoError(t, err)
	callTx.PaddingTxIndex(0)
	block := CreateBlock(blockNum, common.Hash{}, []*model.Transaction{signCallTx}, gasLimit)
	log.Info("callTx info", "callTx", callTx)
	txConfig := &TxProcessConfig{
		Tx:       signCallTx,
		Header:   block.Header().(*model.Header),
		GetHash:  fakeGetBlockHash,
		GasLimit: &gasLimit,
		GasUsed:  &gasUsed2,
	}
	err = processor.ProcessTxNew(txConfig)
	if funcName == "getBalance" {
		receipt,err := callTx.GetReceipt()
		assert.NoError(t,err)
		log.Info("receipt  log", "receipt log", receipt.Logs)
	}
	assert.NoError(t, err)
	processor.Commit()
	return err
}

func CreateProcessorAndInitAccount(err error, t *testing.T, addressSlice []common.Address) (*AccountStateDB, error) {
	db, root := CreateTestStateDB()
	processor, err := NewAccountStateDB(root, NewStateStorageWithCache(db))
	assert.NoError(t, err)
	processor.NewAccountState(addressSlice[0])
	err = processor.AddNonce(addressSlice[0], 0)
	processor.AddBalance(addressSlice[0], new(big.Int).SetInt64(int64(1000000000000000000)))
	for i := 1; i < len(addressSlice); i++ {
		fmt.Println("xxxxxxxxxxxxxxxxx", addressSlice[i])
		processor.NewAccountState(addressSlice[i])
		err = processor.AddNonce(addressSlice[i], 0)
		processor.AddBalance(addressSlice[i], new(big.Int).SetInt64(int64(10000000)))

	}
	return processor, err
}

func TestGetByteFromAbiFile(t *testing.T) {
	bytes, err := ioutil.ReadFile("../../vm/event/example.cpp.abi.json")
	assert.NoError(t, err)
	fmt.Println(bytes)
}

func getSignedTx(t *testing.T,priKey *ecdsa.PrivateKey, tx *model.Transaction, singer model.Signer) *model.Transaction {
	signCreateTx, err := tx.SignTx(priKey, singer)
	assert.NoError(t, err)
	return signCreateTx
}

