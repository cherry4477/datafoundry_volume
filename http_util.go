package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)



const (
	ErrCodeOK           = 1200


	ErrCodeUnknownError = 1600
)

func retHttpCodef(code, bodyCode int, w http.ResponseWriter, format string, a ...interface{}) {

	w.WriteHeader(code)
	msg := fmt.Sprintf(`{"code":%d,"msg":"%s"}`, bodyCode, fmt.Sprintf(format, a...))

	fmt.Fprintf(w, msg)
	return
}

func retHttpCode(code int, bodyCode int, w http.ResponseWriter, a ...interface{}) {
	w.WriteHeader(code)
	msg := fmt.Sprintf(`{"code":%d,"msg":"%s"}`, bodyCode, fmt.Sprint(a...))

	fmt.Fprintf(w, msg)
	return
}

func retHttpCodeJson(code int, bodyCode int, w http.ResponseWriter, a ...interface{}) {
	w.WriteHeader(code)
	msg := fmt.Sprintf(`{"code":%d,"msg":%s}`, bodyCode, fmt.Sprint(a...))

	fmt.Fprintf(w, msg)
	return
}

func RespError(w http.ResponseWriter, err error, httpCode int) {
	resp := genRespJson(httpCode, err)

	if body, err := json.MarshalIndent(resp, "", "  "); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpCode)
		w.Write(body)
	}
}

func RespOK(w http.ResponseWriter, data interface{}) {
	if data == nil {
		data = genRespJson(http.StatusOK, nil)
	}

	if body, err := json.MarshalIndent(data, "", "  "); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}
}

func RespAccepted(w http.ResponseWriter, data interface{}) {
	if data == nil {
		data = genRespJson(http.StatusAccepted, nil)
	}

	if body, err := json.MarshalIndent(data, "", "  "); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		w.Write(body)
	}
}

func genRespJson(httpCode int, err error) *APIResponse {
	resp := new(APIResponse)
	var msgCode int
	var message string

	if err == nil {
		msgCode = ErrCodeOK
		message = "OK"
	} else {
		msgCode = ErrCodeUnknownError
		message = err.Error()
	}

	resp.Code = msgCode
	resp.Message = message
	resp.Status = http.StatusText(httpCode)
	return resp
}
