# Good Job

an job scheduler in kubernetes

## 生成调用 k8s CRD api 代码的生成
调用 k8s CRD 的代码是由 [k8s-code-generate](https://github.com/kubernetes/code-generator) 项目生成，
我们在修改实体的类型时，不可以自己修改生成的代码，应该修改 `pkg/apis/goodjobcontroller/v1alpha1/types.go` 里面的实体类型然后
重新运行 `./hack/update-codegen.sh` 脚本进行重新生成


