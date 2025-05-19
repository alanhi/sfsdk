package sfsdk

import (
	"encoding/json"
	"fmt"
)

type ApiResult struct {
	ApiResultCode string `json:"apiResultCode"`
	ApiErrorMsg   string `json:"apiErrorMsg"`
	ApiResponseID string `json:"apiResponseID"`
	ApiResultData string `json:"apiResultData"`
}

func (r *ApiResult) IsSuccess() bool {
	return r.ApiResultCode == "A1000"
}

func (r *ApiResult) String() string {
	return fmt.Sprintf("%#v", r)
}

func (r *ApiResult) Json() string {
	bytes, _ := json.Marshal(r)
	return string(bytes)
}
