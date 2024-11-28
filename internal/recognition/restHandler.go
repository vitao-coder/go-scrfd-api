package recognition

import (
	"encoding/json"
	"go-scrfd-api/internal/infrastructure/rest"
	"go.uber.org/fx"
	"io"
	"net/http"
)

type HandlerRest struct {
	fx.Out
	http.Handler
	uc UseCase
}

func NewRESTHandler(recognitionUC UseCase) rest.Handler {
	return &HandlerRest{uc: recognitionUC}
}

func (h *HandlerRest) Pattern() string {
	return "/recognize"
}

func (h *HandlerRest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	allBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	boxes, err := h.uc.Detect(allBytes)
	if err != nil {
		http.Error(w, "Detection failed", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(boxes)
}
