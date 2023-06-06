package main

import (
	"github.com/uuosio/chain"
)

const (
	KEY_TX_REQUEST_SEQ = 1
)

// table processes
type Process struct {
	contract chain.Name //primary : t.contract.N
	process  chain.Uint128
}

// table logs
type TxLog struct {
	id        uint64 //primary : t.id
	nonce     uint64
	contract  chain.Name
	process   chain.Uint128
	asset     chain.Uint128
	members   []chain.Uint128
	threshold int32
	amount    chain.Uint128
	extra     []byte
	timestamp uint64
}

// table signers singleton
type Signers struct {
	public_keys []chain.PublicKey
}

// table counters
type Counter struct {
	id    uint64 //primary : t.id
	count uint64
}

// contract mtg.xin
type Contract struct {
	self, firstReceiver, action chain.Name
}

func NewContract(receiver, firstReceiver, action chain.Name) *Contract {
	return &Contract{receiver, firstReceiver, action}
}

// action setup
func (c *Contract) Setup(signers []chain.PublicKey) {
	chain.RequireAuth(c.self)

	check(!hasDuplicates(signers), "invalid signers")

	db := NewSignersTable(c.self, c.self)
	db.Set(&Signers{public_keys: signers}, c.self)
}

// action addprocess
func (c *Contract) AddProcess(contract chain.Name, process chain.Uint128, signatures []chain.Signature) {
	check(chain.IsAccount(contract), "contract account does not exists!")

	enc := chain.NewEncoder(8 + 16)
	enc.PackName(contract)
	enc.Pack(&process)
	data := enc.GetBytes()
	VerifySignatures(c.self, data, signatures)

	db := NewProcessTable(c.self, c.self)
	it := db.Find(contract.N)
	check(!it.IsOk(), "process already exists!")
	item := &Process{
		contract: contract,
		process:  process,
	}
	db.Store(item, c.self)
}

// action txrequest
func (c *Contract) TxRequest(nonce uint64,
	contract chain.Name,
	process chain.Uint128,
	asset chain.Uint128,
	members []chain.Uint128,
	threshold int32,
	amount chain.Uint128,
	extra []byte) {

	chain.RequireAuth(contract)
	db := NewProcessTable(c.self, c.self)
	it, item := db.GetByKey(contract.N)
	check(it.IsOk(), "process not found!")
	check(item.process == process, "invalid process!")

	seq := c.GetNextSeq()
	log := TxLog{
		id:        seq,
		nonce:     nonce,
		contract:  contract,
		process:   process,
		asset:     asset,
		members:   members,
		threshold: threshold,
		amount:    amount,
		extra:     extra,
		timestamp: chain.CurrentTime().Elapsed * 1000,
	}

	chain.NewAction(
		&chain.PermissionLevel{c.self, chain.ActiveName},
		c.self,
		chain.NewName("ontxlog"),
		&log,
	).Send()
	//TODO: emit transfer event so block explorer can show it
}

// action ontxlog
func (c *Contract) OnTxLog(log *TxLog) {
	chain.RequireAuth(c.self)
}

func (c *Contract) GetNextIndex(key uint64) uint64 {
	db := NewCounterTable(c.self, c.self)
	if it, item := db.GetByKey(key); it.IsOk() {
		index := item.count
		item.count += 1
		db.Update(it, item, chain.SamePayer)
		return index
	} else {
		item := Counter{id: key, count: 1}
		db.Store(&item, c.self)
		return 0
	}
}

func (c *Contract) GetNextSeq() uint64 {
	return c.GetNextIndex(KEY_TX_REQUEST_SEQ)
}
