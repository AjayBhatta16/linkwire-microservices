package myfunction

type Request struct {
	ResetRequestID string `json:"resetRequestId"`
	NewPassword    string `json:"newPassword"`
}