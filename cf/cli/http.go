package cli

import (
	"bytes"
	"code.cloudfoundry.org/cli/plugin"
	"crypto/tls"
	"github.com/happytobi/cf-puppeteer/ui"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	print "github.com/happytobi/cf-puppeteer/cf/utils"
)

//Calls interface
type HttpCalls interface {
	PostJSON(path string, body []byte) (string, error)
}

//HttpConnection
type HttpConnection struct {
	httpClient    *http.Client
	cliConnection plugin.CliConnection
	traceLogging  bool
}

//NewHttpClient ff
func NewHttpClient(cliConnection plugin.CliConnection, traceLogging bool, timeout int, skipSSLValidation bool) *HttpConnection {
	timeoutDuration, _ := time.ParseDuration(string(timeout))

	return &HttpConnection{
		cliConnection: cliConnection,
		traceLogging:  traceLogging,
		httpClient:    setupHttpClient(timeoutDuration, skipSSLValidation),
	}
}

//setup default http client to send http requests
func setupHttpClient(timout time.Duration, skipSSLValidation bool) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: skipSSLValidation,
		},
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			KeepAlive: 30 * time.Second,
			Timeout:   timout,
		}).DialContext,
	}

	return &http.Client{Transport: tr}
}

func (conn *HttpConnection) PostJSON(path string, body []byte) (string, error) {
	request, err := http.NewRequest(
		"POST",
		path,
		bytes.NewReader(body),
	)

	if err != nil {
		return "", err
	}

	//get access token from cli connection
	token, _ := conn.cliConnection.AccessToken()

	request.Header = http.Header{}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", "client.userAgent")
	request.Header.Set("Content-Type", "application/x-yaml")
	request.Header.Set("Authorization", token)

	if conn.traceLogging {
		ui.Say("try to call path: %s", path)
	}

	res, err := conn.httpClient.Do(request)
	if err != nil {
		ui.Failed("Error while calling the apply manifest url %s - error: %s", path, err)
		return "", err
	}

	result, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	jsonResp := string(result)
	if conn.traceLogging {
		ui.Say("response from get call to path: %s was: %s\n", path, print.PrettyJSON(jsonResp))
	}

	return jsonResp, nil
}
