package dynamodb

import (
	"dochq.co.uk.answerservice/internal/domain"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	awsDynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type answerEventRepo struct {
	db        *awsDynamodb.DynamoDB
	tableName string
}

// NewAnswerEventRepository creates a new repository.
func NewAnswerEventRepository(session *awsSession.Session, tableName string) domain.AnswerEventRepository {

	// Create a new dynamodb client.
	//
	db := awsDynamodb.New(session)

	// Ensure is table exist.
	//
	answerEventTableMustExist(db, tableName)

	return &answerEventRepo{
		db:        db,
		tableName: tableName,
	}
}

func answerEventTableMustExist(db *awsDynamodb.DynamoDB, tableName string) {
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
			{
				AttributeName: aws.String(domain.JSONFieldEventType),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*awsDynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(domain.JSONFieldAnswerKey),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String(domain.JSONFieldEventType),
				KeyType:       aws.String("RANGE"),
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

func (r *answerEventRepo) Create(answerEvent *domain.AnswerEvent) error {

	// Marshal Go value type to a map of AttributeValues.
	//
	attributes, err := dynamodbattribute.MarshalMap(answerEvent)
	if err != nil {
		return err
	}

	// Add missing attributes.
	//
	attributes[domain.JSONFieldAnswerKey] = &awsDynamodb.AttributeValue{
		S: aws.String(string(answerEvent.Data.Key)),
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

func (r *answerEventRepo) ListEvents(key domain.AnswerKey) ([]*domain.AnswerEvent, error) {

	// Build expression.
	//
	filter := expression.Name(domain.JSONFieldAnswerKey).Equal(expression.Value(key))

	expr, err := expression.NewBuilder().WithProjection(r.getProjection()).WithFilter(filter).Build()
	if err != nil {
		return nil, err
	}

	params := &awsDynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(r.tableName),
	}

	// Make the DynamoDB Query API call
	result, err := r.db.Scan(params)
	if err != nil {
		return nil, err
	}
	var items []*domain.AnswerEvent
	for _, i := range result.Items {
		item := &domain.AnswerEvent{}
		err = dynamodbattribute.UnmarshalMap(i, item)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *answerEventRepo) getProjection() expression.ProjectionBuilder {
	return expression.NamesList(
		expression.Name(domain.JSONFieldAnswerKey),
		expression.Name(domain.JSONFieldEventType),
		expression.Name(domain.JSONFieldData))
}
