package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/asiainfoLDP/datahub_commons/common"
)

//======================================================
// 
//======================================================

const (
	StringParamType_General        = 0
	StringParamType_UrlWord        = 1
	StringParamType_UnicodeUrlWord = 2
	StringParamType_Email          = 3
)

//======================================================
//
//======================================================

func MustBoolParam(params httprouter.Params, paramName string) (bool, error) {
	bool_str := params.ByName(paramName)
	if bool_str == "" {
		return false, fmt.Errorf("%s can't be blank", paramName)
	}

	b, err := strconv.ParseBool(bool_str)
	if err != nil {
		return false, fmt.Errorf("%s=%s", paramName, bool_str)
	}

	return b, nil
}

func MustBoolParamInMap(m map[string]interface{}, paramName string) (bool, error) {
	v, ok := m [paramName]
	if ok {
		b, ok := v.(bool)
		if ok {
			return b, nil
		}
		
		return false, fmt.Errorf("param %s is not bool", paramName)
	}
	
	return false, fmt.Errorf("param %s is not found", paramName)
}

func MustBoolParamInQuery(r *http.Request, paramName string) (bool, error) {
	bool_str := r.Form.Get(paramName)
	if bool_str == "" {
		return false, fmt.Errorf("%s can't be blank", paramName)
	}

	b, err := strconv.ParseBool(bool_str)
	if err != nil {
		return false, fmt.Errorf("%s=%s", paramName, bool_str)
	}

	return b, nil
}

func optionalBoolParamInQuery(r *http.Request, paramName string, defaultValue bool) bool {
	bool_str := r.Form.Get(paramName)
	if bool_str == "" {
		return defaultValue
	}

	b, err := strconv.ParseBool(bool_str)
	if err != nil {
		return defaultValue
	}

	return b
}

func _optionalIntParam(intStr string, defaultInt int64) int64 {
	if intStr == "" {
		return defaultInt
	}

	i, err := strconv.ParseInt(intStr, 10, 64)
	if err != nil {
		return defaultInt
	} else {
		return i
	}
}

func optionalIntParamInQuery(r *http.Request, paramName string, defaultInt int64) int64 {
	return _optionalIntParam(r.Form.Get(paramName), defaultInt)
}

func _mustIntParam(paramName string, int_str string) (int64, error) {
	if int_str == "" {
		return 0, fmt.Errorf("%s can't be blank", paramName)
	}

	i, err := strconv.ParseInt(int_str, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s=%s", paramName, int_str)
	}

	return i, nil
}

func MustIntParamInQuery(r *http.Request, paramName string) (int64, error) {
	return _mustIntParam(paramName, r.Form.Get(paramName))
}

func MustIntParamInPath(params httprouter.Params, paramName string) (int64, error) {
	return _mustIntParam(paramName, params.ByName(paramName))
}

func MustIntParamInMap(m map[string]interface{}, paramName string) (int64, error) {
	v, ok := m[paramName]
	if ok {
		i, ok := v.(float64)
		if ok {
			return int64(i), nil
		}

		return 0, fmt.Errorf("param %s is not int", paramName)
	}

	return 0, fmt.Errorf("param %s is not found", paramName)
}

func optionalIntParamInMap(m map[string]interface{}, paramName string, defaultValue int64) int64 {
	v, ok := m[paramName]
	if ok {
		i, ok := v.(float64)
		if ok {
			return int64(i)
		}
	}

	return defaultValue
}

func MustFloatParam(params httprouter.Params, paramName string) (float64, error) {
	float_str := params.ByName(paramName)
	if float_str == "" {
		return 0.0, fmt.Errorf("%s can't be blank", paramName)
	}

	f, err := strconv.ParseFloat(float_str, 64)
	if err != nil {
		return 0.0, fmt.Errorf("%s=%s", paramName, float_str)
	}

	return f, nil
}

func _mustStringParam(paramName string, str string, paramType int) (string, error) {
	if str == "" {
		return "", fmt.Errorf("param: %s can't be blank", paramName)
	}

	if paramType == StringParamType_UrlWord {
		str2, ok := common.ValidateUrlWord(str)
		if !ok {
			return "", fmt.Errorf("param: %s=%s", paramName, str)
		}
		str = str2
	} else if paramType == StringParamType_UnicodeUrlWord {
		str2, ok := common.ValidateUnicodeUrlWord(str)
		if !ok {
			return "", fmt.Errorf("param: %s=%s", paramName, str)
		}
		str = str2
	} else if paramType == StringParamType_Email {
		str2, ok := common.ValidateEmail(str)
		if !ok {
			return "", fmt.Errorf("param: %s=%s", paramName, str)
		}
		str = str2
	} else { // if paramType == StringParamType_General
		str2, ok := common.ValidateGeneralWord(str)
		if !ok {
			return "", fmt.Errorf("param: %s=%s", paramName, str)
		}
		str = str2
	}

	return str, nil
}

func MustStringParamInPath(params httprouter.Params, paramName string, paramType int) (string, error) {
	return _mustStringParam(paramName, params.ByName(paramName), paramType)
}

func MustStringParamInQuery(r *http.Request, paramName string, paramType int) (string, error) {
	return _mustStringParam(paramName, r.Form.Get(paramName), paramType)
}

func MustStringParamInMap(m map[string]interface{}, paramName string, paramType int) (string, error) {
	v, ok := m[paramName]
	if ok {
		str, ok := v.(string)
		if ok {
			return _mustStringParam(paramName, str, paramType)
		}

		return "", fmt.Errorf("param %s is not string", paramName)
	}

	return "", fmt.Errorf("param %s is not found", paramName)
}

func optionalTimeParamInQuery(r *http.Request, paramName string, timeLayout string, defaultTime time.Time) time.Time {
	str := r.Form.Get(paramName)
	if str == "" {
		return defaultTime
	}

	t, err := time.Parse(timeLayout, str)
	if err != nil {
		return defaultTime
	}

	return t
}
