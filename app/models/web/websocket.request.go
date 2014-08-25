package web_models

const (
	// Requests
	RequestGetPlatforms          = "GET_PLATFORMS"
	RequestGetUser               = "GET_USER"
	RequestRegisterPlatform      = "REGISTER_PLATFORM"
	RequestRegisterFullSite      = "REGISTER_FULL_SITE"
	RequestGetSiteTemplates      = "GET_SITE_TEMPLATES"
	RequestGetServerTemplates    = "GET_SERVER_TEMPLATES"
	RequestGetServerCertificates = "GET_SERVER_CERTIFICATES"

	// Responses
	ResponsePlatforms          = "PLATFORMS"
	ResponseUser               = "USER"
	ResponsePlatformRegistered = "PLATFORM_REGISTERED"
	ResponseStatusMessage      = "STATUS_MESSAGE"
	ResponseSiteTemplates      = "SITE_TEMPLATES"
	ResponseServerTemplates    = "SERVER_TEMPLATES"
	ResponseServerCertificates = "SERVER_CERTIFICATES"
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
