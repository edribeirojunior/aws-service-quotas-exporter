package servicequotas

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/pkg/errors"
)

const (
	maximumNumberOfTablesName = "maximum_number_of_tables"
	maximumNumberOfTablesDesc = "Maximum number of tables"
)

// MaximumNumberTablesCheck implements the UsageCheck interface
// for maximum number of tables
type MaximumNumberTablesUsageCheck struct {
	client dynamodbiface.DynamoDBAPI
}

// Usage returns the usage for Dynamodb table
func (c *MaximumNumberTablesUsageCheck) Usage() ([]QuotaUsage, error) {
	quotaUsages := []QuotaUsage{}

	params := &dynamodb.ListTablesInput{}
	err := c.client.ListTablesPages(params,
		func(page *dynamodb.ListTablesOutput, lastPage bool) bool {
			if page != nil {

				tags := map[string]string{}

				rName := "dynamodb_table_name"

				tablesUsage := QuotaUsage{
					Name:         maximumNumberOfTablesName,
					ResourceName: &rName,
					Description:  maximumNumberOfTablesDesc,
					Usage:        float64(len(page.TableNames)),
					Tags:         tags,
				}

				quotaUsages = append(quotaUsages, []QuotaUsage{tablesUsage}...)
			}
			return !lastPage
		},
	)
	if err != nil {
		return nil, errors.Wrapf(ErrFailedToGetUsage, "%w", err)
	}

	return quotaUsages, nil
}
