# \BookingAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetPackages**](BookingAPI.md#GetPackages) | **Get** /booking/packages | Http Get Packages



## GetPackages

> []PackageOut GetPackages(ctx).Provider(provider).Region(region).Execute()

Http Get Packages

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/qdrant/qdrant-cloud-cluster-api/pypi/go-client-programmatic-access"
)

func main() {
	provider := "provider_example" // string |  (optional)
	region := "region_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BookingAPI.GetPackages(context.Background()).Provider(provider).Region(region).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BookingAPI.GetPackages``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetPackages`: []PackageOut
	fmt.Fprintf(os.Stdout, "Response from `BookingAPI.GetPackages`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetPackagesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **provider** | **string** |  | 
 **region** | **string** |  | 

### Return type

[**[]PackageOut**](PackageOut.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

