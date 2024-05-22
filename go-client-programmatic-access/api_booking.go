/*
Qdrant Cloud API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.1.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package qdrant_cloud_programmatic_access

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
)


type BookingAPI interface {

	/*
	GetPackages Http Get Packages

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiGetPackagesRequest
	*/
	GetPackages(ctx context.Context) ApiGetPackagesRequest

	// GetPackagesExecute executes the request
	//  @return []PackageOut
	GetPackagesExecute(r ApiGetPackagesRequest) ([]PackageOut, *http.Response, error)
}

// BookingAPIService BookingAPI service
type BookingAPIService service

type ApiGetPackagesRequest struct {
	ctx context.Context
	ApiService BookingAPI
	provider *string
	region *string
}

func (r ApiGetPackagesRequest) Provider(provider string) ApiGetPackagesRequest {
	r.provider = &provider
	return r
}

func (r ApiGetPackagesRequest) Region(region string) ApiGetPackagesRequest {
	r.region = &region
	return r
}

func (r ApiGetPackagesRequest) Execute() ([]PackageOut, *http.Response, error) {
	return r.ApiService.GetPackagesExecute(r)
}

/*
GetPackages Http Get Packages

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiGetPackagesRequest
*/
func (a *BookingAPIService) GetPackages(ctx context.Context) ApiGetPackagesRequest {
	return ApiGetPackagesRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
//  @return []PackageOut
func (a *BookingAPIService) GetPackagesExecute(r ApiGetPackagesRequest) ([]PackageOut, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  []PackageOut
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "BookingAPIService.GetPackages")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/booking/packages"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.provider != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "provider", r.provider, "")
	}
	if r.region != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "region", r.region, "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	if r.ctx != nil {
		// API Key Authentication
		if auth, ok := r.ctx.Value(ContextAPIKeys).(map[string]APIKey); ok {
			if apiKey, ok := auth["ApiKeyAuth"]; ok {
				var key string
				if apiKey.Prefix != "" {
					key = apiKey.Prefix + " " + apiKey.Key
				} else {
					key = apiKey.Key
				}
				localVarHeaderParams["Authorization"] = key
			}
		}
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if localVarHTTPResponse.StatusCode == 422 {
			var v HTTPValidationError
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}