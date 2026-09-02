package api

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// このコードの責務は AWS/dynamodbへのインスタンスを持つこと
// あくまでもインスタンスを持つことのみである

type DynamoDBClient struct {
	Client *dynamodb.Client
}

func NewClient() *DynamoDBClient {
	// --- 設定読み込みブロック ---
	region := os.Getenv("AWS_REGION")

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))

	if err != nil {
		log.Fatal(err)
	}

	// --- クライアント組み立てブロック ---
	client := &DynamoDBClient{}

	client.Client = dynamodb.NewFromConfig(cfg)

	log.Println("DynamoDBクライアント初期化完了")

	return client
}
