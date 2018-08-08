package todos

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	uuid "github.com/satori/go.uuid"
)

func createUUID() (string, error) {
	b := uuid.NewV4()
	s := string(b[:])
	return s, nil
}

type Todo struct {
	ID          string `json:"id"`
	Completed   bool   `json:"completed"`
	Description string `json:"description"`
}

func GetTodo(id string) (Todo, error) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	todo := Todo{}

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String(id),
			},
		},
	})

	if err != nil {
		fmt.Println(err.Error())
		return todo, err
	}

	// Unmarshall the result in to an Item
	err = dynamodbattribute.UnmarshalMap(result.Item, &todo)
	if err != nil {
		fmt.Println(err.Error())
		return todo, err
	}

	return todo, nil
}

func DeleteTodo(id string) error {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String(id),
			},
		},
	}

	_, err := svc.DeleteItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func PostTodo(description string) (Todo, error) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	id, err := createUUID()

	// Marshall the request body
	todo := Todo{
		Description: "",
		ID:          id,
		Completed:   false,
	}
	json.Unmarshal([]byte(description), &todo)

	// Marshall the Item into a Map DynamoDB can deal with
	av, err := dynamodbattribute.MarshalMap(todo)
	if err != nil {
		fmt.Println("Got error marshalling map:")
		fmt.Println(err.Error())
		return todo, err
	}

	// Create Item in table and return
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	}
	_, err = svc.PutItem(input)
	return todo, err
}
