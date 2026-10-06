package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
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

// Creates new coinbase transaction
// Coinbase transactions don't require prev existing outputs to create outputs
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

// Checks if transaction is a coinbase transaction
func (tx Transaction) IsCoinbase() bool {
	return len(tx.TXout) == 1 && len(tx.TXin[0].TXid) == 0 && tx.TXin[0].TXout == -1
}

func (input *TXInput) CanUnlockOutput(key string) bool {
	return input.ScriptSig == key
}

func (output *TXOutput) CanBeUnlocked(key string) bool {
	return output.ScriptPubKey == key
}

func (bc *Blockchain) FindUnspentTransactionsHelper(address string) []Transaction {
	var unspentTransactions []Transaction
	spentTransactions := make(map[string][]int)
	bc_iterator := bc.Iterator()

	for {
		block := bc_iterator.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outputID, out := range tx.TXout {
				if spentTransactions[txID] != nil {
					for _, spentOut := range spentTransactions[txID] {
						if spentOut == outputID { // Check if output was referenced in some other input, if so, skip
							continue Outputs
						}
					}
				}

				if out.CanBeUnlocked(address) {
					unspentTransactions = append(unspentTransactions, *tx)
				}

			}

			if tx.IsCoinbase() == false { // coinbase transactions can't unlock outputs, skip them
				for _, in := range tx.TXin {
					if in.CanUnlockOutput(address) {
						inTransactionID := hex.EncodeToString(in.TXid)
						spentTransactions[inTransactionID] = append(spentTransactions[inTransactionID], in.TXout)
					}
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}

	}
	return unspentTransactions

}

func (bc *Blockchain) FindUnspentTransactions(address string) []TXOutput {
	var UTXOs []TXOutput
	unspentTransactions := bc.FindUnspentTransactionsHelper(address)

	for _, transaction := range unspentTransactions {
		for _, output := range transaction.TXout {
			if output.CanBeUnlocked(address) {
				UTXOs = append(UTXOs, output)
			}
		}
	}

	return UTXOs
}

func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	total, validOutputs := bc.FindSpendableOutputs(from, amount)

	if total < amount {
		log.Panic("ERROR: Not enough funds")
	}

	// Create input list
	for txID, outs := range validOutputs {
		txID, err := hex.DecodeString(txID)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TXInput{txID, out, from}
			inputs = append(inputs, input)
		}

	}

	// Create output list
	outputs = append(outputs, TXOutput{amount, to})
	if total > amount {
		outputs = append(outputs, TXOutput{total - amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx

}

// Iterates over all unspent transactions to accumulate total value
// When accumulated value is equal or greater to desired total for transfer, return amount and relevant transaction IDs
func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTransactions := bc.FindUnspentTransactionsHelper(address)
	total := 0

Work:
	for _, transaction := range unspentTransactions {
		txID := hex.EncodeToString(transaction.ID)

		for outputID, output := range transaction.TXout {
			if output.CanBeUnlocked(address) && total < amount {
				total += output.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outputID)

				if total >= amount {
					break Work
				}
			}
		}

	}

	return total, unspentOutputs
}
