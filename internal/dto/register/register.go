package register

type RegisterResponse struct {
	Status                          string `json:"status"`
	ThreadAmount                    int    `json:"thread_amount"`
	ClientSendUpdateStatusInSeconds int    `json:"CLIENT_SEND_UPDATE_STATUS_IN_SECONDS"`
}
