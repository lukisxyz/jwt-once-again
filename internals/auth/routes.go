package authentication

import (
	"encoding/json"
	"errors"
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
	r.Post("/api/login", h.login)
}

func (h *HttpHandler) login(
	w http.ResponseWriter,
	r *http.Request,
) {
	var body domain.RegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}
	ctx := r.Context()
	acc, err := h.svc.Login(ctx, body.Email, body.Password)
	if err != nil {
		if errors.Is(err, domain.ErrPasswordNotMatch) {
			httpresponse.WriteError(w, http.StatusUnauthorized, err)
			return
		}
		httpresponse.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	httpresponse.WriteResponse(w, http.StatusOK, acc)
}
