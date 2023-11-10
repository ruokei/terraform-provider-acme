package tests

import (
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {

	tfOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/basic",
	})

	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, tfOptions)

	_, err := retry.DoWithRetryE(t, "Return value after 5 retries", 10, 5*time.Second, func() (string, error) {
		// Run "terraform init" and "terraform apply". Fail the test if there are any errors after 5 retries.
		terraform.InitAndApplyE(t, tfOptions)

		return "", nil
	})

	assert.Equal(t, err, nil)
}
