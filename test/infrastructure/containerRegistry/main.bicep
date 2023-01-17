targetScope =  'resourceGroup'

param name string
param location string
param tags object
param skuName string
param properties object

resource cr 'Microsoft.ContainerRegistry/registries@2022-02-01-preview' = {
  name: name
  location: location
  tags: tags
  sku: {
    name: skuName
  }
  identity: null
  properties: properties
}
