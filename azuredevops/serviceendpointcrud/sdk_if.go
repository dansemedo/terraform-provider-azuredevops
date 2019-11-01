package serviceendpointcrud

import (
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
)

// Make the Azure DevOps API call to create the endpoint
func createServiceEndpoint(clients *aggregatedClient, endpoint *serviceendpoint.ServiceEndpoint, project *string) (*serviceendpoint.ServiceEndpoint, error) {
	createdServiceEndpoint, err := clients.ServiceEndpointClient.CreateServiceEndpoint(
		clients.ctx,
		serviceendpoint.CreateServiceEndpointArgs{
			Endpoint: endpoint,
			Project:  project,
		})

	return createdServiceEndpoint, err
}

func deleteServiceEndpoint(clients *aggregatedClient, project *string, endPointID *uuid.UUID) error {
	err := clients.ServiceEndpointClient.DeleteServiceEndpoint(
		clients.ctx,
		serviceendpoint.DeleteServiceEndpointArgs{
			Project:    project,
			EndpointId: endPointID,
		})

	return err
}

func updateServiceEndpoint(clients *aggregatedClient, endpoint *serviceendpoint.ServiceEndpoint, project *string) (*serviceendpoint.ServiceEndpoint, error) {
	updatedServiceEndpoint, err := clients.ServiceEndpointClient.UpdateServiceEndpoint(
		clients.ctx,
		serviceendpoint.UpdateServiceEndpointArgs{
			Endpoint:   endpoint,
			Project:    project,
			EndpointId: endpoint.Id,
		})

	return updatedServiceEndpoint, err
}
