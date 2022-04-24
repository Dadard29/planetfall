package planetfall

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/errorreporting"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/gorilla/mux"
)

// need the following services activate in console:
// - Secret Manager
// - Error Reporting
type Server struct {
	// link to gcloud serviceName
	projectID   string
	serviceName string

	// access to gcloud services
	secretManager  *secretmanager.Client
	errorReporting *errorreporting.Client

	// router
	router *mux.Router
}

type Route struct {
	Endpoint string
	Handler  func(w http.ResponseWriter, r *http.Request)
	Methods  []string
}

func NewServer(
	projectID string, serviceName string,
	routeList []Route,
) (*Server, error) {

	ctx := context.Background()

	// init secret manager
	secretManager, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create secretmanager client: %v", err)
	}

	// init error reporting
	errorReporting, err := errorreporting.NewClient(ctx, projectID, errorreporting.Config{
		ServiceName: serviceName,
		OnError: func(err error) {
			log.Printf("Could not log error: %v", err)
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create error reporting: %v", err)
	}

	// init router
	router := mux.NewRouter()

	for _, route := range routeList {
		router.HandleFunc(route.Endpoint, route.Handler).Methods(route.Methods...)
	}

	router.HandleFunc("/health", handlerHealth).Methods(http.MethodGet)

	serv := &Server{
		projectID:   projectID,
		serviceName: serviceName,

		secretManager:  secretManager,
		errorReporting: errorReporting,

		router: router,
	}

	return serv, nil
}

func handlerHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

func (s *Server) Listen(addr string) {
	if err := http.ListenAndServe(addr, s.router); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) Close() {
	s.secretManager.Close()
}
