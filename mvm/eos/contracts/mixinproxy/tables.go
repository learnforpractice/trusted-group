package main

import "github.com/uuosio/chain"

//table txevents
type TxEvent struct {
	nonce      uint64 //primary : t.nonce
	process    chain.Uint128
	asset      chain.Uint128
	members    []chain.Uint128
	threshold  int32
	amount     chain.Uint128
	extra      []byte
	timestamp  uint64
	signatures []chain.Signature
}

//table errorevents
type ErrorTxEvent struct {
	event  TxEvent //primary : t.event.nonce
	reason string
}

//table submittedevs
type SubmittedEvent struct {
	nonce uint64 //primary : t.nonce
}

//table counters
type Counter struct {
	id    uint64 //primary : t.id
	count uint64
}

//table accountcache singleton
type AccountCache struct {
	id      uint64
	account chain.Name
}

//table bindaccounts
type MixinAccount struct {
	eos_account chain.Name    //primary : t.eos_account.N
	client_id   chain.Uint128 //secondary
}

//table mixinassets
type MixinAsset struct {
	symbol   chain.Symbol  //primary : t.symbol.Code()
	asset_id chain.Uint128 //secondary
}

//table transferfees
type TransferFee struct {
	fee chain.Asset //primary : t.fee.Symbol.Code()
}

//table totalfees
type TotalFee struct {
	total chain.Asset //primary : t.total.Symbol.Code()
}

//table createaccfee singleton
type CreateAccountFee struct {
	fee chain.Asset
}

//table processes ignore
type Process struct {
	contract chain.Name //primary : t.contract.N
	process  chain.Uint128
}
