package web_models

const (
	// Requests
	RequestGetPlatforms     = "GET_PLATFORMS"
	RequestGetUser          = "GET_USER"
	RequestRegisterPlatform = "REGISTER_PLATFORM"

	// Responses
	ResponsePlatforms          = "PLATFORMS"
	ResponseUser               = "USER"
	ResponsePlatformRegistered = "PLATFORM_REGISTERED"
)

type WebSocketRequest struct {
	Type       string
	CallbackId int64
	Data       interface{}
}

type WebSocketResponse struct {
	WebSocketRequest
	ErrorCode int
	Error     string
	Refresh   bool
}

func (wr *WebSocketResponse) SetError(code int, msg string) *WebSocketResponse {
	wr.ErrorCode = code
	wr.Error = msg

	return wr
}

func (wr *WebSocketResponse) RefreshContent() *WebSocketResponse {
	wr.Refresh = true

	return wr
}

func (wr *WebSocketResponse) SetCallback(reg *WebSocketRequest) *WebSocketResponse {
	wr.CallbackId = reg.CallbackId

	return wr
}

func (wr *WebSocketResponse) SetData(dataType string, i interface{}) *WebSocketResponse {
	wr.Type = dataType
	wr.Data = i

	return wr
}
