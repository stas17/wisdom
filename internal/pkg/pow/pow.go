package pow

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"time"
)

const targetBits = 24
const maxNonce = math.MaxInt

//ProofOfWork implementation of proof of work service. Uses the SHA256 algorithm for validation
type ProofOfWork interface {
	PrepareData(int) []byte
	Run() (int, []byte)
	Validate() bool
}

type proofOfWork struct {
	block  *Block
	target *big.Int
}

//NewPoofOfWork gets the service which implements prow of work
func NewPoofOfWork(b *Block) ProofOfWork {
	t := big.NewInt(1)
	t.Lsh(t, uint(256-targetBits))

	return &proofOfWork{
		block:  b,
		target: t,
	}
}

var _ ProofOfWork = (*proofOfWork)(nil)

//NewBlock gets the block. During creating a new block computed nonce and hash
func NewBlock(data string) *Block {
	block := &Block{int(time.Now().Unix()), data, string([]byte{}), 0}
	pow := NewPoofOfWork(block)

	nonce, hash := pow.Run()
	block.Hash = string(hash[:])
	block.Nonce = nonce

	return block
}

//PrepareData prepares data for hashing
func (pow *proofOfWork) PrepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			[]byte(pow.block.Data),
			[]byte(fmt.Sprintf("%x", pow.block.Timestamp)),
			[]byte(fmt.Sprintf("%x", int64(nonce))),
		},
		[]byte{},
	)

	return data
}

//Run does a main work. Increases a nonce each time and checks the result of cashing
func (pow *proofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	for nonce < maxNonce {
		data := pow.PrepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		}
		nonce++
	}

	return nonce, hash[:]
}

//Validate checks work. The number representation of the computed hash must be less than the target
func (pow *proofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.PrepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}
