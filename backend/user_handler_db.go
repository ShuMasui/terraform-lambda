package api

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// ここの責務は、あくまでもuser_handlerにまつわるDB操作処理をインターフェースとして提供するだけ

// --- 型定義ブロック ---
type TableBasic struct {
	Client    *DynamoDBClient
	TableName string
}

type User struct {
	UserID   string `dynamodbav:"user_id"`
	UserName string `dynamodbav:"user_name"`
}

type UserKey struct {
	UserID string `dynamodbav:"user_id"`
}

func NewTableBasic(c *DynamoDBClient) *TableBasic {
	if c == nil {
		return nil
	}

	return &TableBasic{Client: c, TableName: os.Getenv("DYNAMO_TABLE_USERS")}
}

func (tb *TableBasic) scanUsers(ctx context.Context) ([]User, error) {

	// --- 変数宣言ブロック ---
	var users []User
	var err error
	var response *dynamodb.ScanOutput

	response, err = tb.Client.Client.Scan(ctx, &dynamodb.ScanInput{TableName: aws.String(tb.TableName)})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = attributevalue.UnmarshalListOfMaps(response.Items, &users)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return users, nil
}

func (tb *TableBasic) putUser(ctx context.Context, user *User) error {

	// --- 変数宣言ブロック ---
	var err error

	av, err := attributevalue.MarshalMap(user)

	if err != nil {
		return err
	}

	if _, err = tb.Client.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tb.TableName),
		Item:      av,
	}); err != nil {
		return err
	}

	log.Println("正常に保存できました")

	return nil
}
