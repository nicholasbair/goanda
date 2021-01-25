package goanda

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func checkApiErr(body []byte, route string) error {
	bodyString := string(body[:])
	if strings.Contains(bodyString, "errorMessage") {
		return fmt.Errorf("OANDA API Error:  %s. Route: %+v", bodyString, route)
	}

	return nil
}

func unmarshalJson(body []byte, data interface{}) {
	jsonErr := json.Unmarshal(body, &data)
	checkErr(jsonErr)
}

func createUrl(host string, endpoint string) string {
	var buffer bytes.Buffer
	// Generate the auth header
	buffer.WriteString(host)
	buffer.WriteString(endpoint)

	url := buffer.String()
	return url
}

func makeRequest(c *OandaConnection, endpoint string, client http.Client, req *http.Request) ([]byte, error) {
	req.Header.Set("User-Agent", c.headers.agent)
	req.Header.Set("Authorization", c.headers.auth)
	req.Header.Set("Content-Type", c.headers.contentType)

	res, getErr := client.Do(req)
	if getErr != nil {
		return []byte("{}"), getErr
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if getErr != nil {
		return []byte("{}"), readErr
	}

	err := checkApiErr(body, endpoint)
	return body, err
}
