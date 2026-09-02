package main

import (
	"api"
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

// --- グローバル変数ブロック(コールドスタート間で使い回す) ---
var adaptor *httpadapter.HandlerAdapter

func init() {
	mux := http.NewServeMux()

	// DI
	dynamoDBClient := api.NewClient()
	tablebasic := api.NewTableBasic(dynamoDBClient)
	service := api.NewService(tablebasic)

	// エンドポイント定義
	mux.HandleFunc("GET /v1/users", service.GetUsers)
	mux.HandleFunc("POST /v1/users", service.PostUsers)

	// 本来ならListenAndServeするところをそのままhttpadaptorに渡している
	adaptor = httpadapter.New(mux)
}

// HTTPアクションをAPIGatewayの独自のものから、他のAPI系のHandlerで扱える一般的なAPIリクエストに変換
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return adaptor.ProxyWithContext(ctx, req)
}

// エントリーポイント
// lambda関数自体はここから実行される
func main() {
	lambda.Start(Handler)
}
