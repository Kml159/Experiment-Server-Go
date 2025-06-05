package register

type RegisterResponse struct {
    Status       string `json:"status"`
    ThreadAmount int    `json:"thread_amount"`
}
