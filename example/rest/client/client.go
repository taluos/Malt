package httpclient

import (
	"context"
	"time"

	restClient "github.com/taluos/Malt/client/rest"
	// 改为导入fasthttp客户端包
	restfasthttp "github.com/taluos/Malt/client/rest/rest-fasthttp"
)

func ClientInit() restClient.Client {
	client, err := restClient.NewClient("fasthttp",
		"http://127.0.0.1:8080",
		// 使用fasthttp的WithTimeout选项
		restfasthttp.WithTimeout(5*time.Second),
	)
	if err != nil {
		panic(err)
	}
	return client
}

func ClientGet(cli restClient.Client, ctx context.Context, path string) (restClient.Response, error) {
	response, err := cli.Get(ctx, path)
	return response, err
}

func ClientStop(cli restClient.Client, ctx context.Context) error {
	return cli.Close(context.Background())
}
