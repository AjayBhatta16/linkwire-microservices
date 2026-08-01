package myfunction

type Request struct {
	ResetRequestID string `json:"resetRequestId"`
	OldPassword    string `json:"oldPassword"`
	NewPassword    string `json:"newPassword"`
}