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

// table signers singleton
type Signers struct {
	public_keys []chain.PublicKey
}

func VerifySignatures(data []byte, signatures []chain.Signature) bool {
	digest := chain.Sha256(data)
	signerTable := NewSignersTable(MTG_XIN, MTG_XIN)
	signers := signerTable.Get()
	assert(signers != nil, "no signers")

	threshold := len(signers.public_keys)/3*2 + 1
	validSignatures := 0

	verfiedSignatures := make([]*chain.Signature, 0, len(signers.public_keys))

	for i := 0; i < len(signatures); i++ {
		signature := signatures[i]
		CheckDuplicatedSignature(verfiedSignatures, &signature)
		verfiedSignatures = append(verfiedSignatures, &signature)

		pub_key := chain.RecoverKey(digest, &signature)
		for _, public_key := range signers.public_keys {
			if public_key == *pub_key {
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
