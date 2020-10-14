package external

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"transport/lib/thttp/httpclient"
	"transport/lib/thttp/httpclient/body"
	"transport/lib/thttp/httpclient/hook"
	"transport/lib/utils/logger"

	"telegram-bot/app/models"
	"telegram-bot/app/schema"
	"telegram-bot/pkg/utils"
)

type External interface {
	CallAPI(ctx context.Context, message *models.InMessage) (*schema.ExternalResponse, error)
}

type external struct {
	httpClient httpclient.Client
}

func New(httpClient httpclient.Client) External {
	return &external{
		httpClient: httpClient,
	}
}

func (e *external) CallAPI(ctx context.Context, message *models.InMessage) (*schema.ExternalResponse, error) {
	res := schema.ExternalResponse{
		HandleTime: time.Now().UTC().Format(time.RFC3339Nano),
	}
	function := e.getFunction(message.RoutingKey.APIMethod)
	url := message.RoutingKey.APIUrl + e.getUrlSuffix(message)
	authentication, err := e.setAuthentication(&message.RoutingKey)
	if err != nil {
		return nil, err
	}

	status, err := function(
		ctx,
		url,
		httpclient.WithHookFn(authentication),
		httpclient.WithBodyProvider(body.NewJson(message.Payload)),
		httpclient.WithHookFn(hook.UnmarshalResponse(&res)),
	)

	if err != nil {
		logger.Errorf("Failed to send request: %s, error: %s", url, err)
		return e.parseError(err), utils.EMSCallAPIError.New()
	}
	//service error
	if status != http.StatusOK {
		logger.Error(fmt.Sprintf("Failed to call API: %s, status: ", url), res.Status)
		return &res, utils.EMSCallAPIError.New()
	}

	// check by pass status
	if utils.InList(message.RoutingKey.ByPassStatus, res.Status) != -1 {
		logger.Info(fmt.Sprintf("Call API successfully: %s, status: ", url), res.Status)
		return &res, nil
	}

	// check response status
	if res.Status != "OK" && res.Status != 1 {
		logger.Error(fmt.Sprintf("Failed to call API: %s, status: ", url), res.Status)
		return &res, utils.EMSCallAPIError.New()
	}

	logger.Info(fmt.Sprintf("Call API successfully: %s, status: ", url), res.Status)
	return &res, nil
}

func (e *external) getFunction(method string) func(ctx context.Context, url string, opts ...httpclient.RequestOption) (int, error) {
	switch method {
	case "GET":
		return e.httpClient.Get
	case "POST":
		return e.httpClient.Post
	case "PUT":
		return e.httpClient.Put
	case "DELETE":
		return e.httpClient.Delete
	default:
		return nil
	}
}

func (e *external) getUrlSuffix(message *models.InMessage) string {
	suffix := ""
	if message.ExternalID != "" {
		suffix += "/" + message.ExternalID
	}
	if message.Query != "" {
		suffix += "?" + message.Query
	}

	return suffix
}

func (e *external) setAuthentication(routing *models.RoutingKey) (httpclient.HookFn, error) {
	switch routing.AuthType {
	case models.BasicAuthentication:
		auth := strings.Split(routing.AuthKey, ":")
		if len(auth) != 2 {
			logger.Errorf("Auth key is invalid")
			return nil, utils.EMSDataInvalidError.Newm("auth key is invalid")
		}
		return hook.BasicAuth(auth[0], auth[1]), nil
	case models.APIKeyAuthentication:
		if routing.AuthKey == "" {
			return nil, utils.EMSDataInvalidError.Newm("auth key is invalid")
		}
		return e.apiKeyAuth(routing.AuthKey), nil
	default:
		return nil, nil
	}
}

func (e *external) apiKeyAuth(apiKey string) httpclient.HookFn {
	return func(ctx context.Context, reqChain httpclient.Chain) (*http.Response, error) {
		req := reqChain.GetRequest(ctx)
		req.Header.Add("x-api-key", apiKey)
		return reqChain.Proceed(ctx, req)
	}
}

func (e *external) parseError(err error) *schema.ExternalResponse {
	res := schema.ExternalResponse{
		Status:     int(utils.EMSCallAPIError),
		Message:    err.Error(),
		HandleTime: time.Now().UTC().Format(time.RFC3339Nano),
	}

	return &res
}
