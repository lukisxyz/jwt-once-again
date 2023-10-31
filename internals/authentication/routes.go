package authentication

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/flukis/inboice/services/domain"
	custommiddleware "github.com/flukis/inboice/services/internals/custom_middleware"
	httpresponse "github.com/flukis/inboice/services/utils/http_response"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

type HttpHandler struct {
	svc    RegisterAccount
	secret string
}

func NewHttpHandler(svc RegisterAccount, secret string) *HttpHandler {
	return &HttpHandler{svc: svc, secret: secret}
}

func (h *HttpHandler) Route(r *chi.Mux) {
	r.Post("/api/login", h.login)
	r.Post("/api/request-token", h.requestaccesstoken)
	r.Group(func(route chi.Router) {
		route.Use(custommiddleware.AuthJwtMiddleware)
		route.Get("/api/protected", h.protected)
		route.Post("/api/logout", h.logout)
	})
}

func (h *HttpHandler) requestaccesstoken(
	w http.ResponseWriter,
	r *http.Request,
) {
	reqToken := r.Header.Get("Authorization")
	splittedToken := strings.Split(reqToken, "Bearer ")
	ctx := r.Context()
	if len(splittedToken) != 2 {
		httpresponse.WriteError(w, http.StatusUnauthorized, errors.New(http.StatusText((http.StatusUnauthorized))))
		ctx.Done()
		return
	}

	jwtToken := splittedToken[1]
	claims := &domain.ClaimResponse{}
	token, err := jwt.ParseWithClaims(
		jwtToken,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(h.secret), nil
		},
	)

	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		httpresponse.WriteError(w, http.StatusUnauthorized, err)
		ctx.Done()
		return
	}
	claims = token.Claims.(*domain.ClaimResponse)

	accessToken, err := h.svc.RefreshToken(ctx, claims.ID, claims.Name, claims.Email)
	if err != nil {
		httpresponse.WriteError(w, http.StatusUnauthorized, err)
		ctx.Done()
		return
	}

	if err != nil {
		httpresponse.WriteError(w, http.StatusUnauthorized, err)
		ctx.Done()
		return
	}

	httpresponse.WriteResponse(w, http.StatusOK, accessToken)
}

func (h *HttpHandler) protected(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()

	key := custommiddleware.UserValueKey
	data := ctx.Value(key).(*domain.ClaimResponse)

	httpresponse.WriteResponse(w, http.StatusOK, data.Email)
}

func (h *HttpHandler) logout(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()

	key := custommiddleware.UserValueKey
	data := ctx.Value(key).(*domain.ClaimResponse)

	if err := h.svc.Logout(ctx, data.ID); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}

	httpresponse.WriteMessage(w, http.StatusOK, "success logout")
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
		if errors.Is(err, domain.ErrAlreadyLogin) {
			httpresponse.WriteError(w, http.StatusUnprocessableEntity, err)
			return
		}
		httpresponse.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	httpresponse.WriteResponse(w, http.StatusOK, acc)
}
