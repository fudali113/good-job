package storage

type Pipeline struct {
	Id   string           `json:"id"`
	Name string           `json:"name"`
	Jobs map[string][]Job `json:"jobs"`
}

// Job 执行 Job 的配置
type Job struct {
	Id string `json:"id"`
	// Job 的名字
	Name string `json:"name"`
	// 储存分片后的数据
	Shards []interface{} `json:"shards"`
	// 储存执行的程序
	Exec ExecConfig `json:"exec"`
	// 分片的配置
	Shard ShardConfig `json:"shard"`
	// 指定并行度
	Parallel int `json:"parallel"`
}

// ExecConfig 执行 job 程序的配置
type ExecConfig struct {
	// 镜像
	Image string `json:"image"`
	// 启动命令
	Cmd []string `json:"cmd"`
	// 启动参数
	Args []string `json:"args"`
	// 环境变量
	Env []string `json:"env"`
}

// ShardConfig 分片程序的配置
type ShardConfig struct {
	// 分片的类型
	Type string `json:"type"`
	// 执行分片程序的配置
	Exec ExecConfig `json:"exec"`
	// 手动设置分片
	Shards []interface{} `json:"shards"`
}

type JobStorage interface {
	Get(id string) (Job, error)
	Create(job Job) error
	Update(job Job) error
	Delete(id string) error
}
