package azuredevops

import (
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/serviceendpointcrud"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"

	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/tfhelper"
)

func resourceServiceEndpointGitHub() *schema.Resource {

	expand := func(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string) {
		serviceEndpoint, projectID := serviceendpointcrud.BaseExpander(d)
		serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"accessToken": d.Get("github_service_endpoint_pat").(string),
			},
			Scheme: converter.String("PersonalAccessToken"),
		}

		return serviceEndpoint, projectID
	}

	flatten := func(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
		serviceendpointcrud.BaseFlattener(d, serviceEndpoint, projectID)
		tfhelper.HelpFlattenSecret(d, "github_service_endpoint_pat")
		d.Set("github_service_endpoint_pat", (*serviceEndpoint.Authorization.Parameters)["accessToken"])
	}

	s := serviceendpointcrud.MakeBaseResource()

	s.Schema["github_service_endpoint_pat"] = &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		DefaultFunc:      schema.EnvDefaultFunc("AZDO_GITHUB_SERVICE_CONNECTION_PAT", nil),
		Description:      "The GitHub personal access token which should be used.",
		Sensitive:        true,
		DiffSuppressFunc: tfhelper.DiffFuncSupressSecretChanged,
	}

	patHashKey, patHashSchema := tfhelper.GenerateSecreteMemoSchema("github_service_endpoint_pat")
	s.Schema[patHashKey] = patHashSchema

	return &s
}
