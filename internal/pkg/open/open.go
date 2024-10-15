package open

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/UniqueStudio/UniqueSSOBackend/internal/tracer"
	"github.com/xylonx/zapx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
	"net/http"
)

type OpenClient struct {
	*http.Client
	BaseAddr string
	Token    string
}

type PushSMSRequest struct {
	Phone      string   `json:"phone_number"`
	TemplateId uint     `json:"template_id"`
	Params     []string `json:"template_param_set"`
	SignID     *uint    `json:"sign_id,omitempty"`
}

func NewOpenClient(baseAddr string, token string) *OpenClient {
	return &OpenClient{
		Client:   http.DefaultClient,
		BaseAddr: baseAddr,
		Token:    token,
	}
}

func (client *OpenClient) PushSMS(ctx context.Context, smsRequest *PushSMSRequest) (*http.Response, error) {
	apmCtx, span := tracer.Tracer.Start(ctx, "PushSMS")
	defer span.End()

	zlog := zapx.WithContext(apmCtx)
	reqBytes, err := json.Marshal(smsRequest)

	//type oldPushSMSRequest struct {
	//	Phone      string   `json:"phone_number"`
	//	TemplateId string   `json:"template_id"`
	//	Params     []string `json:"template_param_set"`
	//}
	//oldReq := oldPushSMSRequest{
	//	Phone:      smsRequest.Phone,
	//	TemplateId: strconv.Itoa(int(smsRequest.TemplateId)),
	//	Params:     smsRequest.Params,
	//}
	//reqBytes, err := json.Marshal(oldReq)

	if err != nil {
		zlog.With(zap.Error(err)).Error("json marshal error")
		return nil, err
	}
	path := "/sms/send_single"
	req, err := http.NewRequest(http.MethodPost, client.BaseAddr+path, bytes.NewBuffer(reqBytes))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		zlog.With(zap.Error(err)).Error("new sms request error")
		return nil, err
	}

	req.Header.Set("AccessKey", client.Token)

	span.SetAttributes(attribute.String("PushSMSReq", fmt.Sprintf("%v", smsRequest)))
	resp, err := client.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		zlog.With(zap.Error(err)).Error("send sms request to open-platform error")
		return nil, err
	}

	return resp, nil
}
