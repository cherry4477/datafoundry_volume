package main

type APIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status,omitempty"`
	//Data    interface{} `json:"data,omitempty"`
}
