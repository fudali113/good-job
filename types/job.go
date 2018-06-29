package types

// Job 执行 Job 的配置
type Job struct {
	Id string
	// Job 的名字
	Name string
	// 储存分片后的数据
	Shards []interface{}
	// 储存执行的程序
	Exec ExecConfig
	// 分片的配置
	Shard ShardConfig
}

// ExecConfig 执行 job 程序的配置
type ExecConfig struct {
	// 镜像
	Image string
	// 启动命令
	Cmd []string
	// 启动参数
	Args []string
	// 环境变量
	Env []string
}

// ShardConfig 分片程序的配置
type ShardConfig struct {
	// 分片的类型
	Type   string
	// 执行分片程序的配置
	Exec   ExecConfig
	// 手动设置分片
	Shards []interface{}
}

