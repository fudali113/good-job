package main

import (
	"net/http"
	"fmt"
	"flag"
)

var id, token, triggerResource string

func main()  {

	flag.StringVar(&id, "id", "", "触发资源的 id")
	flag.StringVar(&token, "token", "", "你的相关触发 job 或者 pipeline 的 token")
	flag.StringVar(&triggerResource, "resource", "job", "触发 job 或者是 pipeline")

	flag.Parse()

	uriTemplate := "good-job-server/v1/%s/%s/start?token=%s"

	if token == "" || id == "" {
		panic("token 或者 id 不能为空")
	}

	uri := fmt.Sprintf(uriTemplate, triggerResource, id, token)
	resp, err := http.Get(uri)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode >= 300 {
		panic(fmt.Sprintf("请求返回结果出错，StatusCode: %d", resp.StatusCode))
	}

	fmt.Println("触发请求成功")

}
