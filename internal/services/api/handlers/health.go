package handlers

import (
	"net/http"

	"go.uber.org/zap"
)

type healthHandler struct {
	log *zap.SugaredLogger
}

func NewHealthHandler(log *zap.SugaredLogger) *healthHandler {
	return &healthHandler{
		log: log,
	}
}

func (h *healthHandler) CheckStatus(rw http.ResponseWriter, r *http.Request) {
	h.log.Info("Checking API status")
	rw.Write([]byte("status ok"))
	// _, err := rw.Write([]byte("status ok"))
	// if err != nil {
	// 	log.Printf("CheckStatus - write failed: %v", err)
	// }
}
