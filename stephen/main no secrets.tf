# Make sure to set the following environment variables:
#   AZDO_PERSONAL_ACCESS_TOKEN
#   AZDO_ORG_SERVICE_URL
#   AZDO_GITHUB_SERVICE_CONNECTION_PAT
provider "azuredevops" {
  version = ">= 0.0.1"
}

resource "azuredevops_project" "project" {
  project_name       = "Test Project"
  description        = "Test Project Description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_build_definition" "build_definition" {
  project_id      = azuredevops_project.project.id
  name            = "Test Pipeline"
  agent_pool_name = "Hosted Ubuntu 1604"

  repository {
    repo_type             = "GitHub"
    repo_name             = "nmiodice/terraform-azure-devops-hack"
    branch_name           =. "master"
    yml_path              = "azdo-api-samples/azure-pipeline.yml"
    service_connection_id = azuredevops_serviceendpoint_github.github_serviceendpoint.id
  }
}

resource "azuredevops_serviceendpoint_github" "github_serviceendpoint" {
  project_id      = azuredevops_project.project.id
  connection_name = "GitHub Service Connection"
  github_pat      = "authorization.parameters.accessToken"
}

resource "azuredevops_serviceendpoint_dockerhub" "dockerhub_serviceendpoint" {
  project_id      = azuredevops_project.project.id
  connection_name = "dockerhub test connection"
  docker_username = "authorization.parameters.accessToken"
  docker_password = "secret!"
  docker_email    = "nancy@example.com"
}

resource "azuredevops_serviceendpoint_azurerm" "azurerm_serviceendpoint" {
  project_id        = azuredevops_project.project.id
  connection_name   = "Example Deployment Admin"
  sp_id             = "aaaaaaaa-9999-8888-7777-bbbbbbbbbbbb"
  sp_key            = "bbbbbbbb-3333-4444-5555-cccccccccccc"
  tenant_id         = "cccccccc-7777-4444-aaaa-eeeeeeeeeeee"
  subscription_id   = "dddddddd-2452-4444-9ae2-333333333333"
  subscription_name = "Example Developer"
}

resource "azuredevops_serviceendpoint_acr" "acr_serviceendpoint" {
  project_id        = azuredevops_project.project.id
  connection_name   = "acr connection test"
  login_server      = "examplepp3isolatedserviceacr.azurecr.io"
  sp_role           = "8311e382-0749-4cb8-b61a-304f252e45ec"
  registry_id       = "/subscriptions/dddddddd-2452-4444-9ae2-333333333333/resourceGroups/example-service-app-rg/providers/Microsoft.ContainerRegistry/registries/examplepp3isolatedserviceacr"
  sp_id             = "00000000-8c8e-49ba-9999-4dfdc813c313"
  tenent_id         = "cccccccc-7777-4444-aaaa-eeeeeeeeeeee"
  subscription_id   = "dddddddd-2452-4444-9ae2-333333333333"
  subscription_name = "Example Developer"
}
