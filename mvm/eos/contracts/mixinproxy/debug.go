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

var (
	MTG_XIN         = chain.NewName("mtgxinmtgxin")
	MIXIN_WTOKENS   = chain.NewName("mixinwtokens")
	ACCOUNT_OWNER   = chain.NewName("mtgxinmtgxin")
	ACCOUNT_CREATER = chain.NewName("mixincrossss")
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

// action updateauth
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

// action clear
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

// action test
func (c *Contract) test() {
	c.GetNextIndex(KEY_COUNTER_TEST, 1)
}

// action testname
func (c *Contract) testName() {
	// aaaaaaaaamvm
	// name := GetAccountNameFromId((uint64)i)
	name := GetAccountNameFromId(uint64(30))
	chain.Check(name == chain.NewName("aaaaaaaa5mvm"), "bad value")

	name = GetAccountNameFromId(uint64(31))
	chain.Check(name == chain.NewName("aaaaaaabamvm"), "bad value")

	// for i := 0; i <= 100; i++ {
	// 	name := GetAccountNameFromId(uint64(i))
	// 	chain.Println("++++++++:", i, name)
	// }
}