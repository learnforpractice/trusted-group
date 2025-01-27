package machine

import (
	"context"
	"encoding/base64"

	"github.com/MixinNetwork/mixin/common"
	"github.com/MixinNetwork/mixin/logger"
	"github.com/MixinNetwork/nfo/mtg"
	"github.com/MixinNetwork/trusted-group/mvm/encoding"
	"github.com/fox-one/mixin-sdk-go"
)

type Asset struct {
	Id     string
	Symbol string
	Name   string
}

var (
	// because the sdk bug, this output is skipped, and should always be in the future
	InvalidCollectibleOutputHackMap = map[string]bool{
		"271d7ef5-6bf3-3b96-9c0c-701f7a989435": true,
		"8f96c027-fbf0-39dc-99b7-6ba6cdf9c66c": true,
		"5bb0997f-669a-3e65-8c2f-5cc15de4c9ca": true,
		"175620b8-f4b8-3e66-84bb-669c62d5140d": true,
	}
)

func (m *Machine) ProcessOutput(ctx context.Context, out *mtg.Output) {
	op, err := parseOperation(out.Memo)
	logger.Verbosef("Machine.ProcessOutput(%v) => %v %v", out, op, err)
	if err != nil {
		return
	}
	switch op.Purpose {
	case encoding.OperationPurposeAddProcess:
		m.AddProcess(ctx, op.Process, op.Platform, op.Address, out, op.Extra)
	case encoding.OperationPurposeGroupEvent:
		m.WriteGroupEvent(ctx, op.Process, out, op.Extra)
	}
}

func (m *Machine) ProcessCollectibleOutput(ctx context.Context, out *mtg.CollectibleOutput) {
	if InvalidCollectibleOutputHackMap[out.OutputId] {
		return
	}
	op, err := parseOperation(out.Memo)
	logger.Verbosef("Machine.ProcessCollectibleOutput(%v) => %v %v", out, op, err)
	if err != nil {
		return
	}
	switch op.Purpose {
	case encoding.OperationPurposeGroupEvent:
		m.WriteNFOGroupEvent(ctx, op.Process, out, op.Extra)
	}
}

func parseOperation(memo string) (*encoding.Operation, error) {
	b, err := base64.RawURLEncoding.DecodeString(memo)
	if err != nil {
		return nil, err
	}
	return encoding.DecodeOperation(b)
}

func (m *Machine) checkAssetOrCollectible(ctx context.Context, id string) (string, error) {
	cat, err := m.store.ReadAssetOrCollectible(id)
	if err != nil || cat != "" {
		return cat, err
	}

	asset, err := m.fetchAssetMeta(ctx, id)
	if mixin.IsErrorCodes(err, 404, 10002) {
		err = nil
	}
	if err != nil {
		return "", err
	} else if asset != nil {
		return "ASSET", m.store.WriteAssetOrCollectible(id, "ASSET")
	}

	token, err := m.fetchCollectibleToken(ctx, id)
	if err != nil {
		return "", err
	} else if token != nil {
		return "COLLECTIBLE", m.store.WriteAssetOrCollectible(id, "COLLECTIBLE")
	}

	panic(id)
}

func (m *Machine) fetchAssetMeta(ctx context.Context, id string) ([]byte, error) {
	old, err := m.store.ReadAsset(id)
	if err != nil {
		return nil, err
	} else if old != nil {
		return encodeAssetMeta(old.Symbol, old.Name), nil
	}
	asset, err := m.mixin.ReadAsset(ctx, id)
	if err != nil || asset == nil {
		return nil, err
	}
	err = m.store.WriteAsset(&Asset{
		Id:     id,
		Symbol: asset.Symbol,
		Name:   asset.Name,
	})
	return encodeAssetMeta(asset.Symbol, asset.Name), err
}

func encodeAssetMeta(symbol, name string) []byte {
	enc := common.NewEncoder()
	enc.WriteInt(len(symbol))
	enc.Write([]byte(symbol))
	enc.WriteInt(len(name))
	enc.Write([]byte(name))
	return enc.Bytes()
}
