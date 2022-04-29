package servicequotas

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type Repository struct {
	Name string
}


func (m *mockECRClient) DescribeRepositoriesPages(input *ecr.DescribeRepositoriesInput, fn func(*ecr.DescribeRepositoriesOutput, bool) bool) error {
	fn(m.DescribeECRResponse, true)
	return m.err
} 

func TestECRReposUsageWithError(t *testing.T) {
	mockClient := &mockECRClient{
		err: errors.New("some err"),
		DescribeECRResponse: nil,
	}

	check := ImagesPerRepositoryUsageCheck{mockClient}
	usage, err := check.Usage()

	assert.Error(t, err)
	assert.Error(t, errors.Is(err, ErrFailedToGetUsage))
	assert.Nil(t, usage)
}

func TestECRReposUsage(t * testing.T) {
	testCases := []struct {
		name string
		repositories []*Repository{},
		expectedUsage: []QuotaUsage{} 
	}{
		{
			name: "WithNoRepositories",
			repositories: []*Repository{},
			expectedUsage: []QuotaUsage{},	
		},
		{
			name: "WithRepositories",
			repositories: []*Repository{
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
					Name: imagesPerRepositoryName,
					ResourceName: "ecr_repos",
					Description: imagesPerRepositoryDesc,
					Usage: 0,
				},
				{
					Name: imagesPerRepositoryName,
					ResourceName: "ecr_repos",
					Description: imagesPerRepositoryDesc,
					Usage: 3,
				}
			}
		}
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &mockECRClient{
				err: nil,
				DescribeECRResponse: &ecr.DescribeRepositoriesOutput{
					Repositories tc.repositories,
				},
			}

			check := ImagesPerRepositoryUsageCheck{mockClient}
			usage, err := check.Usage()

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedUsage, usage)
		})
	}
}
