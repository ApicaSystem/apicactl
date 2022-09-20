package ui

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"

	"github.com/logiqai/logiqctl/utils"
	"github.com/spf13/viper"
)

var (
	client *http.Client = nil
)

func getHttpClient() *http.Client {
	if client != nil {
		return client
	}

	api_key := viper.GetString(utils.AuthToken)
	httpTransport, ok := http.DefaultTransport.(*http.Transport)
	if ok {
		httpTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	if api_key != "" {
		client = &http.Client{}
	} else {
		user := utils.GetUIUser()
		password := utils.GetUIPass()

		if user != "" && password != "" {
			cookieJar, _ := cookiejar.New(nil)

			client = &http.Client{
				Jar: cookieJar,
				//Transport: &http.Transport{
				//	MaxIdleConnsPerHost: 20,
				//},
				//Timeout: 1000 * time.Second,
			}

			loginUrl := GetUrlForResource(ResourceLogin)
			u, _ := url.Parse(loginUrl)
			q, _ := url.ParseQuery(u.RawQuery)
			q.Add("remember", "on")
			q.Add("email", user)
			q.Add("password", password)
			u.RawQuery = q.Encode()

			if resp, err := client.Post(u.String(), "application/x-www-form-urlencoded", bytes.NewReader(([]byte)(q.Encode()))); err != nil {
				fmt.Println("Error login with provided credentials, Error:", err.Error())
				os.Exit(-1)
			} else {
				defer resp.Body.Close()
				b, _ := ioutil.ReadAll(resp.Body)
				if match, _ := regexp.Match("Wrong email or password", b); match {
					fmt.Println("Error credentials")
					os.Exit(-1)
				}
			}
		} else {
			fmt.Println("api token or ui credentials must be set. See \"logiqctl config help\" for more details")
			os.Exit(-1)
		}
	}

	return client
}

type Api interface {
	MakeApiCall(method string, url string, payload *bytes.Buffer) (*http.Response, error)
}

type ApiClient struct{}

func getClient() {
	api_key := viper.GetString(utils.AuthToken)
	var client *http.Client
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	if api_key != "" {
		client = &http.Client{}
	} else {
		user := utils.GetUIUser()
		password := utils.GetUIPass()

		if user != "" && password != "" {
			cookieJar, _ := cookiejar.New(nil)

			client = &http.Client{
				Jar: cookieJar,
				//Transport: &http.Transport{
				//	MaxIdleConnsPerHost: 20,
				//},
				//Timeout: 1000 * time.Second,
			}

			loginUrl := GetUrlForResource(ResourceLogin)
			u, _ := url.Parse(loginUrl)
			q, _ := url.ParseQuery(u.RawQuery)
			q.Add("remember", "on")
			q.Add("email", user)
			q.Add("password", password)
			u.RawQuery = q.Encode()

			if resp, err := client.Post(u.String(), "application/x-www-form-urlencoded", bytes.NewReader(([]byte)(q.Encode()))); err != nil {
				fmt.Println("Error login with provided credentials, Error:", err.Error())
				os.Exit(-1)
			} else {
				defer resp.Body.Close()
				b, _ := ioutil.ReadAll(resp.Body)
				if match, _ := regexp.Match("Wrong email or password", b); match {
					fmt.Println("Error credentials")
					os.Exit(-1)
				}
			}
		} else {
			fmt.Println("api token or ui credentials must be set. See \"logiqctl config help\" for more details")
			os.Exit(-1)
		}
	}
}

func (c *ApiClient) MakeApiCall(method string, url string, payload *bytes.Buffer) (*http.Response, error) {
	var req *http.Request
	var err error
	var res *http.Response

	if method == http.MethodPost {
		req, err = utils.CreateHttpRequest(method, url, payload)
	} else {
		req, err = utils.CreateHttpRequest(method, url, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("Error sending request")
	}

	if utils.FlagNetTrace == true {
		req = utils.AddNetTrace(req)
	}

	client := getHttpClient()
	res, err = client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error: %s", err.Error())
	}

	return res, nil
}
