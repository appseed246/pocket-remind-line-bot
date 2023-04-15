package datasource

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Db struct {
	client *dynamodb.DynamoDB
}

type PocketReminderUser struct {
	UserId            string `dynamodbav:"userId"`
	LineUserId        string `dynamodbav:"lineUserId"`
	PocketAccessToken string `dynamodbav:"pocketAccessToken"`
}

func New() (*Db, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1")},
	)

	if err != nil {
		fmt.Println("Error creating session:")
		fmt.Println(err.Error())
		return nil, err
	}

	db := &Db{
		client: dynamodb.New(sess),
	}

	return db, nil
}

func (db *Db) GetPocketReminderUser(lineUserId string) (*PocketReminderUser, error) {
	// lineのuserIdからアプリのユーザIDを取得する
	q := &dynamodb.QueryInput{
		TableName: aws.String("PocketReminderTable"),
		IndexName: aws.String("LineUserIdIndex"),
		KeyConditions: map[string]*dynamodb.Condition{
			"lineUserId": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(lineUserId),
					},
				},
			},
		},
	}

	result, err := db.client.Query(q)
	if err != nil {
		fmt.Println("Failed to query item: " + err.Error())
		return nil, err
	}

	if len(result.Items) > 0 {
		head := result.Items[0]
		return &PocketReminderUser{
			UserId:            *head["userId"].S,
			LineUserId:        *head["lineUserId"].S,
			PocketAccessToken: *head["pocketAccessToken"].S,
		}, nil
	} else {
		return nil, nil
	}
}

func (db *Db) SaveAccessToken(userId string, lineUserId string, AccessToken string) error {

	pocketUser := &PocketReminderUser{
		UserId:            userId,
		LineUserId:        lineUserId,
		PocketAccessToken: AccessToken,
	}

	av, err := dynamodbattribute.MarshalMap(pocketUser)
	if err != nil {
		fmt.Println("Error marshalling struct:")
		fmt.Println(err.Error())
		return err
	}

	putInput := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("PocketReminderTable"),
	}

	_, err = db.client.PutItem(putInput)
	if err != nil {
		fmt.Println("Error saving item to DynamoDB:")
		fmt.Println(err.Error())
		return err
	}
	return nil
}
