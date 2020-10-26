package external

import (
	"context"
	"fmt"
	"net/http"

	"transport/lib/thttp/httpclient"
	"transport/lib/thttp/httpclient/body"
	"transport/lib/thttp/httpclient/hook"
	"transport/lib/utils/logger"

	"telegram-bot/app/models"
	"telegram-bot/app/schema"
	"telegram-bot/pkg/utils"
)

type External interface {
	Call(ctx context.Context, message *models.Message, action string) (*schema.ExternalResponse, error)
}

type external struct {
	httpClient httpclient.Client
}

func New(httpClient httpclient.Client) External {
	return &external{
		httpClient: httpClient,
	}
}

func (e *external) getUrl(ctx context.Context, msg *models.Message, action string) string {
	switch action {
	case "confirm":
		return msg.Action.UrlConfirm
	case "reject":
		return msg.Action.UrlReject
	case "cancel":
		return msg.Action.UrlCancel
	default:
		return ""
	}
}

func (e *external) Call(ctx context.Context, msg *models.Message, action string) (*schema.ExternalResponse, error) {
	res := schema.ExternalResponse{}

	url := e.getUrl(ctx, msg, action)
	if url == "" {
		logger.Errorf("Cannot get url for action %s", action)
	}

	status, err := e.httpClient.Put(
		ctx,
		url,
		httpclient.WithBodyProvider(body.NewJson(msg.Data)),
		httpclient.WithHookFn(hook.UnmarshalResponse(&res)),
	)

	if err != nil {
		logger.Errorf("Failed to send request: %s, error: %s", url, err)
		return nil, utils.BOTCallAPIError.New()
	}
	//service error
	if status != http.StatusOK {
		logger.Error(fmt.Sprintf("Failed to call API: %s, status: ", url), res.Status)
		return &res, utils.BOTCallAPIError.New()
	}

	// check response status
	if res.Status != "OK" && res.Status != 1 {
		logger.Error(fmt.Sprintf("Failed to call API: %s, status: ", url), res.Status)
		return &res, utils.BOTCallAPIError.New()
	}

	logger.Info(fmt.Sprintf("Call API successfully: %s, status: ", url), res.Status)
	return &res, nil
}
