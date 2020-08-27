package secrets

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Password struct {
	StagingAdminPassword           string `json:"STAGING_ADMIN_PASSWORD"`
	StagingUserPassword            string `json:"STAGING_USER_PASSWORD"`
	StagingReadOnlyUserPassword    string `json:"STAGING_READONLY_USER_PASSWORD"`
	ProductionAdminPassword        string `json:"PRODUCTION_ADMIN_PASSWORD"`
	ProductionUserPassword         string `json:"PRODUCTION_USER_PASSWORD"`
	ProductionReadOnlyUserPassword string `json:"PRODUCTION_READONLY_USER_PASSWORD"`
}

func getSession() *session.Session {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	if err == nil {
		return sess
	} else {
		fmt.Println(err.Error())
		return nil
	}
}

func getSecretId(d *schema.ResourceData) string {
	return d.Get("secret_id").(string)
}

func generateRandomPassword(svc *secretsmanager.SecretsManager) string {
	gpi := &secretsmanager.GetRandomPasswordInput{
		ExcludePunctuation: aws.Bool(true),
		PasswordLength:     aws.Int64(32),
	}

	gpo, err := svc.GetRandomPassword(gpi)

	if err != nil {
		fmt.Println(err.Error())
		return "FAILED_TO_GENERATE_RANDOM_PASSWORD"
	}

	return *gpo.RandomPassword
}