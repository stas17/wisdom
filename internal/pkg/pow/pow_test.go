package pow

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestPow_PrepareData(t *testing.T) {
	pow := getPow()
	data := pow.PrepareData(42)

	assert.Equal(t, "test2ac2a", string(data))
}

func TestPow_Run(t *testing.T) {
	var hashInt big.Int
	exp, _ := hashInt.SetString("13525537110524281909028781806727420424168970506817147624579040245994845563", 10)
	pow := getPow()

	nonce, hash := pow.Run()

	assert.Equal(t, 4639, nonce)
	assert.Equal(t, exp, hashInt.SetBytes(hash))
	fmt.Printf("Hash: %x\n", hash)
	fmt.Printf("Hash: %x\n", hash)

}

func TestPow_Validate(t *testing.T) {
	pow := getPow()
	assert.False(t, pow.Validate())

	b := &Block{
		42,
		"test",
		string([]byte{}),
		4639,
	}
	pow = NewPoofOfWork(b)
	assert.True(t, pow.Validate())

}

func getPow() ProofOfWork {
	b := &Block{
		42,
		"test",
		string([]byte{}),
		0,
	}
	return NewPoofOfWork(b)
}
