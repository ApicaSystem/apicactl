package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/logiqai/logiqctl/defines"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

var (
	client *ApiClient
)

type TokenType int

const (
	TokenType_APIKEY = iota
	TokenType_BEARER = iota
)

type Api interface {
	MakeApiCall(method string, url string, payload *bytes.Buffer) (*http.Response, error)
}

type ApiClient struct {
	AuthToken string
	AuthType  TokenType
	Url       string
	client    *http.Client
	TraceFlag bool
}

func (c *ApiClient) initHttpClient() error {
	httpTransport, ok := http.DefaultTransport.(*http.Transport)
	if ok {
		httpTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	if c.AuthToken != "" {
		c.client = &http.Client{}
	} else {
		httpClient, err := loginWithEmailAndPassword(c.Url)
		if err != nil {
			return err
		}
		c.client = httpClient
	}
	return nil
}

func (c *ApiClient) MakeApiCall(method string, url string, payload *bytes.Buffer) (*http.Response, error) {
	var req *http.Request
	var err error
	var res *http.Response

	if method == http.MethodPost {
		req, err = c.createHttpRequest(method, url, payload)
	} else {
		req, err = c.createHttpRequest(method, url, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("error sending request")
	}

	if c.TraceFlag == true {
		req = AddNetTrace(req)
	}

	res, err = c.client.Do(req)
	if err != nil {
		return nil, err
	}
	contentType := strings.Split(res.Header.Get("Content-Type"), ";")[0]
	if contentType != "application/json" {
		return nil, fmt.Errorf("Unexpected Response. '%s'  is not a logiq endpoint. Please try chanding the endpoint using config command", viper.GetString(KeyCluster))
	}

	return res, nil
}

func (c *ApiClient) GetResponseString(resp *http.Response) ([]byte, error) {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("error:%s", err.Error())
	}
	if bodyBytes == nil {
		return []byte{}, fmt.Errorf("error: Response is Empty")
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return []byte{}, fmt.Errorf("error: %s", string(bodyBytes))
	}
	return bodyBytes, nil
}

func (c *ApiClient) createHttpRequest(method string, uri string, payload *bytes.Buffer) (*http.Request, error) {
	var req *http.Request
	var err error
	uri = fmt.Sprintf("%s/%s", c.Url, uri)
	if method == http.MethodGet || method == http.MethodDelete {
		req, err = http.NewRequest(method, uri, nil)
	} else {
		req, err = http.NewRequest(method, uri, payload)
	}
	if err != nil {
		return nil, err
	}
	tokenType := "Key"
	if c.AuthType == TokenType_BEARER {
		tokenType = "Bearer"
	}
	if api_key := c.AuthToken; api_key != "" {
		req.Header.Add("Authorization", fmt.Sprintf("%s %s", tokenType, api_key))
	}
	req.Header.Add("Content-Type", "application/json")
	return req, err
}

func loginWithEmailAndPassword(endpoint string) (*http.Client, error) {
	user := GetUIUser()
	password := GetUIPass()

	if user != "" && password != "" {
		cookieJar, _ := cookiejar.New(nil)

		c := &http.Client{
			Jar: cookieJar,
			//Transport: &http.Transport{
			//	MaxIdleConnsPerHost: 20,
			//},
			//Timeout: 1000 * time.Second,
		}

		loginUrl := endpoint + "/" + GetUrlForResource(defines.ResourceLogin)
		u, _ := url.Parse(loginUrl)
		q, _ := url.ParseQuery(u.RawQuery)
		q.Add("remember", "on")
		q.Add("email", user)
		q.Add("password", password)
		u.RawQuery = q.Encode()

		if resp, err := c.Post(u.String(), "application/x-www-form-urlencoded", bytes.NewReader(([]byte)(q.Encode()))); err != nil {
			return nil, err
		} else {
			defer resp.Body.Close()
			b, _ := ioutil.ReadAll(resp.Body)
			if match, _ := regexp.Match("Wrong email or password", b); match {
				return nil, fmt.Errorf("Wrong email or password")
			}
		}
		return c, nil
	} else {
		return nil, fmt.Errorf("api token or ui credentials must be set. See \"logiqctl config help\" for more details")
	}
}

func InitApiClient(authToken string, authType TokenType, endpoint string, trace bool) error {
	c := ApiClient{
		AuthToken: authToken,
		AuthType:  authType,
		Url:       endpoint,
		TraceFlag: trace,
	}
	if !strings.Contains(endpoint, "://") {
		protocolType := getProtocol(endpoint)
		if protocolType == defines.UriHttp {
			c.Url = "http://" + c.Url
		} else if protocolType == defines.UriHttps {
			c.Url = "https://" + c.Url
		} else {
			return fmt.Errorf("Unsupported protocol in cluster endpoint")
		}
	}
	err := c.initHttpClient()
	if err != nil {
		return err
	}
	client = &c
	return nil
}

func GetApiClient() *ApiClient {
	return client
}

func GetEndpointWithProtocol(endpoint string) string {
	protocol := getProtocol(endpoint)
	if protocol == defines.UriHttp {
		return "http://" + endpoint
	} else if protocol == defines.UriHttps {
		return "https://" + endpoint
	}
	return endpoint
}
