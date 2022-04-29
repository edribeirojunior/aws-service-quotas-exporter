package servicequotas

import (
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/pkg/errors"
)

// Not all quota limits here are reported under "ecr", but all of the
// usage checks are using the ecr service
const (
	imagesPerRepositoryName = "images_per_repository"
	imagesPerRepositoryDesc = "Images per Repository"
)

// ImagesPerRespositoryUsageCheck implements the UsageCheck interface
// for images per repos
type ImagesPerRepositoryUsageCheck struct {
	client ecriface.ECRAPI
}

// Usage returns the usage for each image with the usage
// value being the sum of their inbound and outbound rules or an error
func (c *ImagesPerRepositoryUsageCheck) Usage() ([]QuotaUsage, error) {
	quotaUsages := []QuotaUsage{}

	params := &ecr.DescribeRepositoriesInput{}
	err := c.client.DescribeRepositoriesPages(params,
		func(page *ecr.DescribeRepositoriesOutput, lastPage bool) bool {
			if page != nil {
				//tags := map[string]string{}

				rName := "ecr_repos"

				reposUsage := QuotaUsage{
					Name:         imagesPerRepositoryName,
					ResourceName: &rName,
					Description:  imagesPerRepositoryDesc,
					Usage:        float64(len(page.Repositories)),
					//Tags:         tags,
				}

				quotaUsages = append(quotaUsages, []QuotaUsage{reposUsage}...)
			}
			return !lastPage
		},
	)
	if err != nil {
		return nil, errors.Wrapf(ErrFailedToGetUsage, "%w", err)
	}

	return quotaUsages, nil
}
