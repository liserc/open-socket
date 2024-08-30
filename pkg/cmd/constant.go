package cmd

const (
	FlagConf          = "config_folder_path"
	FlagTransferIndex = "index"
)

var (
	FileName          string
	LogConfigFileName string
)

var ConfigEnvPrefixMap map[string]string

func init() {
	FileName = "config.yaml"
	LogConfigFileName = "log.yml"
}
