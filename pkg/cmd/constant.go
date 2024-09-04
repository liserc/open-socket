package cmd

import "strings"

var (
	FileName             string
	LogCfgFileName       string
	ShareFileName        string
	KafkaCfgFileName     string
	DiscoveryCfgFilename string
	GatewayCfgFileName   string
	ApiCfgFileName       string
	PushCfgFileName      string
)

var ConfigEnvPrefixMap map[string]string

func init() {
	FileName = "config.yaml"
	LogCfgFileName = "log.yml"
	ShareFileName = "share.yml"
	KafkaCfgFileName = "kafka.yml"
	DiscoveryCfgFilename = "discovery.yml"
	GatewayCfgFileName = "gateway.yml"
	ApiCfgFileName = "api.yml"
	PushCfgFileName = "push.yml"
	ConfigEnvPrefixMap = make(map[string]string)
	fileNames := []string{
		FileName, LogCfgFileName, ShareFileName,
		KafkaCfgFileName, DiscoveryCfgFilename,
		GatewayCfgFileName, ApiCfgFileName, PushCfgFileName,
	}

	for _, fileName := range fileNames {
		envKey := strings.TrimSuffix(strings.TrimSuffix(fileName, ".yml"), ".yaml")
		envKey = "IMENV_" + envKey
		envKey = strings.ToUpper(strings.ReplaceAll(envKey, "-", "_"))
		ConfigEnvPrefixMap[fileName] = envKey
	}
}

const (
	FlagConf          = "config_folder_path"
	FlagTransferIndex = "index"
)
