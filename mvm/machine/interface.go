package machine

import (
	"github.com/MixinNetwork/mixin/common"
	"github.com/MixinNetwork/trusted-group/mvm/encoding"
)

type Store interface {
	WriteGroupEventAndNonce(pid string, event *encoding.Event) error
	ListSignedGroupEvents(pid string, limit int) ([]*encoding.Event, error)
	ExpireGroupEventsWithCost(events []*encoding.Event, cost common.Integer) error
	ListPendingGroupEvents(limit int) ([]*encoding.Event, error)
	WriteGroupEventState(pid string, nonce uint64) error

	ReadAccount(pid string, asset string) (*Account, error)
	WriteAccountChange(pid string, asset string, amount common.Integer, credit bool) error

	ReadEngineGroupEventsOffset(pid string) (uint64, error)
	WriteEngineGroupEventsOffset(pid string, offset uint64) error

	ListProcesses() ([]*Process, error)
	WriteProcess(p *Process) error
}

type Engine interface {
	VerifyAddress(addr string) error
	SetupNotifier(addr string) error
	EstimateCost(events []*encoding.Event) (common.Integer, error)
	EnsureSendGroupEvents(address string, events []*encoding.Event) error
	ReceiveGroupEvents(address string, offset uint64, limit int) ([]*encoding.Event, error)
}