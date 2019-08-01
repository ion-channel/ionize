package ionic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	"github.com/ion-channel/ionic/dependencies"
)

const (
	// RubyEcosystem represents the ruby ecosystem for resolving dependencies
	RubyEcosystem = "ruby"
)

// ResolveDependenciesInFile takes a dependency file location and token to send
// the specified file to the API. All dependencies that are able to be resolved will
// be with their info returned, and a list of any errors encountered during the
// process.
func (ic *IonClient) ResolveDependenciesInFile(o dependencies.DependencyResolutionRequest, token string) (*dependencies.DependencyResolutionResponse, error) {
	params := &url.Values{}
	params.Set("type", o.Ecosystem)
	if o.Flatten {
		params.Set("flatten", "true")
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	fw, err := w.CreateFormFile("file", o.File)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %v", err.Error())
	}

	fh, err := os.Open(o.File)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err.Error())
	}

	_, err = io.Copy(fw, fh)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file contents: %v", err.Error())
	}

	w.Close()

	h := http.Header{}
	h.Set("Content-Type", w.FormDataContentType())

	b, err := ic.Post(dependencies.ResolveDependenciesInFileEndpoint, token, params, buf, h)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve dependencies: %v", err.Error())
	}

	var resp dependencies.DependencyResolutionResponse
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err.Error())
	}

	return &resp, nil
}

// GetLatestVersionForDependency takes a package name, an ecosystem to find the
// package in, and a token for accessing the API. It returns a dependency
// representation of the latest version and any errors it encounters with the
// API.
func (ic *IonClient) GetLatestVersionForDependency(packageName, ecosystem, token string) (*dependencies.Dependency, error) {
	params := &url.Values{}
	params.Set("name", packageName)
	params.Set("type", ecosystem)

	b, err := ic.Get(dependencies.GetLatestVersionForDependencyEndpoint, token, params, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest version for dependency: %v", err.Error())
	}

	var dep dependencies.Dependency
	err = json.Unmarshal(b, &dep)
	if err != nil {
		return nil, fmt.Errorf("cannot parse dependency: %v", err.Error())
	}

	dep.Name = packageName
	return &dep, nil
}

// GetVersionsForDependency takes a package name, an ecosystem to find the
// package in, and a token for accessing the API. It returns a dependency
// representation of the latest versions and any errors it encounters with the
// API.
func (ic *IonClient) GetVersionsForDependency(packageName, ecosystem, token string) ([]dependencies.Dependency, error) {
	params := &url.Values{}
	params.Set("name", packageName)
	params.Set("type", ecosystem)

	b, err := ic.Get(dependencies.GetVersionsForDependencyEndpoint, token, params, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest version for dependency: %v", err.Error())
	}

	var vs []string
	err = json.Unmarshal(b, &vs)
	if err != nil {
		return nil, fmt.Errorf("cannot parse dependency: %v", err.Error())
	}

	deps := []dependencies.Dependency{}
	for index := range vs {
		dep := dependencies.Dependency{
			Name:    packageName,
			Version: vs[index],
		}
		deps = append(deps, dep)
	}

	return deps, nil
}
