package planetfall

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/errorreporting"
)

func (s *Server) InternalServerError(w http.ResponseWriter, r *http.Request, err error, message string) {
	err = fmt.Errorf("%s: %v", message, err)
	s.errorReporting.Report(errorreporting.Entry{
		Error: err,
		Req:   r,
	})
	log.Printf("%s: %v", message, err)

	http.Error(w, message, http.StatusInternalServerError)
}

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
