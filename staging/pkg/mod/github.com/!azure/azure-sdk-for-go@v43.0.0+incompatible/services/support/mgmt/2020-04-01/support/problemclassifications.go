package support

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// ProblemClassificationsClient is the microsoft Azure Support Resource Provider.
type ProblemClassificationsClient struct {
	BaseClient
}

// NewProblemClassificationsClient creates an instance of the ProblemClassificationsClient client.
func NewProblemClassificationsClient(subscriptionID string) ProblemClassificationsClient {
	return NewProblemClassificationsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewProblemClassificationsClientWithBaseURI creates an instance of the ProblemClassificationsClient client using a
// custom endpoint.  Use this when interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds,
// Azure stack).
func NewProblemClassificationsClientWithBaseURI(baseURI string, subscriptionID string) ProblemClassificationsClient {
	return ProblemClassificationsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Get get problem classification details for a specific Azure service.
// Parameters:
// serviceName - name of the Azure service available for support.
// problemClassificationName - name of problem classification.
func (client ProblemClassificationsClient) Get(ctx context.Context, serviceName string, problemClassificationName string) (result ProblemClassification, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ProblemClassificationsClient.Get")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetPreparer(ctx, serviceName, problemClassificationName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "support.ProblemClassificationsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "support.ProblemClassificationsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "support.ProblemClassificationsClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client ProblemClassificationsClient) GetPreparer(ctx context.Context, serviceName string, problemClassificationName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"problemClassificationName": autorest.Encode("path", problemClassificationName),
		"serviceName":               autorest.Encode("path", serviceName),
	}

	const APIVersion = "2020-04-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/providers/Microsoft.Support/services/{serviceName}/problemClassifications/{problemClassificationName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client ProblemClassificationsClient) GetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client ProblemClassificationsClient) GetResponder(resp *http.Response) (result ProblemClassification, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// List lists all the problem classifications (categories) available for a specific Azure service. Always use the
// service and problem classifications obtained programmatically. This practice ensures that you always have the most
// recent set of service and problem classification Ids.
// Parameters:
// serviceName - name of the Azure service for which the problem classifications need to be retrieved.
func (client ProblemClassificationsClient) List(ctx context.Context, serviceName string) (result ProblemClassificationsListResult, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ProblemClassificationsClient.List")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.ListPreparer(ctx, serviceName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "support.ProblemClassificationsClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "support.ProblemClassificationsClient", "List", resp, "Failure sending request")
		return
	}

	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "support.ProblemClassificationsClient", "List", resp, "Failure responding to request")
	}

	return
}

// ListPreparer prepares the List request.
func (client ProblemClassificationsClient) ListPreparer(ctx context.Context, serviceName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"serviceName": autorest.Encode("path", serviceName),
	}

	const APIVersion = "2020-04-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/providers/Microsoft.Support/services/{serviceName}/problemClassifications", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client ProblemClassificationsClient) ListSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client ProblemClassificationsClient) ListResponder(resp *http.Response) (result ProblemClassificationsListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
