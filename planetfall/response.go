package planetfall

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/errorreporting"
)

// Behavior to use when an error in encountered in a http handler
// Returns the message as a response to the request performer
// Log the message, the error, and the request using the Error Reporting tool from GCloud
func (s *Server) InternalServerError(w http.ResponseWriter, r *http.Request, err error, message string) {
	err = fmt.Errorf("%s: %v", message, err)
	s.errorReporting.Report(errorreporting.Entry{
		Error: err,
		Req:   r,
	})
	log.Println(err)

	http.Error(w, message, http.StatusInternalServerError)
}

// Returns an object as a JSON response
func (s *Server) JsonResponse(w http.ResponseWriter, r *http.Request, v interface{}) {

	out, err := json.Marshal(&v)
	if err != nil {
		s.InternalServerError(w, r, err, "failed to format JSON response")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(out))
}
