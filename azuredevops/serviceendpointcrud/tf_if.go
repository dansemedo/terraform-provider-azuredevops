package serviceendpointcrud

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
)

type aggregatedClient struct {
	ctx                   context.Context
	ServiceEndpointClient serviceendpoint.Client
}

// MarshallingFuncs houses the funcs needed to transform Terraform and AzDO structures
type MarshallingFuncs struct {
	Flatten func(*schema.ResourceData, *serviceendpoint.ServiceEndpoint, *string)
	Expand  func(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string)
}

// CrudFuncs houses the funcs the are generatored here.
type CrudFuncs struct {
	Create schema.CreateFunc
	Read   schema.ReadFunc
	Update schema.UpdateFunc
	Delete schema.DeleteFunc
}

// GenerateCrudFuncs creates a set of schema CRUD functions to manage a Service Endpoint, leveraging
// customized FlattenerFuncs and ExpanderFuncs.
func GenerateCrudFuncs(fns *MarshallingFuncs) *CrudFuncs {
	create := makeCreaterFunc(fns)
	read := makeReaderFunc(fns)
	update := makeUpdaterFunc(fns)
	delete := makeDeleterFunc(fns)
	return &CrudFuncs{create, read, update, delete}
}

// MakeBaseResource makes  a new Schema struct that suitable for all
// Service Endpoints to start with
func MakeBaseResource() schema.Resource {
	return schema.Resource{}
}

// BaseExpander performs the expansions work that all service endpoints will need
func BaseExpander(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string) {
	// an "error" is OK here as it is expected in the case that the ID is not set in the resource data
	var serviceEndpointID *uuid.UUID
	parsedID, err := uuid.Parse(d.Id())
	if err == nil {
		serviceEndpointID = &parsedID
	}
	log.Printf("Updating github_service_endpoint_pat to %s", d.Get("github_service_endpoint_pat").(string))
	projectID := converter.String(d.Get("project_id").(string))
	serviceEndpoint := &serviceendpoint.ServiceEndpoint{
		Id:    serviceEndpointID,
		Name:  converter.String(d.Get("service_endpoint_name").(string)),
		Type:  converter.String(d.Get("service_endpoint_type").(string)),
		Url:   converter.String(d.Get("service_endpoint_url").(string)),
		Owner: converter.String(d.Get("service_endpoint_owner").(string)),
		Authorization: &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"accessToken": d.Get("github_service_endpoint_pat").(string),
			},
			Scheme: converter.String("PersonalAccessToken"),
		},
	}

	return serviceEndpoint, projectID
}

// BaseFlattener performs the flattening chores that all service endpoitns share
func BaseFlattener(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	d.SetId(serviceEndpoint.Id.String())
	d.Set("service_endpoint_name", *serviceEndpoint.Name)
	d.Set("service_endpoint_type", *serviceEndpoint.Type)
	d.Set("service_endpoint_url", *serviceEndpoint.Url)
	d.Set("service_endpoint_owner", *serviceEndpoint.Owner)
	d.Set("project_id", projectID)
}

func makeCreaterFunc(fns *MarshallingFuncs) schema.CreateFunc {
	return func(d *schema.ResourceData, m interface{}) error {
		clients := m.(*aggregatedClient)
		serviceEndpoint, projectID := fns.Expand(d)

		createdServiceEndpoint, err := createServiceEndpoint(clients, serviceEndpoint, projectID)
		if err != nil {
			return fmt.Errorf("Error creating service endpoint in Azure DevOps: %+v", err)
		}
		fns.Flatten(d, createdServiceEndpoint, projectID)
		return nil
	}
}

func makeReaderFunc(fns *MarshallingFuncs) schema.ReadFunc {
	return func(d *schema.ResourceData, m interface{}) error {
		clients := m.(*aggregatedClient)

		var serviceEndpointID *uuid.UUID
		parsedServiceEndpointID, err := uuid.Parse(d.Id())
		if err != nil {
			return fmt.Errorf("Error parsing the service endpoint ID from the Terraform resource data: %v", err)
		}
		serviceEndpointID = &parsedServiceEndpointID
		projectID := converter.String(d.Get("project_id").(string))

		serviceEndpoint, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(
			clients.ctx,
			serviceendpoint.GetServiceEndpointDetailsArgs{
				EndpointId: serviceEndpointID,
				Project:    projectID,
			},
		)
		if err != nil {
			return fmt.Errorf("Error looking up service endpoint given ID (%v) and project ID (%v): %v", serviceEndpointID, projectID, err)
		}

		fns.Flatten(d, serviceEndpoint, projectID)
		return nil
	}
}

func makeUpdaterFunc(fns *MarshallingFuncs) schema.UpdateFunc {
	return func(d *schema.ResourceData, m interface{}) error {
		clients := m.(*aggregatedClient)
		serviceEndpoint, projectID := fns.Expand(d)

		updatedServiceEndpoint, err := updateServiceEndpoint(clients, serviceEndpoint, projectID)
		if err != nil {
			return fmt.Errorf("Error updating service endpoint in Azure DevOps: %+v", err)
		}

		fns.Flatten(d, updatedServiceEndpoint, projectID)
		return nil
	}
}

func makeDeleterFunc(fns *MarshallingFuncs) schema.DeleteFunc {
	return func(d *schema.ResourceData, m interface{}) error {
		clients := m.(*aggregatedClient)
		serviceEndpoint, projectID := fns.Expand(d)

		return deleteServiceEndpoint(clients, projectID, serviceEndpoint.Id)
	}
}
