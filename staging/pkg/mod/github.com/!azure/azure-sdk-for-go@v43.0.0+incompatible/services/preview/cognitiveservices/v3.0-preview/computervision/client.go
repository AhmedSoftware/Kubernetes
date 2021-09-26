// Package computervision implements the Azure ARM Computervision service API version 3.0-preview.
//
// The Computer Vision API provides state-of-the-art algorithms to process images and return information. For example,
// it can be used to determine if an image contains mature content, or it can be used to find all the faces in an
// image.  It also has other features like estimating dominant and accent colors, categorizing the content of images,
// and describing an image with complete English sentences.  Additionally, it can also intelligently generate images
// thumbnails for displaying large images effectively.
package computervision

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
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/Azure/go-autorest/tracing"
	"github.com/satori/go.uuid"
	"io"
	"net/http"
)

// BaseClient is the base client for Computervision.
type BaseClient struct {
	autorest.Client
	Endpoint string
}

// New creates an instance of the BaseClient client.
func New(endpoint string) BaseClient {
	return NewWithoutDefaults(endpoint)
}

// NewWithoutDefaults creates an instance of the BaseClient client.
func NewWithoutDefaults(endpoint string) BaseClient {
	return BaseClient{
		Client:   autorest.NewClientWithUserAgent(UserAgent()),
		Endpoint: endpoint,
	}
}

// GetReadResult this interface is used for getting OCR results of Read operation. The URL to this interface should be
// retrieved from 'Operation-Location' field returned from Read interface.
// Parameters:
// operationID - id of read operation returned in the response of the 'Read' interface.
func (client BaseClient) GetReadResult(ctx context.Context, operationID uuid.UUID) (result ReadOperationResult, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/BaseClient.GetReadResult")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetReadResultPreparer(ctx, operationID)
	if err != nil {
		err = autorest.NewErrorWithError(err, "computervision.BaseClient", "GetReadResult", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetReadResultSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "computervision.BaseClient", "GetReadResult", resp, "Failure sending request")
		return
	}

	result, err = client.GetReadResultResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "computervision.BaseClient", "GetReadResult", resp, "Failure responding to request")
	}

	return
}

// GetReadResultPreparer prepares the GetReadResult request.
func (client BaseClient) GetReadResultPreparer(ctx context.Context, operationID uuid.UUID) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"Endpoint": client.Endpoint,
	}

	pathParameters := map[string]interface{}{
		"operationId": autorest.Encode("path", operationID),
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithCustomBaseURL("{Endpoint}/vision/v3.0-preview", urlParameters),
		autorest.WithPathParameters("/read/analyzeResults/{operationId}", pathParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetReadResultSender sends the GetReadResult request. The method will close the
// http.Response Body if it receives an error.
func (client BaseClient) GetReadResultSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// GetReadResultResponder handles the response to the GetReadResult request. The method always
// closes the http.Response Body.
func (client BaseClient) GetReadResultResponder(resp *http.Response) (result ReadOperationResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Read use this interface to get the result of a Read operation, employing the state-of-the-art Optical Character
// Recognition (OCR) algorithms optimized for text-heavy documents. When you use the Read interface, the response
// contains a field called 'Operation-Location'. The 'Operation-Location' field contains the URL that you must use for
// your 'GetReadResult' operation to access OCR results.​
// Parameters:
// imageURL - a JSON document with a URL pointing to the image that is to be analyzed.
// language - the BCP-47 language code of the text to be detected in the image. In future versions, when
// language parameter is not passed, language detection will be used to determine the language. However, in the
// current version, missing language parameter will cause English to be used. To ensure that your document is
// always parsed in English without the use of language detection in the future, pass “en” in the language
// parameter
func (client BaseClient) Read(ctx context.Context, imageURL ImageURL, language OcrDetectionLanguage) (result autorest.Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/BaseClient.Read")
		defer func() {
			sc := -1
			if result.Response != nil {
				sc = result.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: imageURL,
			Constraints: []validation.Constraint{{Target: "imageURL.URL", Name: validation.Null, Rule: true, Chain: nil}}}}); err != nil {
		return result, validation.NewError("computervision.BaseClient", "Read", err.Error())
	}

	req, err := client.ReadPreparer(ctx, imageURL, language)
	if err != nil {
		err = autorest.NewErrorWithError(err, "computervision.BaseClient", "Read", nil, "Failure preparing request")
		return
	}

	resp, err := client.ReadSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "computervision.BaseClient", "Read", resp, "Failure sending request")
		return
	}

	result, err = client.ReadResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "computervision.BaseClient", "Read", resp, "Failure responding to request")
	}

	return
}

// ReadPreparer prepares the Read request.
func (client BaseClient) ReadPreparer(ctx context.Context, imageURL ImageURL, language OcrDetectionLanguage) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"Endpoint": client.Endpoint,
	}

	queryParameters := map[string]interface{}{}
	if len(string(language)) > 0 {
		queryParameters["language"] = autorest.Encode("query", language)
	} else {
		queryParameters["language"] = autorest.Encode("query", "en")
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithCustomBaseURL("{Endpoint}/vision/v3.0-preview", urlParameters),
		autorest.WithPath("/read/analyze"),
		autorest.WithJSON(imageURL),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ReadSender sends the Read request. The method will close the
// http.Response Body if it receives an error.
func (client BaseClient) ReadSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// ReadResponder handles the response to the Read request. The method always
// closes the http.Response Body.
func (client BaseClient) ReadResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByClosing())
	result.Response = resp
	return
}

// ReadInStream use this interface to get the result of a Read operation, employing the state-of-the-art Optical
// Character Recognition (OCR) algorithms optimized for text-heavy documents. When you use the Read interface, the
// response contains a field called 'Operation-Location'. The 'Operation-Location' field contains the URL that you must
// use for your 'GetReadResult' operation to access OCR results.​
// Parameters:
// imageParameter - an image stream.
// language - the BCP-47 language code of the text to be detected in the image. In future versions, when
// language parameter is not passed, language detection will be used to determine the language. However, in the
// current version, missing language parameter will cause English to be used. To ensure that your document is
// always parsed in English without the use of language detection in the future, pass “en” in the language
// parameter
func (client BaseClient) ReadInStream(ctx context.Context, imageParameter io.ReadCloser, language OcrDetectionLanguage) (result autorest.Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/BaseClient.ReadInStream")
		defer func() {
			sc := -1
			if result.Response != nil {
				sc = result.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.ReadInStreamPreparer(ctx, imageParameter, language)
	if err != nil {
		err = autorest.NewErrorWithError(err, "computervision.BaseClient", "ReadInStream", nil, "Failure preparing request")
		return
	}

	resp, err := client.ReadInStreamSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "computervision.BaseClient", "ReadInStream", resp, "Failure sending request")
		return
	}

	result, err = client.ReadInStreamResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "computervision.BaseClient", "ReadInStream", resp, "Failure responding to request")
	}

	return
}

// ReadInStreamPreparer prepares the ReadInStream request.
func (client BaseClient) ReadInStreamPreparer(ctx context.Context, imageParameter io.ReadCloser, language OcrDetectionLanguage) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"Endpoint": client.Endpoint,
	}

	queryParameters := map[string]interface{}{}
	if len(string(language)) > 0 {
		queryParameters["language"] = autorest.Encode("query", language)
	} else {
		queryParameters["language"] = autorest.Encode("query", "en")
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/octet-stream"),
		autorest.AsPost(),
		autorest.WithCustomBaseURL("{Endpoint}/vision/v3.0-preview", urlParameters),
		autorest.WithPath("/read/analyze"),
		autorest.WithFile(imageParameter),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ReadInStreamSender sends the ReadInStream request. The method will close the
// http.Response Body if it receives an error.
func (client BaseClient) ReadInStreamSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// ReadInStreamResponder handles the response to the ReadInStream request. The method always
// closes the http.Response Body.
func (client BaseClient) ReadInStreamResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByClosing())
	result.Response = resp
	return
}
