package state_processor

import (
	"github.com/dipperin/dipperin-core/common"
	"math/big"
)

type Fullstate struct{
	state *AccountStateDB
}

func (f  *Fullstate) CreateAccount(address common.Address) {
	f.state.setBalance(address,big.NewInt(0))
	f.state.setNonce(address,uint64(0))
}

func (f  *Fullstate) SubBalance(addr common.Address, amount *big.Int) {
	panic("can not call sub Balance in vm")
}

func (f  *Fullstate) AddBalance(addr common.Address, amount *big.Int) {
	panic("can not call sub Balance in vm")
}

func (f  *Fullstate) GetBalance(addr common.Address) *big.Int {
	balance, err :=f.state.GetBalance(addr)
	if err != nil{
		return big.NewInt(0)
	}
	return balance
}

func (f  *Fullstate) GetNonce(addr common.Address) uint64 {
	nonce, err := f.state.GetNonce(addr)
	if err != nil{
		return uint64(0)
	}
	return nonce
}

func (f  *Fullstate) SetNonce(common.Address, uint64) {
	panic("implement me")
}

func (f  *Fullstate) GetCodeHash(common.Address) common.Hash {
	panic("implement me")
}

func (f  *Fullstate) GetCode(addr common.Address) (result []byte) {
	//f.state.contractTrieCache
	ct, err := f.state.getContractTrie(addr)
	if err!=nil{
		return
	}
	return ct.GetKey(GetContractFieldKey(addr,"code"))
}

func (f  *Fullstate) SetCode(addr common.Address, code []byte) {
	ct, err := f.state.getContractTrie(addr)
	if err!=nil{
		return
	}
	err = ct.TryUpdate(GetContractFieldKey(addr,"code"),code)
	if err!=nil{
		return
	}
}

func (f  *Fullstate) GetCodeSize(addr common.Address) (size int) {
	ct, err := f.state.getContractTrie(addr)
	if err!=nil{
		return
	}
	code := ct.GetKey(GetContractFieldKey(addr,"code"))
	return len(code)
}

func (f  *Fullstate) GetAbiHash(common.Address) common.Hash {
	panic("implement me")
}

func (f  *Fullstate) GetAbi(addr common.Address) (abi []byte) {
	ct, err := f.state.getContractTrie(addr)
	if err!=nil{
		return
	}
	return ct.GetKey(GetContractFieldKey(addr,"abi"))
}

func (f  *Fullstate) SetAbi(common.Address, []byte) {
	panic("implement me")
}

func (f  *Fullstate) AddRefund(uint64) {
	panic("implement me")
}

func (f  *Fullstate) SubRefund(uint64) {
	panic("implement me")
}

func (f  *Fullstate) GetRefund() uint64 {
	panic("implement me")
}

func (f  *Fullstate) GetCommittedState(common.Address, []byte) []byte {
	panic("implement me")
}

func (f  *Fullstate) GetState(addr common.Address,key []byte) (data []byte) {
	ct, err := f.state.getContractTrie(addr)
	if err!=nil{
		return
	}
	return ct.GetKey(GetContractFieldKey(addr,string(key)))
}

func (f  *Fullstate) SetState(addr common.Address,key []byte, value []byte) {
	ct, err := f.state.getContractTrie(addr)
	if err!=nil{
		return
	}
	err = ct.TryUpdate(GetContractFieldKey(addr,string(key)),value)
	if err!= nil{
		panic("can not update contract field")
	}
}

func (f  *Fullstate) AddLog(address common.Address, topics []common.Hash, data []byte, bn uint64) {
	panic("implement me")
}

func (f  *Fullstate) Suicide(common.Address) bool {
	panic("implement me")
}

func (f  *Fullstate) HasSuicided(common.Address) bool {
	panic("implement me")
}

func (f  *Fullstate) Exist(common.Address) bool {
	panic("implement me")
}

func (f  *Fullstate) Empty(common.Address) bool {
	panic("implement me")
}

func (f  *Fullstate) RevertToSnapshot(int) {
	panic("implement me")
}

func (f  *Fullstate) Snapshot() int {
	panic("implement me")
}

func (f  *Fullstate) AddPreimage(common.Hash, []byte) {
	panic("implement me")
}

func (f  *Fullstate) ForEachStorage(common.Address, func(common.Hash, common.Hash) bool) {
	panic("implement me")
}

func (f  *Fullstate) TxHash() common.Hash {
	panic("implement me")
}

func (f  *Fullstate) TxIdx() uint32 {
	panic("implement me")
}





