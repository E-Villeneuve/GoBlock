package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
)

type Transaction struct {
	ID    []byte
	TXin  []TXInput
	TXout []TXOutput
}

const subsidy = 10

type TXInput struct {
	TXid      []byte //
	TXout     int
	ScriptSig string // Provides data to be used in an output's ScriptPubKey
}

type TXOutput struct {
	Value        int
	ScriptPubKey string // Yields destination address of the transaction output when hashed
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encoder := gob.NewEncoder(&encoded)
	encoder.Encode(tx)
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// Returns hash of all transactions in specified block
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions { // Concatenate hashes of every transaction
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{})) // Hash the concatenation

	return txHash[:]
}

func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{subsidy, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()

	return &tx
}

func (input *TXInput) CanUnlockOutput(key string) bool {
	return input.ScriptSig == key
}

func (output *TXOutput) CanBeUnlocked(key string) bool {
	return output.ScriptPubKey == key
}

func (bc *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTransactions []Transaction
	spentTransactions := make(map[string][]int)
	bc_iterator := bc.Iterator()

	for {
		block := bc_iterator.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outputIDs, out := range tx.TXout {
				// TODO: Complete this

			}

		}
	}
}
