package main

import (
	"github.com/uuosio/chain"
)

func VerifySignatures(codeAccount chain.Name, data []byte, signatures []chain.Signature) bool {
	signerDB := NewSignersTable(codeAccount)
	signers := signerDB.Get()
	check(signers != nil, "no signers")

	threshold := len(signers.public_keys)*2/3 + 1
	validSignatures := 0

	verfiedSignatures := make([]*chain.Signature, 0, len(signers.public_keys))

	digest := chain.Sha256(data)
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
	check(false, "Not enough valid signatures: "+string([]byte{'0' + byte(validSignatures)}))
	return false
}

func CheckDuplicatedSignature(signatures []*chain.Signature, signature *chain.Signature) {
	for _, sig := range signatures {
		if *sig == *signature {
			check(false, "duplicated signature")
		}
	}
}

func hasDuplicates(publicKeys []chain.PublicKey) bool {
	for i := range publicKeys {
		for j := i + 1; j < len(publicKeys); j++ {
			if publicKeys[i] == publicKeys[j] {
				return true
			}
		}
	}
	return false
}

func check(b bool, msg string) {
	chain.Check(b, msg)
}
