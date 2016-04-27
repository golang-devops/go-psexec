package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func checkResponse(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		if b, e := ioutil.ReadAll(resp.Body); e != nil {
			return fmt.Errorf("The server returned status code %d but could not read response body. Error: %s", resp.StatusCode, e.Error())
		} else {
			return fmt.Errorf("Server status code %d with response %s", resp.StatusCode, string(b))
		}
	}
	return nil
}
