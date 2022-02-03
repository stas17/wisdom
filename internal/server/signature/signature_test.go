package signature

import (
	"encoding/base64"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"wisdom/internal/mocks/mock_config"
	"wisdom/internal/server/config"
)

func TestSignature_Generate(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := getSignature(ctrl)
	s := service.Generate()

	decodeString, err := base64.StdEncoding.DecodeString(s)
	assert.NoError(t, err)
	
	ps := strings.Split(string(decodeString), delimiter)
	assert.Equal(t, 3, len(ps))
	assert.Equal(t, "42", ps[0])
}

func TestSignature_Validate(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := getSignature(ctrl)
	notOk, err := service.Validate("MzJ8MjAyMi0wMi0wMyAxNzo1ODozOS42ODk0MDUgKzAzMDAgTVNLIG09KzAuMDAxOTMwNDA0fDAuNjA0NjYw")
	assert.NoError(t, err)
	assert.False(t, notOk)

	ok, err := service.Validate("NDJ8MjAyMi0wMi0wMyAxNzo1Nzo0OS4zMzQ3NjUgKzAzMDAgTVNLIG09KzAuMDAxMjA4ODQ1fDAuNjA0NjYw")
	assert.NoError(t, err)
	assert.True(t, ok)
}

func getSignature(ctrl *gomock.Controller) Signature {
	configLoader := mock_config.NewMockLoader(ctrl)

	cfg := config.Config{App: &config.App{Salt: "42"}}
	configLoader.EXPECT().Get().AnyTimes().Return(&cfg)

	return NewSignatureService(configLoader)
}
