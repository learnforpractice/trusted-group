package main

import (
	"context"
	"testing"

	"github.com/uuosio/chaintester"
)

var ctx = context.Background()

func OnApply(receiver, firstReceiver, action uint64) {
	contract_apply(receiver, firstReceiver, action)
}

func init() {
	chaintester.SetApplyFunc(OnApply)
}

func TestHello(t *testing.T) {
	tester := chaintester.NewChainTester()
	defer tester.FreeChain()

	tester.ProduceBlock()
}
