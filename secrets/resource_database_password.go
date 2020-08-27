package secrets

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

func resourceDatabasePassword() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatabasePasswordCreate,
		ReadContext:   resourceDatabasePasswordRead,
		UpdateContext: resourceDatabasePasswordRead,
		DeleteContext: resourceDatabasePasswordDelete,
		Schema: map[string]*schema.Schema{
			"secret_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "id of secret",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(60 * time.Second),
		},
	}
}

func resourceDatabasePasswordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	secretsManager := secretsmanager.New(getSession())

	secret := Password{
		StagingAdminPassword:           generateRandomPassword(secretsManager),
		StagingUserPassword:            generateRandomPassword(secretsManager),
		StagingReadOnlyUserPassword:    generateRandomPassword(secretsManager),
		ProductionAdminPassword:        generateRandomPassword(secretsManager),
		ProductionUserPassword:         generateRandomPassword(secretsManager),
		ProductionReadOnlyUserPassword: generateRandomPassword(secretsManager),
	}

	secretString, err := json.Marshal(secret)

	if err != nil {
		return diag.FromErr(err)
	}

	secretId := getSecretId(d)

	psvi := &secretsmanager.PutSecretValueInput{
		SecretId:     aws.String(secretId),
		SecretString: aws.String(string(secretString)),
	}

	_, err = secretsManager.PutSecretValue(psvi)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(secretId)

	return diags
}

func resourceDatabasePasswordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	d.SetId(getSecretId(d))

	return diags
}

func resourceDatabasePasswordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}