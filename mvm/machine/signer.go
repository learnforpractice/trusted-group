package machine

import (
	"bytes"
	"context"
	"encoding/hex"
	"time"

	"github.com/MixinNetwork/tip/crypto"
	"github.com/MixinNetwork/tip/logger"
	"github.com/MixinNetwork/trusted-group/mvm/encoding"
	"github.com/drand/kyber/pairing/bn256"
	"github.com/drand/kyber/share"
	"github.com/drand/kyber/sign/tbls"
)

func (m *Machine) loopSignGroupEvents(ctx context.Context) {
	for {
		events, err := m.store.ListPendingGroupEvents(100)
		if err != nil {
			panic(err)
		}
		if len(events) == 0 {
			time.Sleep(5 * time.Second)
			continue
		}
		for _, e := range events {
			partials, err := m.store.ReadPendingGroupEventSignatures(e.Process, e.Nonce)
			if err != nil {
				panic(err)
			}
			if len(partials) >= m.group.GetThreshold() {
				e.Signature = m.recoverSignature(e, partials)
				err = m.store.WriteSignedGroupEventAndExpirePending(e)
				if err != nil {
					panic(err)
				}
				continue
			}

			scheme := tbls.NewThresholdSchemeOnG1(bn256.NewSuiteG2())
			partial, err := scheme.Sign(m.share, e.Encode())
			if err != nil {
				panic(err)
			}
			if checkSignedWith(partials, partial) {
				continue
			}

			e.Signature = partial
			err = m.messenger.SendMessage(ctx, e.Encode())
			if err != nil {
				panic(err)
			}
			partials = append(partials, partial)
			err = m.store.WritePendingGroupEventSignatures(e.Process, e.Nonce, partials)
			if err != nil {
				panic(err)
			}
		}
	}
}

func (m *Machine) loopReceiveGroupMessages(ctx context.Context) {
	for {
		b, err := m.messenger.ReceiveMessage(ctx)
		if err != nil {
			logger.Verbosef("ReceiveMessage() => %s", err)
			panic(err)
		}
		evt, err := encoding.DecodeEvent(b)
		if err != nil {
			logger.Verbosef("DecodeEvent(%s) => %s", hex.EncodeToString(b), err)
			continue
		}
		partials, err := m.store.ReadPendingGroupEventSignatures(evt.Process, evt.Nonce)
		if err != nil {
			panic(err)
		}
		if checkSignedWith(partials, evt.Signature) {
			continue
		}
		partials = append(partials, evt.Signature)
		err = m.store.WritePendingGroupEventSignatures(evt.Process, evt.Nonce, partials)
		if err != nil {
			panic(err)
		}
	}
}

func (m *Machine) recoverSignature(e *encoding.Event, partials [][]byte) []byte {
	e.Signature = nil
	msg := e.Encode()
	suite := bn256.NewSuiteG2()
	scheme := tbls.NewThresholdSchemeOnG1(bn256.NewSuiteG2())
	poly := share.NewPubPoly(suite, suite.Point().Base(), m.commitments)
	sig, err := scheme.Recover(poly, msg, partials, m.group.GetThreshold(), len(m.group.GetMembers()))
	if err != nil {
		panic(err)
	}
	err = crypto.Verify(poly.Commit(), msg, sig)
	if err != nil {
		panic(err)
	}
	return sig
}

func checkSignedWith(partials [][]byte, s []byte) bool {
	for _, p := range partials {
		if bytes.Compare(p, s) == 0 {
			return true
		}
	}
	return false
}
