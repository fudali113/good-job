package typed

const (
	JobCrdName      = "GoodJob"
	PipelineCrdName = "Pipeline"
	CornTrigger     = "CronTrigger"
)

const (
	None = iota
	Begin
	Sharding
	ShardFail
	ShardSuccess
	Runing
	RunFail
	RunSuccess
)
