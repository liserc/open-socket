package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadLogConfig(t *testing.T) {
	var log Log
	err := LoadConfig("../../../config/log.yml", "IMENV_LOG", &log)
	assert.Nil(t, err)
	assert.Equal(t, "../../../../logs/", log.StorageLocation)
}
