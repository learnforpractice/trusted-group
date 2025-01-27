package main

import (
	"context"
	"crypto/ed25519"
	"math/big"

	"github.com/MixinNetwork/mixin/crypto"
	"github.com/MixinNetwork/mixin/domains/ethereum"
	"github.com/MixinNetwork/mixin/logger"
	"github.com/fox-one/mixin-sdk-go"
)

type User struct {
	*mixin.User
	Key      *mixin.Keystore `json:"key"`
	PIN      string          `json:"-"`
	Contract string          `json:"contract"`
}

// TODO should verify the signature from MetaMask of the addr
func (p *Proxy) createUser(ctx context.Context, store *Storage, addr, sig string) (*User, error) {
	err := ethereum.VerifyAddress(addr)
	if err != nil {
		return nil, err
	}

	old, err := store.readUserByAddress(addr)
	if err != nil {
		return nil, err
	}
	if old != nil && old.Contract != "" {
		return old, nil
	}

	seed := crypto.NewHash([]byte(ProxyUserSecret + addr))
	signer := ed25519.NewKeyFromSeed(seed[:])
	u, ks, err := p.CreateUser(ctx, signer, addr)
	if err != nil {
		return nil, err
	}
	user := &User{u, ks, "", ""}

	seed = crypto.NewHash(seed[:])
	pin := new(big.Int).SetBytes(seed[:]).String()
	for len(pin) < 6 {
		pin = pin + pin
	}
	user.PIN = pin[:6]

	if !user.HasPin {
		uc, err := mixin.NewFromKeystore(ks)
		if err != nil {
			return nil, err
		}
		err = uc.ModifyPin(ctx, "", user.PIN)
		if err != nil {
			return nil, err
		}
	}

	err = store.writeUser(user)
	if err != nil {
		return nil, err
	}
	return p.readUserWithContract(store, user.UserID)
}

func (p *Proxy) readUserWithContract(store *Storage, id string) (*User, error) {
	user, err := store.readUserById(id)
	if err != nil || user == nil {
		return nil, err
	}
	if user.Contract != "" {
		return user, nil
	}
	ua, err := user.getContract(p)
	if err != nil {
		return nil, err
	}
	if ua.String() == "0x0000000000000000000000000000000000000000" {
		return user, nil
	}
	user.Contract = ua.String()
	err = store.writeUser(user)
	return user, err
}

func (u *User) handle(ctx context.Context, store *Storage, s *mixin.Snapshot, act *Action) error {
	logger.Verbosef("User.handle(%v, %v)", *s, *act)
	if act.Destination != "" {
		return u.submit(ctx, store, s, act)
	}

	traceId := mixin.UniqueConversationID(s.SnapshotID, "HANDLE||TRANSFER")
	input := &mixin.TransferInput{
		AssetID: s.AssetID,
		Amount:  s.Amount,
		TraceID: traceId,
		Memo:    act.Extra,
	}
	if len(act.Receivers) == 1 {
		input.OpponentID = act.Receivers[0]
	} else {
		input.OpponentMultisig.Receivers = act.Receivers
		input.OpponentMultisig.Threshold = uint8(act.Threshold)
	}
	return u.send(ctx, input)
}

func (u *User) pass(ctx context.Context, p *Proxy, s *mixin.Snapshot) error {
	logger.Verbosef("User.pass(%v)", *s)
	return u.bindAndPass(ctx, p, s.SnapshotID, u.FullName, s.AssetID, s.Amount)
}

func (u *User) send(ctx context.Context, in *mixin.TransferInput) error {
	uc, err := mixin.NewFromKeystore(u.Key)
	if err != nil {
		panic(err)
	}
	if len(in.OpponentMultisig.Receivers) > 0 {
		_, err = uc.Transaction(ctx, in, u.PIN)
	} else {
		_, err = uc.Transfer(ctx, in, u.PIN)
	}
	logger.Verbosef("User.send(%v) => %v", *in, err)
	return err
}
