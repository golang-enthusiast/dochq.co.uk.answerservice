package dynamodb

import (
	"fmt"

	"dochq.co.uk.answerservice/internal/domain"
	errors "dochq.co.uk.answerservice/internal/error"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	awsDynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type answerRepo struct {
	db        *awsDynamodb.DynamoDB
	tableName string
}

// NewAnswerRepository creates a new repository.
func NewAnswerRepository(session *awsSession.Session, tableName string) domain.AnswerRepository {

	// Create a new dynamodb client.
	//
	db := awsDynamodb.New(session)

	// Ensure is table exist.
	//
	answerTableMustExist(db, tableName)

	return &answerRepo{
		db:        db,
		tableName: tableName,
	}
}

func answerTableMustExist(db *awsDynamodb.DynamoDB, tableName string) {
	// Check if the table exist.
	//
	listTablesOutput, err := db.ListTables(&awsDynamodb.ListTablesInput{})

	// Returned internal server error.
	// The application does not know if the table exists or not.
	// Thus, it cannot query the server, so we panic this error.
	//
	if err != nil {
		panic(err)
	}
	for _, table := range listTablesOutput.TableNames {
		// the table already exists and there is no reason to continue.
		if *table == tableName {
			return
		}
	}

	// Create table.
	//
	_, err = db.CreateTable(&awsDynamodb.CreateTableInput{
		AttributeDefinitions: []*awsDynamodb.AttributeDefinition{
			{
				AttributeName: aws.String(domain.JSONFieldAnswerKey),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*awsDynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(domain.JSONFieldAnswerKey),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &awsDynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	})
	if aerr, ok := err.(awserr.Error); ok {
		if aerr.Code() != awsErrorResourceInUse {
			panic(aerr)
		}
	}

	// Wait for table.
	_ = db.WaitUntilTableExists(&awsDynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
}

func (r *answerRepo) Create(answer *domain.Answer) error {

	// Marshal Go value type to a map of AttributeValues.
	//
	attributes, err := dynamodbattribute.MarshalMap(answer)
	if err != nil {
		return err
	}

	// Put input.
	//
	input := &awsDynamodb.PutItemInput{
		Item:      attributes,
		TableName: aws.String(r.tableName),
	}

	// Put item in dynamodb storage.
	//
	_, err = r.db.PutItem(input)

	return err
}

func (r *answerRepo) Update(answer *domain.Answer) error {

	// Marshal Go value type to a map of AttributeValues.
	//
	attributes, err := dynamodbattribute.MarshalMap(answer)
	if err != nil {
		return err
	}

	// Put input.
	//
	input := &awsDynamodb.PutItemInput{
		Item:      attributes,
		TableName: aws.String(r.tableName),
	}

	// Put item in dynamodb storage.
	//
	_, err = r.db.PutItem(input)

	return err
}

func (r *answerRepo) Delete(key domain.AnswerKey) error {

	// Delete input.
	//
	input := &awsDynamodb.DeleteItemInput{
		Key: map[string]*awsDynamodb.AttributeValue{
			domain.JSONFieldAnswerKey: {
				S: aws.String(string(key)),
			},
		},
		TableName: aws.String(r.tableName),
	}

	// Delete item.
	//
	_, err := r.db.DeleteItem(input)

	return err
}

func (r *answerRepo) Get(key domain.AnswerKey) (*domain.Answer, error) {

	// Build the query input parameters.
	//
	scanInput := &awsDynamodb.ScanInput{
		TableName: aws.String(r.tableName),
		ScanFilter: map[string]*awsDynamodb.Condition{
			domain.JSONFieldAnswerKey: {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*awsDynamodb.AttributeValue{
					{S: aws.String(string(key))},
				},
			},
		},
		ConsistentRead: aws.Bool(true),
	}

	// Make the DynamoDB Query API call.
	//
	result, err := r.db.Scan(scanInput)
	if err != nil {
		return nil, errors.NewErrInternal(fmt.Sprintf("Query API call failed: %s", err))
	}

	// Return error if nothing found.
	//
	if len(result.Items) == 0 {
		return nil, errors.NewErrNotFound("Answer not found")
	}

	// Unmarshal entity.
	//
	answer := &domain.Answer{}
	err = dynamodbattribute.UnmarshalMap(result.Items[0], &answer)
	if err != nil {
		return nil, errors.NewErrInternal(fmt.Sprintf("Got error unmarshalling: %s", err))
	}

	// Return result.
	//
	return answer, nil
}
