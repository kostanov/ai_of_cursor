package models

// User — запись пользователя в БД.
type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// AddUserRequest — тело POST /adduser.
type AddUserRequest struct {
	Name string `json:"name"`
}

// AddUserResponse — ответ POST /adduser.
type AddUserResponse struct {
	Status string `json:"status"`
	ID     int64  `json:"id"`
	Name   string `json:"name"`
}

// ActivateResponse — ответ GET|POST /activate/{uid}.
type ActivateResponse struct {
	Status string `json:"status"`
	Active []int  `json:"active"`
}

// StatusResponse — простой статус (например /slow).
type StatusResponse struct {
	Status string `json:"status"`
}

// ErrorBody — JSON с полем error.
type ErrorBody struct {
	Error string `json:"error"`
}

// WrongErrorBody — JSON для /wrong.
type WrongErrorBody struct {
	Msg    string `json:"msg"`
	Detail string `json:"detail"`
}
