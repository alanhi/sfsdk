package sfsdk

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type env int

const (
	Test env = iota
	Prod
	HkProd
)

type Client struct {
	env          env
	customerCode string
	checkWord    string
	httpClient   *http.Client
}

func NewClient(customerCode string, checkWord string, env env, httpclient ...*http.Client) *Client {
	var client *http.Client

	if httpclient == nil {
		client = &http.Client{
			Timeout: time.Second * 10,
		}
	} else {
		client = httpclient[0]
	}

	return &Client{
		env:          env,
		customerCode: customerCode,
		checkWord:    checkWord,
		httpClient:   client,
	}
}

func (c *Client) Execute(serviceCode string, msgData string) (*ApiResult, error) {
	timestampStr := strconv.FormatInt(time.Now().UnixMilli(), 10)
	urlEncodedStr := url.QueryEscape(msgData + timestampStr + c.checkWord)

	digest := md5.New()
	digest.Write([]byte(urlEncodedStr))
	md5Str := base64.StdEncoding.EncodeToString(digest.Sum(nil))

	var values = make(url.Values)
	values.Add("requestID", uuid.NewString())
	values.Add("partnerID", c.customerCode)
	values.Add("serviceCode", serviceCode)
	values.Add("timestamp", timestampStr)
	values.Add("msgDigest", md5Str)
	values.Add("msgData", msgData)

	resp, err := c.httpClient.PostForm(c.GetRequestUrl(), values)
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

func (c *Client) GetRequestUrl() string {
	switch c.env {
	case Prod:
		return "https://bspgw.sf-express.com/std/service"
	case HkProd:
		return "https://sfapi-hk.sf-express.com/std/service"
	default:
		return "https://sfapi-sbox.sf-express.com/std/service"
	}
}

func (c *Client) GetEnv() env {
	return c.env
}
