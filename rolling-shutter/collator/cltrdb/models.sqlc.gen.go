// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0

package cltrdb

import ()

type DecryptionKey struct {
	EpochID       []byte
	DecryptionKey []byte
}

type DecryptionTrigger struct {
	EpochID   []byte
	BatchHash []byte
}

type EonPublicKeyCandidate struct {
	Hash                  []byte
	EonPublicKey          []byte
	ActivationBlockNumber int64
	KeyperConfigIndex     int64
	Eon                   int64
	Confirmed             bool
}

type EonPublicKeyVote struct {
	Hash              []byte
	Sender            string
	Signature         []byte
	Eon               int64
	KeyperConfigIndex int64
}

type MetaInf struct {
	Key   string
	Value string
}

type NextEpoch struct {
	EnforceOneRow bool
	EpochID       []byte
}

type Transaction struct {
	TxID        []byte
	EpochID     []byte
	EncryptedTx []byte
}
