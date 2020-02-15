package gorsk

// Success response
// swagger:response ok
type swaggOKResp struct{}

// Error response
// swagger:response err
type swaggErrResp struct{}

// Error response with message
// swagger:response errMsg
type swaggErrMsgResp struct {
	Message string `json:"message"`
}
