# Good Job

an ultralight job scheduler in kubernetes

## 生成调用 k8s CRD api 代码的生成
调用 k8s CRD 的代码是由 [k8s-code-generate](https://github.com/kubernetes/code-generator) 项目生成，
我们在修改实体的类型时，不可以自己修改生成的代码，应该修改 `pkg/apis/goodjobcontroller/v1alpha1/types.go` 里面的实体类型然后
重新运行 `./hack/update-codegen.sh` 脚本进行重新生成


## 思路
### 配置
依托与 kubernetes 的 Job 进行分布式的任务调度，定义了 GoodJob 类型，把一个 Job 分为 Sharding 和 Runing 两个阶段，Sharding 阶段为为任务数据进行分片，有业务代码自主实现，
服务只保存分片结果(所以我们的分片应该是根据我们实际运行的可以划分的区间而不是实际的值，比如`1，2，3，4，5，6，7，8，9`进行分片，分`3`片，我们应该分为`[1-3，4-6，7-9`]而不是`[[1,2,3],[4,5,6],[7,8,9]]`),
我们应该让平台相关的参数区间进行分片，而不是保存实际的值；分片之后平台将会按照分片对任务进行并行执行(每个分片对应一个`k8s job`,我们会将分片参数以参数的形式传入运行的容器)
配置大致如下:
```
apiVersion: goodjob.k8s.io/v1alpha1
kind: GoodJob
metadata:
    name: test
    namespace: good-job
spec:
    template:                   #  参照 PodTemplate
        name: "test"
        image: "test"
        cmd: ["test"]
        args: []
    shard:
        matchPattern: ""    # 正则表达式， 默认 `GOOD_JOB_SHARDS[\s\S\w\W]+`
        type: config        # config or exec
        shards: []          # if type == config
        template:           # if type == exec ; 我们应该在程序中调用相关的 api 来进行更新改 job 的分片信息  (待优化，应该劲量与应用解耦)
            name: test
            image: "test"
            cmd: ["test"]
            args: []
    parallel: 2             # 并行度, 即最多可以同时运行多少个分片，大于分片数实际取值会是分片数
status:

```

PodTemplate 参见 [PodTemplate](https://v1-9.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.9/#podtemplatespec-v1-core)

#### shard type exec

当你使用 type : exec 的 shard 配置的时候，程序将会从你的 shard 容器执行日志当中去需要符合规则的条目并转换为 shard ; 你可以自己配置 `matchPattern`
正则表达式来从 shard 容器日志中匹配相关的日志条目，日志条目以 `\n` 进行切分；(因为我们会读取全部日志并进行匹配，所以建议 shard 容器日志行数不要太大)

### 实现原理
利用 CRD 创建 GoodJob 的资源，监听 GoodJob 资源的添加，将 GoodJob 拆分成具体的 Job 并执行，监听 Job 执行的情况并更新 GoodJob 的状态