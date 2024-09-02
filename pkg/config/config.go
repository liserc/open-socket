package config

type Log struct {
	StorageLocation     string `mapstructure:"storageLocation"`
	RotationTime        uint   `mapstructure:"rotationTime"`
	RemainRotationCount uint   `mapstructure:"remainRotationCount"`
	RemainLogLevel      int    `mapstructure:"remainLogLevel"`
	IsStdout            bool   `mapstructure:"isStdout"`
	IsJson              bool   `mapstructure:"isJson"`
	IsSimplify          bool   `mapstructure:"isSimplify"`
	WithStack           bool   `mapstructure:"withStack"`
}

type Discovery struct {
	Enable string `mapstructure:"enable"`
	Etcd   Etcd   `mapstructure:"etcd"`
}

type Etcd struct {
	RootDirectory string   `mapstructure:"rootDirectory"`
	Address       []string `mapstructure:"address"`
	Username      string   `mapstructure:"username"`
	Password      string   `mapstructure:"password"`
}

type Prometheus struct {
	Enable bool  `mapstructure:"enable"`
	Ports  []int `mapstructure:"ports"`
}

type Gateway struct {
	RPC struct {
		RegisterIP string `mapstructure:"registerIP"`
		Ports      []int  `mapstructure:"ports"`
	} `mapstructure:"rpc"`
	Prometheus Prometheus `mapstructure:"prometheus"`
	ListenIP   string     `mapstructure:"listenIP"`
	WS         struct {
		Ports               []int `mapstructure:"ports"`
		WebsocketMaxConnNum int   `mapstructure:"websocketMaxConnNum"`
		WebsocketMaxMsgLen  int   `mapstructure:"websocketMaxMsgLen"`
		WebsocketTimeout    int   `mapstructure:"websocketTimeout"`
	} `mapstructure:"ws"`
	MultiLoginPolicy int `mapstructure:"multiLoginPolicy"`
}

type Share struct {
	Secret          string          `mapstructure:"secret"`
	RpcRegisterName RpcRegisterName `mapstructure:"rpcRegisterName"`
	AdminUserID     []string        `mapstructure:"adminUserID"`
}

type RpcRegisterName struct {
	Push    string `mapstructure:"push"`
	Gateway string `mapstructure:"gateway"`
}

func (r *RpcRegisterName) GetServiceNames() []string {
	return []string{
		r.Push,
		r.Gateway,
	}
}
