package sfsdk

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type env int

const (
	Test env = iota
	Prod
	HkProd
)

type Client struct {
	Env          env
	CustomerCode string
	CheckWord    string
	HttpClient   http.Client
}

type ApiResult struct {
	ApiResultCode string `json:"apiResultCode"`
	ApiErrorMsg   string `json:"apiErrorMsg"`
	ApiResponseID string `json:"apiResponseID"`
	ApiResultData string `json:"apiResultData"`
}

func NewClient(customerCode string, checkWord string, env env, httpclient ...http.Client) Client {
	var client http.Client

	if httpclient == nil {
		client = http.Client{
			Timeout: time.Second * 10,
		}
	} else {
		client = httpclient[0]
	}

	return Client{
		Env:          env,
		CustomerCode: customerCode,
		CheckWord:    checkWord,
		HttpClient:   client,
	}
}

func (c Client) Execute(serviceCode string, msgData string) (*ApiResult, error) {
	timestampStr := strconv.FormatInt(time.Now().UnixMilli(), 10)
	urlEncodedStr := url.QueryEscape(msgData + timestampStr + c.CheckWord)

	digest := md5.New()
	digest.Write([]byte(urlEncodedStr))
	md5Str := base64.StdEncoding.EncodeToString(digest.Sum(nil))

	var values = make(url.Values)
	values.Add("requestID", uuid.NewString())
	values.Add("partnerID", c.CustomerCode)
	values.Add("serviceCode", serviceCode)
	values.Add("timestamp", timestampStr)
	values.Add("msgDigest", md5Str)
	values.Add("msgData", msgData)

	resp, err := c.HttpClient.PostForm(c.GetRequestUrl(), values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respObj := &ApiResult{}

	err = json.Unmarshal(respBytes, respObj)
	if err != nil {
		return nil, err
	}
	return respObj, nil
}

func (c Client) GetRequestUrl() string {
	switch c.Env {
	case Prod:
		return "https://bspgw.sf-express.com/std/service"
	case HkProd:
		return "https://sfapi-hk.sf-express.com/std/service"
	default:
		return "https://sfapi-sbox.sf-express.com/std/service"
	}
}

func (r ApiResult) IsSuccess() bool {
	return r.ApiResultCode == "A1000"
}

func (r ApiResult) String() string {
	return fmt.Sprintf("%#v", r)
}

func (r ApiResult) Json() string {
	bytes, _ := json.Marshal(r)
	return string(bytes)
}
