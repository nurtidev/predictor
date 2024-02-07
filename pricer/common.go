package pricer

import (
	jsoniter "github.com/json-iterator/go"
	"io"
	"io/ioutil"
	"net/http"
)

func JSONDecode(data []byte, to interface{}) error {
	err := jsoniter.Unmarshal(data, &to)
	if err != nil {
		return err
	}
	return nil
}

func SendHTTPRequest(method, path string, headers map[string]string, body io.Reader, result interface{}) error {
	var client http.Client
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = JSONDecode(respBody, &result)
	if err != nil {
		return err
	}
	return nil
}
