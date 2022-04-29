package servicequotas

import (
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
)

type mockECRClient struct {
	ecriface.ECRAPI

	err                 error
	DescribeECRResponse *ecr.DescribeRepositoriesOutput
}
