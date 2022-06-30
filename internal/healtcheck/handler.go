package healthcheck

import (
	"net/http"

	log "github.com/iliesh/go-templates/logger"
)

func Handler(res http.ResponseWriter, req *http.Request) {
	log.Info("request: %v", req.URL)
}
