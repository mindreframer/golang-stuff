/*http response formats*/
package hamster

type NewDeveloperResponse struct {
	ObjectId    string `json:"object_id"`
	AccessToken string `json:"access_token"`
}

type LoginResponse struct {
	ObjectId    string `json:"object_id"`
	AccessToken string `json:"access_token"`
	Status      string `json:"status"`
}

type VerifyLogin struct {
	Status string `json:"status"`
}

type LogoutResponse struct {
	Status string `json:"status"`
}

type QueryDevResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type DeleteResponse struct {
	Status string `json:"status"`
}

type OKResponse struct {
	Status string `json:"status"`
}

type AppResponse struct {
	ApiToken  string `json:"apitoken"`
	ApiSecret string `json:"apisecret"`
	Name      string `json:"name"`
	OS        string `json:"os"`
}

type AllAppResponse struct {
	Responses []AppResponse `json:"responses"`
}

type SaveFileResponse struct {
	FileId   string `json:"file_id"`
	FileName string `json:"file_name"`
}
