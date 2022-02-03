package signature

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"wisdom/internal/server/config"
)

//Signature implementation of signature
type Signature interface {
	Generate() string
	Validate(string) (bool, error)
}

type signature struct {
	config config.Loader
}

// NewSignatureService gets the service for work with the signature
func NewSignatureService(c config.Loader) Signature {
	return &signature{
		config: c,
	}
}

var _ Signature = (*signature)(nil)

//Generate makes a new signature
func (s *signature) Generate() string {
	return base64.StdEncoding.EncodeToString(bytes.Join(
		[][]byte{
			[]byte(s.config.Get().App.Salt),
			[]byte(delimiter),
			[]byte(time.Now().String()),
			[]byte(delimiter),
			[]byte(fmt.Sprintf("%f", rand.Float32())),
		},
		[]byte{},
	))
}

//Validate checks that signature content right salt
func (s *signature) Validate(signature string) (bool, error) {
	decodeString, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, err
	}
	ps := strings.Split(string(decodeString), delimiter)
	if ps[0] == s.config.Get().App.Salt {
		return true, nil
	}
	return false, nil
}
