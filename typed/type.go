package typed

// RunConfig 启动引用的相关配置
type RunConfig struct {
	Server Server `json:"server"`
}

// Server 启动服务器的相关配置
type Server struct {
	Port int `json:"port"`
}
