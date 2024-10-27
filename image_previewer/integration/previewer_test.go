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
		url: "http://localhost:8585",
		client: &http.Client{
			Transport: tr,
		},
	}

	godogCtx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		return context.WithValue(ctx, godogsResponseCtxKey{}, &response{}), nil
	})

	godogCtx.When(`^I send GET request to "([^"]*)" for an existing image _gopher_original_1024x504.jpg`,
		test.iImageFoundLocal)
	godogCtx.Then(`^the response code (\d+) and the response header X-Previewer-Cache-Hit is "([^"]*)"`,
		test.iImageFoundLocalCheckResponse)

	godogCtx.When(`^I send GET request to "([^"]*)" for non-existing server nginx2`, test.iServerNotExists)
	godogCtx.Then(`^the response code (\d+) and the response payload "([^"]*)" for a non-existing server`,
		test.iServerNotExistsCheckResponse)

	godogCtx.When(`^I send GET request to "([^"]*)" for non-existing image _gopher_fake_1024x504.jpg`,
		test.iImageNotExists)
	godogCtx.Then(`^the response code (\d+) and the response payload "([^"]*)" for a non-existing image`,
		test.iImageNotExistsCheckResponse)

	godogCtx.When(`^I send GET request to "([^"]*)" for an image 7z_100x100.jpg`, test.iImageWrongFormat)
	godogCtx.Then(`^the response code (\d+) and the response payload "([^"]*)" for wrong format image`,
		test.iImageWrongFormatCheckResponse)

	godogCtx.When(`^I send GET request to "([^"]*)" for an existing image gopher_2000x1000.jpg`,
		test.iImageFoundRemote)
	godogCtx.Then(`^the response code (\d+) and the response payload must be JPEG image found remotely`,
		test.iImageFoundRemoteCheckResponse)

	godogCtx.When(`^I send GET request to "([^"]*)" for an image gopher_50x50.jpg`, test.iImageResizeUp)
	godogCtx.Then(`^the response code (\d+) and the res payload must have width (\d+) and height (\d+)$`,
		test.iImageResizeUpCheckResponse)

	godogCtx.When(`^I send GET request to "([^"]*)" for an existing image gopher_333x666.jpg`,
		test.iImageHeaderAbsence)
	godogCtx.Then(`^the response code (\d+), the response payload "([^"]*)", because of auth header absence`,
		test.iImageHeaderAbsenceCheckResponse)
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
	actual := &response{
		status:  res.StatusCode,
		body:    b,
		headers: res2.Header,
	}

	return context.WithValue(ctx, godogsResponseCtxKey{}, actual), nil
}

func (pt *previewerTest) iImageResizeUp(ctx context.Context, route string) (context.Context, error) {
	req, _ := http.NewRequestWithContext(ctx, "", pt.url+route, nil)
	req.Header.Set("Authorization", "Basic ==")
	res, _ := pt.client.Do(req) //nolint:bodyclose
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	b, _ := io.ReadAll(res.Body)
	actual := &response{
		status:  res.StatusCode,
		body:    b,
		headers: res.Header,
	}

	return context.WithValue(ctx, godogsResponseCtxKey{}, actual), nil
}

func (pt *previewerTest) iImageFoundRemote(ctx context.Context, route string) (context.Context, error) {
	req, _ := http.NewRequestWithContext(ctx, "", pt.url+route, nil)
	req.Header.Set("Authorization", "Basic ==")
	res, _ := pt.client.Do(req) //nolint:bodyclose
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	b, _ := io.ReadAll(res.Body)
	actual := &response{
		status:  res.StatusCode,
		body:    b,
		headers: res.Header,
	}

	return context.WithValue(ctx, godogsResponseCtxKey{}, actual), nil
}

func (pt *previewerTest) iImageWrongFormat(ctx context.Context, route string) (context.Context, error) {
	req, _ := http.NewRequestWithContext(ctx, "", pt.url+route, nil)
	req.Header.Set("Authorization", "Basic ==")
	res, _ := pt.client.Do(req) //nolint:bodyclose
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	b, _ := io.ReadAll(res.Body)
	actual := &response{
		status:  res.StatusCode,
		body:    b,
		headers: res.Header,
	}

	return context.WithValue(ctx, godogsResponseCtxKey{}, actual), nil
}

func (pt *previewerTest) iImageNotExists(ctx context.Context, route string) (context.Context, error) {
	req, _ := http.NewRequestWithContext(ctx, "", pt.url+route, nil)
	req.Header.Set("Authorization", "Basic ==")
	res, _ := pt.client.Do(req) //nolint:bodyclose
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	b, _ := io.ReadAll(res.Body)
	actual := &response{
		status:  res.StatusCode,
		body:    b,
		headers: res.Header,
	}

	return context.WithValue(ctx, godogsResponseCtxKey{}, actual), nil
}

func (pt *previewerTest) iServerNotExists(ctx context.Context, route string) (context.Context, error) {
	req, _ := http.NewRequestWithContext(ctx, "", pt.url+route, nil)
	req.Header.Set("Authorization", "Basic ==")
	res, _ := pt.client.Do(req) //nolint:bodyclose
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	b, _ := io.ReadAll(res.Body)
	actual := &response{
		status:  res.StatusCode,
		body:    b,
		headers: res.Header,
	}

	return context.WithValue(ctx, godogsResponseCtxKey{}, actual), nil
}

func (pt *previewerTest) iImageHeaderAbsence(ctx context.Context, route string) (context.Context, error) {
	req, _ := http.NewRequestWithContext(ctx, "", pt.url+route, nil)
	res, _ := pt.client.Do(req) //nolint:bodyclose
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	b, _ := io.ReadAll(res.Body)
	actual := &response{
		status:  res.StatusCode,
		body:    b,
		headers: res.Header,
	}

	return context.WithValue(ctx, godogsResponseCtxKey{}, actual), nil
}

func (pt *previewerTest) iServerNotExistsCheckResponse(ctx context.Context, expectedStatus int, subStr string) error {
	actualResp, ok := ctx.Value(godogsResponseCtxKey{}).(*response)
	if !ok {
		return errors.New("there are no godogs available")
	}
	if err := checkStatus(expectedStatus, actualResp.status); err != nil {
		return err
	}

	body := string(actualResp.body.([]byte))

	if !strings.Contains(body, subStr) {
		return fmt.Errorf("response %v does not contain %v", body, subStr)
	}

	return nil
}

func (pt *previewerTest) iImageNotExistsCheckResponse(ctx context.Context, expectedStatus int, subStr string) error {
	actualResp, ok := ctx.Value(godogsResponseCtxKey{}).(*response)
	if !ok {
		return errors.New("there are no godogs available")
	}
	if err := checkStatus(expectedStatus, actualResp.status); err != nil {
		return err
	}

	body := string(actualResp.body.([]byte))

	if !strings.Contains(body, subStr) {
		return fmt.Errorf("response %v does not contain %v", body, subStr)
	}

	return nil
}

func (pt *previewerTest) iImageWrongFormatCheckResponse(ctx context.Context, expectedStatus int, subStr string) error {
	actualResp, ok := ctx.Value(godogsResponseCtxKey{}).(*response)
	if !ok {
		return errors.New("there are no godogs available")
	}
	if err := checkStatus(expectedStatus, actualResp.status); err != nil {
		return err
	}

	body := string(actualResp.body.([]byte))

	if !strings.Contains(body, subStr) {
		return fmt.Errorf("response %v does not contain %v", body, subStr)
	}

	return nil
}

func (pt *previewerTest) iImageFoundRemoteCheckResponse(ctx context.Context, expectedStatus int) error {
	actualResp, ok := ctx.Value(godogsResponseCtxKey{}).(*response)
	if !ok {
		return errors.New("there are no godogs available")
	}
	if err := checkStatus(expectedStatus, actualResp.status); err != nil {
		return err
	}

	body := actualResp.body.([]byte)
	_, err := jpeg.Decode(bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("response is not a valid JPEG image %w", err)
	}

	return nil
}

func (pt *previewerTest) iImageResizeUpCheckResponse(ctx context.Context, expectedStatus int, expectedWidth,
	expectedHeight int,
) error {
	actualResp, ok := ctx.Value(godogsResponseCtxKey{}).(*response)
	if !ok {
		return errors.New("there are no godogs available")
	}
	if err := checkStatus(expectedStatus, actualResp.status); err != nil {
		return err
	}

	body := actualResp.body.([]byte)
	img, err := jpeg.Decode(bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("response is not a valid JPEG image %w", err)
	}
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	if width != expectedWidth {
		return fmt.Errorf("expectedWidth %d does not equal actual width %v", expectedWidth, width)
	}
	if height != expectedHeight {
		return fmt.Errorf("expectedHeight %d does not equal actual height %v", expectedHeight, height)
	}

	return nil
}

func (pt *previewerTest) iImageFoundLocalCheckResponse(ctx context.Context, expectedStatus int, header string) error {
	actualResp, ok := ctx.Value(godogsResponseCtxKey{}).(*response)
	if !ok {
		return errors.New("there are no godogs available")
	}
	if err := checkStatus(expectedStatus, actualResp.status); err != nil {
		return err
	}

	body := actualResp.body.([]byte)
	_, err := jpeg.Decode(bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("response is not a valid JPEG image %w", err)
	}
	hit := actualResp.headers.Get("X-Previewer-Cache-Hit")
	if hit != header {
		return fmt.Errorf("cache hit header must be true, received %v", hit)
	}

	return nil
}

func (pt *previewerTest) iImageHeaderAbsenceCheckResponse(ctx context.Context, expectedStatus int,
	subStr string,
) error {
	actualResp, ok := ctx.Value(godogsResponseCtxKey{}).(*response)
	if !ok {
		return errors.New("there are no godogs available")
	}
	if err := checkStatus(expectedStatus, actualResp.status); err != nil {
		return err
	}

	body := string(actualResp.body.([]byte))

	if !strings.Contains(body, subStr) {
		return fmt.Errorf("response %v does not contain %v", body, subStr)
	}

	return nil
}

func checkStatus(expectedStatus, actualStatus int) error {
	if expectedStatus != actualStatus {
		return fmt.Errorf("expected response code to be: %d, but actual is: %d", expectedStatus, actualStatus)
	}
	return nil
}
