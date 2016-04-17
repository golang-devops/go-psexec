package client

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/mozillazg/request"

	"github.com/golang-devops/go-psexec/shared"
)

type RequestResponse struct {
	Pid      int
	response *request.Response
}

func (r *RequestResponse) TextResponseChannel() (feedback <-chan string, errors <-chan error) {
	feedbackRW := make(chan string)
	errorsRW := make(chan error)

	scanner := bufio.NewScanner(r.response.Body)
	go func() {
		defer r.response.Body.Close()
		defer close(feedbackRW)
		defer close(errorsRW)

		gotEOF := false
		for scanner.Scan() {
			txt := scanner.Text()
			if strings.Contains(txt, shared.RESPONSE_EOF) {
				gotEOF = true
			}
			feedbackRW <- txt

			/*cipher := scanner.Bytes()
			  plaintextBytes, err := shared.DecryptSymmetric(session.SessionToken(), cipher)
			  if err != nil {
			      return fmt.Errorf("Unable read encrypted server response, error: %s", err.Error())
			  }
			  feedbackRW <- string(plaintextBytes)*/
		}

		if !gotEOF {
			errorsRW <- fmt.Errorf("The EOF string '%s' was not found at the end of the response. Assuming the connection got interrupted.", shared.RESPONSE_EOF)
		}
	}()

	return feedbackRW, errorsRW
}
