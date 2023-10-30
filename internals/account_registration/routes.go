package accountRegistration

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/flukis/inboice/services/domain"
	httpresponse "github.com/flukis/inboice/services/utils/http_response"
	"github.com/go-chi/chi/v5"
	"github.com/oklog/ulid/v2"
)

type HttpHandler struct {
	svc RegisterAccount
}

func NewHttpHandler(svc RegisterAccount) *HttpHandler {
	return &HttpHandler{svc: svc}
}

func (h *HttpHandler) Route(r *chi.Mux) {
	r.Post("/api/register", h.registration)
	r.Delete("/api/register/{id}", h.deleteuser)
	r.Post("/api/login", h.login)
}

func (h *HttpHandler) registration(
	w http.ResponseWriter,
	r *http.Request,
) {
	var body domain.RegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
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
	acc, err := h.svc.GetByEmail(ctx, body.Email)
	if err != nil {
		if errors.Is(err, domain.ErrAccountEmailAlreadyRegistered) {
			httpresponse.WriteError(w, http.StatusConflict, err)
			return
		}
		httpresponse.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	httpresponse.WriteResponse(w, http.StatusCreated, acc)
}

func (h *HttpHandler) deleteuser(
	w http.ResponseWriter,
	r *http.Request,
) {
	idStr := chi.URLParam(r, "id")
	id, err := ulid.Parse(idStr)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}
	ctx := r.Context()
	err = h.svc.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrAccountAlreadyDeleted) {
			httpresponse.WriteError(w, http.StatusUnprocessableEntity, err)
			return
		}
		httpresponse.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	httpresponse.WriteMessage(w, http.StatusOK, "account deletion success")
}
