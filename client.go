package fcm

import (
	"context"
	"fmt"

	"github.com/valyala/fasthttp"
	"golang.org/x/oauth2"
)

// SimpleClient abstracts the interaction between the application server and the
// FCM server via HTTP protocol.
// It uses Service Account Private Key for authentication.
//
// If the `HTTP` field is nil, a zeroed http.SimpleClient will be allocated and used
// to send messages.
type Client interface {
	// Send sends a message to the FCM server without retrying in case of service
	// unavailability. A non-nil error is returned if a non-recoverable error
	// occurs (i.e. if the sendResponse status code is not between 200 and 299).
	Send(ctx context.Context, msg *Message) error
}

var _ Client = (*SimpleClient)(nil)

type SimpleClient struct {
	client      FastHTTPDoer
	url         urlConfig
	tokenSource oauth2.TokenSource
	sendPath    []byte
}

// NewClient creates new Firebase Cloud Messaging SimpleClient based on API key and
// with default endpoint and http client.
func NewClient(serviceAccountJSONData []byte, opts ...Option) *SimpleClient {
	defaultOpts := []Option{
		WithEndpoint(DefaultEndpoint),
		WithCredentialsData(serviceAccountJSONData),
		WithHTTPClient(DefaultHTTPAdapter),
	}

	opts = append(defaultOpts, opts...)
	return newClient(opts...)
}

func newClient(opts ...Option) *SimpleClient {
	c := SimpleClient{
		tokenSource: &NoopTokenSource{},
	}

	if err := applyOptions(&c, opts...); err != nil {
		panic(err)
	}

	return &c
}

// Send implementation of Client interface.
// Docs for the reference: https://firebase.google.com/docs/reference/fcm/rest/v1/projects.messages/send
func (c *SimpleClient) Send(ctx context.Context, msg *Message) error {
	if err := msg.Validate(); err != nil {
		return fmt.Errorf("invalid message: %w", err)
	}

	sendReq := sendRequest{
		Message: msg,
	}

	body, err := sendReq.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	authHeaderValue, err := c.authHeaderValue()
	if err != nil {
		return err
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.Header.SetMethod(fasthttp.MethodPost)
	uri := req.URI()
	uri.SetSchemeBytes(c.url.Scheme)
	uri.SetHostBytes(c.url.Host)
	uri.SetPathBytes(c.sendPath)
	req.Header.SetBytesKV(contentTypeHeader, contentTypeHeaderV)
	req.Header.SetBytesKV(authorizationHeader, authHeaderValue)
	req.SetBody(body)

	if err := c.client.Do(ctx, req, resp); err != nil {
		return fmt.Errorf("failed to perform request: %w", err)
	}

	return handleResponse(resp.StatusCode(), resp.Body())
}

func (c *SimpleClient) authHeaderValue() ([]byte, error) {
	// TODO: consider to regenerate this value only when token is expired
	//  e.g. cache and reuse if not expired
	token, err := c.tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to grab oauth2 token: %w", err)
	}

	headerValue := token.Type() + " " + token.AccessToken
	return []byte(headerValue), nil
}

// Handle sendResponse
func handleResponse(statusCode int, respBody []byte) error {
	// Success sendResponse
	if statusCode >= 200 && statusCode <= 299 {
		return nil
	}

	var resp sendResponse
	if err := resp.UnmarshalJSON(respBody); err != nil {
		return fmt.Errorf("unmarshal sendResponse with status code %d: %w", statusCode, err)
	}

	// In case if error sendResponse is not like we expect.
	if resp.Error == nil {
		return fmt.Errorf("empty error in sendResponse for status code %d: %s", statusCode, string(respBody))
	}

	// Extract errorCode of google.firebase.fcm type.
	// Details could be different types of structs
	// and some of the errorDetails elements could not have ErrorCode.
	var errCode errorCode
	if len(resp.Error.Details) > 0 {
		for _, detail := range resp.Error.Details {
			if detail.ErrorCode != "" {
				errCode = detail.ErrorCode
				break
			}
		}
	}

	switch errCode {
	case errorCodeUnregistered:
		return ErrUnregistered
	default:
		return fmt.Errorf("unsuccessful sendResponse with status code: %d: %s", statusCode, string(respBody))
	}
}
