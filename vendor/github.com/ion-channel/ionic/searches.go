package ionic

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/ion-channel/ionic/products"
)

const (
	searchByProductIdentifiersEndpoint = "v1/search/productidentifiers"
	searchByRepositoriesEndpoint       = "v1/search/repositories"
)

// GetSearchProductIdentifiers takes a product identifier, version and vendor to perform
// a productidentifier search against the Ion API, assembling a slice of Ionic
// products.ProductSearchResponse objects
func (ic *IonClient) GetSearchProductIdentifiers(productIdentifer, version, vendor, token string) ([]products.SoftwareEntity, error) {
	params := &url.Values{}
	params.Set("productidentifiers", productIdentifer)
	params.Set("vendor", vendor)
	params.Set("version", version)

	b, err := ic.Get(searchByProductIdentifiersEndpoint, token, params, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get productidentifiers search: %v", err.Error())
	}

	var results []products.SoftwareEntity
	err = json.Unmarshal(b, &results)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal product search results: %v", err.Error())
	}

	return results, nil
}

// GetSearchRepositories takes a repository string and performs
// a repository search against the Ion API, assembling a slice of Ionic
// products.ProductSearchResponse objects
func (ic *IonClient) GetSearchRepositories(repo, token string) ([]products.SoftwareEntity, error) {
	params := &url.Values{}
	params.Set("github", repo)

	b, err := ic.Get(searchByRepositoriesEndpoint, token, params, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get repo search: %v", err.Error())
	}
	var results []products.SoftwareEntity
	err = json.Unmarshal(b, &results)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal repo search results: %v", err.Error())
	}
	return results, nil
}
