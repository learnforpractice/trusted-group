package main

import (
	"github.com/uuosio/chain"
)

func check(b bool, msg string) {
	chain.Check(b, msg)
}

func assert(b bool, msg string) {
	chain.Assert(b, msg)
}

// table signers
type Signer struct {
	account    chain.Name //primary : t.account.N
	public_key chain.PublicKey
}

func VerifySignatures(data []byte, signatures []chain.Signature) bool {
	digest := chain.Sha256(data)
	signerTable := NewSignerTable(MTG_XIN, MTG_XIN)
	signers := make([]*Signer, 0, 10)
	it := signerTable.Lowerbound(0)
	for it.IsOk() {
		item := signerTable.GetByIterator(it)
		signers = append(signers, item)
		it, _ = signerTable.Next(it)
	}

	threshold := len(signers)/3*2 + 1
	validSignatures := 0

	verfiedSignatures := make([]*chain.Signature, 0, len(signers))

	for i := 0; i < len(signatures); i++ {
		signature := signatures[i]
		CheckDuplicatedSignature(verfiedSignatures, &signature)
		verfiedSignatures = append(verfiedSignatures, &signature)

		pub_key := chain.RecoverKey(digest, &signature)
		for _, signer := range signers {
			if signer.public_key == *pub_key {
				validSignatures += 1
				break
			}
		}
		if validSignatures >= threshold {
			return true
		}
	}
	assert(false, "Not enough valid signatures")
	return false
}

func CheckDuplicatedSignature(signatures []*chain.Signature, signature *chain.Signature) {
	for _, sig := range signatures {
		if *sig == *signature {
			assert(false, "duplicated signature")
		}
	}
}

// table processes ignore
type Process struct {
	contract chain.Name //primary : t.contract.N
	process  chain.Uint128
}

func GetProcessId(contract chain.Name) chain.Uint128 {
	db := NewProcessTable(MTG_XIN, MTG_XIN)
	it, record := db.GetByKey(contract.N)
	assert(it.IsOk(), "process not found!")
	return record.process
}

func VerifyProcess(contract chain.Name, process chain.Uint128) {
	db := NewProcessTable(MTG_XIN, MTG_XIN)
	it, record := db.GetByKey(contract.N)
	assert(it.IsOk(), "process not found!")
	assert(record.process == process, "invalid process!")
}
