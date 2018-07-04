package typed

// RunConfig 启动引用的相关配置
type RunConfig struct {
	Server ServerConfig `json:"server"`
	Runtime RuntimeConfig `json:"controller"`
}

// ServerConfig 启动服务器的相关配置
type ServerConfig struct {
	Port int `json:"port"`
}

// RuntimeConfig 运行时的相关配置
type RuntimeConfig struct {
}

