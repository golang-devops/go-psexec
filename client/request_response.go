package client

import (
	"github.com/mozillazg/request"
	"io/ioutil"
)

type RequestResponse struct {
	*request.Response
}

func (r *RequestResponse) WaitAndReadAll() ([]byte, error) {
	defer r.Body.Close()
	return ioutil.ReadAll(r.Body)
}
