package repository

import (
	"github.com/BrunoDM2943/go-todo-lambda/internal/constants/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const TableName = "todo"

type dynamoDBRepo struct {
	client *dynamodb.DynamoDB
}

func NewDynamoDB() TodoRepository {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	client := dynamodb.New(sess)

	return &dynamoDBRepo{
		client: client,
	}
}

func (repo *dynamoDBRepo) Save(item *model.Item) error {
	marshalled, _ := dynamodbattribute.MarshalMap(item)
	_, err := repo.client.PutItem(&dynamodb.PutItemInput{
		Item:      marshalled,
		TableName: aws.String(TableName),
	})
	return err
}
func (repo *dynamoDBRepo) FindByID(id string) (*model.Item, error) {
	result, err := repo.client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, nil
	}
	item := &model.Item{}
	dynamodbattribute.UnmarshalMap(result.Item, item)
	return item, nil
}
func (repo *dynamoDBRepo) ListAll() ([]*model.Item, error) {

	result, err := repo.client.Scan(&dynamodb.ScanInput{
		TableName: aws.String(TableName),
	})
	if err != nil {
		return nil, err
	}

	items := make([]*model.Item, len(result.Items))

	for _, scannedItem := range result.Items {
		item := &model.Item{}
		_ = dynamodbattribute.UnmarshalMap(scannedItem, item)
		items = append(items, item)
	}

	return items, nil
}
func (repo *dynamoDBRepo) DeleteByID(id string) error {
	_, err := repo.client.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}
