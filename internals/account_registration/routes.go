package accountRegistration

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/flukis/inboice/services/domain"
	httpresponse "github.com/flukis/inboice/services/utils/http_response"
	"github.com/go-chi/chi/v5"
)

type HttpHandler struct {
	svc RegisterAccount
}

func NewHttpHandler(svc RegisterAccount) *HttpHandler {
	return &HttpHandler{svc: svc}
}

func (h *HttpHandler) Route(r *chi.Mux) {
	r.Post("/api/register", h.registration)
}

func (h *HttpHandler) registration(
	w http.ResponseWriter,
	r *http.Request,
) {
	var body domain.RegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		fmt.Println(err)
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	_, err := h.svc.Register(ctx, body)
	if err != nil {
		if errors.Is(err, domain.ErrAccountEmailAlreadyRegistered) {
			httpresponse.WriteError(w, http.StatusConflict, err)
			return
		}
		httpresponse.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	httpresponse.WriteMessage(w, http.StatusCreated, "registration success, please check your email for verification")
}
