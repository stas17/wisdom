//go:build unit
// +build unit

package config

import (
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"wisdom/internal/server/logger"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	testConfigPath, err := makeConfigFile()
	assert.NoError(t, err)

	t.Run("success testing conf", func(t *testing.T) {
		os.Setenv(ENVConfigFile, testConfigPath)
		defer os.Clearenv()
		loader := NewConfigLoader(os.Getenv(ENVConfigFile))
		err := loader.Load()
		assert.NoError(t, err, "cant load cfg:%+v", err)
		cfg := loader.Get()
		assert.Truef(t, cfg.Logging.Enabled, "incorrect config value: %+v", cfg)
	})

	realConfig, err := getRealConfig()
	require.NoErrorf(t, err, "can't create get default config")

	t.Run("success default conf", func(t *testing.T) {
		loader := NewConfigLoader(realConfig)
		err := loader.Load()
		cfg := loader.Get()
		assert.NoError(t, err, "cant load cfg:%+v", err)
		require.NotNil(t, cfg.Logging)
		assert.Truef(t, cfg.Logging.Enabled, "incorrect config value: %+v", cfg)
		assert.Equal(t, cfg.Logging.Level, zerolog.Level(logger.DebugLevel).String())
		assert.Equal(t, "hyiujnbh42bhyujnbhyuj", cfg.App.Salt)
		assert.Equal(t, ":8080", cfg.App.HTTPWisdomAddr)
	})
}

func makeConfigFile() (string, error) {
	cfg := Config{
		App: &App{
			HTTPWisdomAddr: ":8080",
			Salt:           "Salt",
		},
		Logging: &Logging{
			Enabled: true,
		},
	}

	b, err := yaml.Marshal(cfg)
	if err != nil {
		return "", errors.Wrap(err, "marshal error")
	}

	p := "/tmp/_cfg.yaml"

	if err := ioutil.WriteFile(p, b, 0777); err != nil {
		return "", errors.Wrap(err, "save error")
	}

	return p, nil
}

func getRealConfig() (string, error) {
	curDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return path.Join(curDir, "/..", "/..", "/..", defaultConfigPath), nil
}
