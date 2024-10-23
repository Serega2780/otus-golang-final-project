package integration

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cucumber/godog"
)

type godogsResponseCtxKey struct{}

type previewerTest struct {
	url    string
	client *http.Client
}

type response struct {
	status  int
	body    any
	headers http.Header
}

func InitializeScenario(godogCtx *godog.ScenarioContext) {
	tr := &http.Transport{
		MaxIdleConns:    100,
		IdleConnTimeout: 10 * time.Second,
	}

	test := &previewerTest{
		url: "http://previewer:8585",
		client: &http.Client{
			Transport: tr,
		},
	}
	godogCtx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		return context.WithValue(ctx, godogsResponseCtxKey{}, &response{status: 0, body: nil}), nil
	})
	godogCtx.Step(`^I send GET request to "([^"]*)" for non-existing server nginx2`, test.iServerNotExists)
	godogCtx.Step(`^the response code should be (\d+)$`, test.theResponseCodeShouldBe)
	godogCtx.Step(`^the response payload must contain "([^"]*)"`, test.iServerNotExistsCheckResponse)

	godogCtx.Step(`^I send GET request to "([^"]*)" for non-existing image _gopher_fake_1024x504.jpg`,
		test.iImageNotExists)
	godogCtx.Step(`^the response code should be (\d+)$`, test.theResponseCodeShouldBe)
	godogCtx.Step(`^the response payload must contain "([^"]*)"`, test.iImageNotExistsCheckResponse)

	godogCtx.Step(`^I send GET request to "([^"]*)" for an image 7z_100x100.jpg`, test.iImageWrongFormat)
	godogCtx.Step(`^the response code should be (\d+)$`, test.theResponseCodeShouldBe)
	godogCtx.Step(`^the response payload must contain "([^"]*)"`, test.iImageWrongFormatCheckResponse)

	godogCtx.Step(`^I send GET request to "([^"]*)" for an existing image _gopher_original_1024x504.jpg`,
		test.iImageFoundRemote)
	godogCtx.Step(`^the response code should be (\d+)$`, test.theResponseCodeShouldBe)
	godogCtx.Step(`^the response payload must be a valid JPEG image`, test.iImageFoundRemoteCheckResponse)

	godogCtx.Step(`^I send GET request to "([^"]*)" for an image gopher_50x50.jpg`, test.iImageWrongDimensions)
	godogCtx.Step(`^the response code should be (\d+)$`, test.theResponseCodeShouldBe)
	godogCtx.Step(`^the response payload must contain "([^"]*)"`, test.iImageWrongDimensionsCheckResponse)

	godogCtx.Step(`^I send GET request to "([^"]*)" for an existing image _gopher_original_1024x504.jpg`,
		test.iImageFoundLocal)
	godogCtx.Step(`^the response code should be (\d+)$`, test.theResponseCodeShouldBe)
	godogCtx.Step(`^the response payload must be a valid JPEG image`, test.iImageFoundLocalCheckResponse)
}

func (pt *previewerTest) iImageFoundLocal(ctx context.Context, route string) (context.Context, error) {
	req, _ := http.NewRequestWithContext(ctx, "", pt.url+route, nil)
	req.Header.Set("Authorization", "Basic ==")
	res, _ := pt.client.Do(req) //nolint:bodyclose
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	res2, _ := pt.client.Do(req) //nolint:bodyclose
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res2.Body)
	b, _ := io.ReadAll(res2.Body)
	actual := response{
		status:  res.StatusCode,
		body:    b,
		headers: res2.Header,
	}

	return context.WithValue(ctx, godogsResponseCtxKey{}, &actual), nil
}

func (pt *previewerTest) iImageWrongDimensions(ctx context.Context, route string) (context.Context, error) {
	req, _ := http.NewRequestWithContext(ctx, "", pt.url+route, nil)
	req.Header.Set("Authorization", "Basic ==")
	res, _ := pt.client.Do(req) //nolint:bodyclose
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	b, _ := io.ReadAll(res.Body)
	actual := response{
		status: res.StatusCode,
		body:   b,
	}

	return context.WithValue(ctx, godogsResponseCtxKey{}, &actual), nil
}

func (pt *previewerTest) iImageFoundRemote(ctx context.Context, route string) (context.Context, error) {
	req, _ := http.NewRequestWithContext(ctx, "", pt.url+route, nil)
	req.Header.Set("Authorization", "Basic ==")
	res, _ := pt.client.Do(req) //nolint:bodyclose
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	b, _ := io.ReadAll(res.Body)
	actual := response{
		status: res.StatusCode,
		body:   b,
	}

	return context.WithValue(ctx, godogsResponseCtxKey{}, &actual), nil
}

func (pt *previewerTest) iImageWrongFormat(ctx context.Context, route string) (context.Context, error) {
	req, _ := http.NewRequestWithContext(ctx, "", pt.url+route, nil)
	req.Header.Set("Authorization", "Basic ==")
	res, _ := pt.client.Do(req) //nolint:bodyclose
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	b, _ := io.ReadAll(res.Body)
	actual := response{
		status: res.StatusCode,
		body:   b,
	}

	return context.WithValue(ctx, godogsResponseCtxKey{}, &actual), nil
}

func (pt *previewerTest) iImageNotExists(ctx context.Context, route string) (context.Context, error) {
	req, _ := http.NewRequestWithContext(ctx, "", pt.url+route, nil)
	req.Header.Set("Authorization", "Basic ==")
	res, _ := pt.client.Do(req) //nolint:bodyclose
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	b, _ := io.ReadAll(res.Body)
	actual := response{
		status: res.StatusCode,
		body:   b,
	}

	return context.WithValue(ctx, godogsResponseCtxKey{}, &actual), nil
}

func (pt *previewerTest) iServerNotExists(ctx context.Context, route string) (context.Context, error) {
	req, _ := http.NewRequestWithContext(ctx, "", pt.url+route, nil)
	req.Header.Set("Authorization", "Basic ==")
	res, _ := pt.client.Do(req) //nolint:bodyclose
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	b, _ := io.ReadAll(res.Body)
	actual := response{
		status: res.StatusCode,
		body:   b,
	}

	return context.WithValue(ctx, godogsResponseCtxKey{}, &actual), nil
}

func (pt *previewerTest) theResponseCodeShouldBe(ctx context.Context, expectedStatus int) error {
	resp, ok := ctx.Value(godogsResponseCtxKey{}).(*response)
	if !ok {
		return errors.New("there are no godogs available")
	}

	if expectedStatus != resp.status {
		return fmt.Errorf("expected response code to be: %d, but actual is: %d", expectedStatus, resp.status)
	}

	return nil
}

func (pt *previewerTest) iServerNotExistsCheckResponse(ctx context.Context, subStr string) error {
	actualResp, ok := ctx.Value(godogsResponseCtxKey{}).(*response)
	if !ok {
		return errors.New("there are no godogs available")
	}
	body := string(actualResp.body.([]byte))

	if !strings.Contains(body, subStr) {
		return fmt.Errorf("response %v does not contain %v", body, subStr)
	}

	return nil
}

func (pt *previewerTest) iImageNotExistsCheckResponse(ctx context.Context, subStr string) error {
	actualResp, ok := ctx.Value(godogsResponseCtxKey{}).(*response)
	if !ok {
		return errors.New("there are no godogs available")
	}
	body := string(actualResp.body.([]byte))

	if !strings.Contains(body, subStr) {
		return fmt.Errorf("response %v does not contain %v", body, subStr)
	}

	return nil
}

func (pt *previewerTest) iImageWrongFormatCheckResponse(ctx context.Context, subStr string) error {
	actualResp, ok := ctx.Value(godogsResponseCtxKey{}).(*response)
	if !ok {
		return errors.New("there are no godogs available")
	}
	body := string(actualResp.body.([]byte))

	if !strings.Contains(body, subStr) {
		return fmt.Errorf("response %v does not contain %v", body, subStr)
	}

	return nil
}

func (pt *previewerTest) iImageFoundRemoteCheckResponse(ctx context.Context) error {
	actualResp, ok := ctx.Value(godogsResponseCtxKey{}).(*response)
	if !ok {
		return errors.New("there are no godogs available")
	}
	body := actualResp.body.([]byte)
	_, err := jpeg.Decode(bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("response is not a valid JPEG image %w", err)
	}

	return nil
}

func (pt *previewerTest) iImageWrongDimensionsCheckResponse(ctx context.Context, subStr string) error {
	actualResp, ok := ctx.Value(godogsResponseCtxKey{}).(*response)
	if !ok {
		return errors.New("there are no godogs available")
	}
	body := string(actualResp.body.([]byte))

	if !strings.Contains(body, subStr) {
		return fmt.Errorf("response %v does not contain %v", body, subStr)
	}

	return nil
}

func (pt *previewerTest) iImageFoundLocalCheckResponse(ctx context.Context) error {
	actualResp, ok := ctx.Value(godogsResponseCtxKey{}).(*response)
	if !ok {
		return errors.New("there are no godogs available")
	}
	body := actualResp.body.([]byte)
	_, err := jpeg.Decode(bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("response is not a valid JPEG image %w", err)
	}
	hit := actualResp.headers.Get("X-Previewer-Cache-Hit")
	if hit != "true" {
		return fmt.Errorf("cache hit header must be true, received %v", hit)
	}

	return nil
}
