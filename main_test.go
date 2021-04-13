package main

import (
	"os"
	"testing"

	"github.com/RedHatInsights/rhc-worker-catalog/internal/common"
	"github.com/RedHatInsights/rhc-worker-catalog/internal/towerapiworker"
	"github.com/stretchr/testify/assert"
)

type FakeRequestHandler struct {
	catalogConfig common.CatalogConfig
	workHandler   towerapiworker.WorkHandler
}

func (frh *FakeRequestHandler) StartHandlingRequests(config *common.CatalogConfig, wh towerapiworker.WorkHandler) {
	frh.catalogConfig = *config
	frh.workHandler = wh
}

func TestMain(t *testing.T) {
	os.Args = []string{"catalog_worker", "--config", "./testdata/catalog_sample.toml"}

	frh := &FakeRequestHandler{}

	initConfig()
	logf := configLogger()
	startRun(makeConfig(), frh)

	info, err := logf.Stat()
	assert.NoError(t, err)
	logf.Close()
	os.Remove(info.Name())

	assert.True(t, info.Size() > 0)
	assert.Equal(t, "info", frh.catalogConfig.Level)
	assert.Equal(t, "<<Your Tower URL>>", frh.catalogConfig.URL)
	assert.Equal(t, "<<Your Tower Token>>", frh.catalogConfig.Token)
	assert.Equal(t, os.Getenv("HTTP_PROXY"), "http://myproxy:3128")
	assert.Equal(t, os.Getenv("HTTPS_PROXY"), "http://myproxy:3128")
	assert.Equal(t, os.Getenv("NO_PROXY"), "localhost")
	assert.Equal(t, &towerapiworker.DefaultAPIWorker{}, frh.workHandler)
}

func TestConfigSearchPath(t *testing.T) {
	assert.Panics(t, func() { initConfig() })

	configFile, err := getConfigFile()
	assert.NotNil(t, err)

	os.MkdirAll("./rhc/workers", os.ModePerm)
	os.Create("./rhc/workers/catalog.toml")
	configFile, err = getConfigFile()
	assert.Equal(t, "./rhc/workers/catalog.toml", configFile)
	os.RemoveAll("./rhc/workers")
}
