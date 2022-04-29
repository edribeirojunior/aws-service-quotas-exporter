package servicequotas

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type TablesNames struct {
	Name string
}


func (m *mockDynamodbClient) ListTablesPages(input *dynamodb.ListTablesInput, fn func(*dynamodb.ListTablesOutput, bool) bool) error {
	fn(m.ListTableResponse, true)
	return m.err
} 

func TestDynamodbTableUsageWithError(t *testing.T) {
	mockClient := &mockDynamodbClient{
		err: errors.New("some err"),
		ListTableResponse: nil,
	}

	check := MaximumNumberTablesUsageCheck{mockClient}
	usage, err := check.Usage()

	assert.Error(t, err)
	assert.Error(t, errors.Is(err, ErrFailedToGetUsage))
	assert.Nil(t, usage)
}

func TestDynamodbTablesUsage(t * testing.T) {
	testCases := []struct {
		name string
		tableNames []*TableNames{},
		expectedUsage: []QuotaUsage{} 
	}{
		{
			name: "WithNoTables",
			tablesNames: []*TablesNames{},
			expectedUsage: []QuotaUsage{},	
		},
		{
			name: "WithTables",
			repositories: []*TablesNames{
				{
					Name: "test1"
				},
				{
					Name: "test2"
				},
				{
					Name: "test3"
				},
			},
			expectedUsage: []QuotaUsage{
				{
					Name: maximumNumberOfTablesName,
					ResourceName: "dynamodb",
					Description: maximumNumberOfTablesDesc,
					Usage: 0,
				},
				{
					Name: maximumNumberOfTablesName,
					ResourceName: "dynamodb",
					Description: maximumNumberOfTablesDesc,
					Usage: 3,
				}
			}
		}
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &mockDynamoDbClient{
				err: nil,
				ListTableResponse: &dynamodb.ListTablesOutput{
					TablesNames tc.tablesNames,
				},
			}

			check := MaximumNumberTablesUsageCheck{mockClient}
			usage, err := check.Usage()

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedUsage, usage)
		})
	}
}
