package account

import (
	"context"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

//retorna un handler de http con nuestros endpoints
func NewHTTPServer(ctx context.Context, endpoints Endpoints) http.Handler {

	r := mux.NewRouter()
	r.Use(commonMiddleware)
	// Esto nos permite usar los endpoints en los request http
	r.Methods("POST").Path("/newuser").Handler(httptransport.NewServer(
		endpoints.CreateUser,
		decodeUserReq,
		encodeResponse,
	))

	r.Methods("GET").Path("/user/{id}").Handler(httptransport.NewServer(
		endpoints.GetUser,
		decodeEmailReq,
		encodeResponse,
	))

	r.Methods("GET").Path("/users").Handler(httptransport.NewServer(
		endpoints.GetUsers,
		decodeUsersReq,
		encodeResponse,
	))

	r.Methods("PUT").Path("/auth").Handler(httptransport.NewServer(
		endpoints.ValidateUser,
		decodeValidateReq,
		encodeResponse,
	))

	r.Methods("PUT").Path("/validate").Handler(httptransport.NewServer(
		endpoints.ValidateToken,
		decodeValidateTokenReq,
		encodeResponse,
	))

	r.Methods("PUT").Path("/repassword").Handler(httptransport.NewServer(
		endpoints.NewPassword,
		decodeNewPasswordReq,
		encodeResponse,
	))

	r.Methods("GET").Path("/close").Handler(httptransport.NewServer(
		endpoints.Close,
		decodeCloseReq,
		encodeResponse,
	))
	return r

}

// middleware que se asegura de que usemos JSON
func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
