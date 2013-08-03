/*http request formats*/
package hamster

type LogoutRequest struct {
	Email string `json:"email"`
}

type UpdateRequest struct {
	Name string `json:"name"`
}

type UpdateAppRequest struct {
	Name string `json:"name"`
	OS   string `json:"os"`
}
