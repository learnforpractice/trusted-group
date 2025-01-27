//go:build debug
// +build debug

package main

import (
	"github.com/uuosio/chain"
	"github.com/uuosio/chain/database"
)

const (
	DEBUG            = true
	KEY_COUNTER_TEST = 7
)

func ClearTable(db database.MultiIndexInterface) {
	for {
		it := db.Lowerbound(0)
		if !it.IsOk() {
			break
		}
		db.Remove(it)
	}
}

func ClearSingletonTable(db *database.SingletonTable) {
	db.Remove()
}

//action updateauth
func (c *Contract) UpdateAuth(account chain.Name) {
	chain.RequireAuth(c.self)
	auth := Authority{
		uint32(1),     //threshold
		[]KeyWeight{}, //keys
		[]PermissionLevelWeight{
			PermissionLevelWeight{
				PermissionLevel{
					c.self,
					chain.ActiveName,
				},
				uint16(1),
			},
			PermissionLevelWeight{
				PermissionLevel{
					c.self,
					chain.NewName("multisig"),
				},
				uint16(1),
			},
		}, //accounts
		[]WaitWeight{},
	}
	chain.NewAction(
		&chain.PermissionLevel{account, chain.ActiveName},
		chain.NewName("eosio"),
		chain.NewName("updateauth"),
		account,          //account
		chain.ActiveName, //permission
		chain.OwnerName,  //parent
		&auth,
	).Send()
}

//action clear
func (c *Contract) clear() {
	chain.RequireAuth(c.self)
	ClearSingletonTable(NewAccountCacheTable(c.self, c.self).db)
	ClearTable(NewMixinAccountTable(c.self, c.self).MultiIndexInterface)
	// return
	ClearTable(NewCounterTable(c.self, c.self).MultiIndexInterface)

	// ClearTable(NewMixinAssetTable(c.self, c.self).MultiIndexInterface)

	ClearTable(NewTxEventTable(c.self, c.self).MultiIndexInterface)
	// ClearTable(NewMTGWorkTable(c.self, c.self).MultiIndexInterface)
}

//action test
func (c *Contract) test() {
	c.GetNextIndex(KEY_COUNTER_TEST, 1)
}
