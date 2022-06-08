package account

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/shaj13/go-guardian/auth"
	"github.com/shaj13/go-guardian/auth/strategies/ldap"
	"github.com/shaj13/go-guardian/store"
)

var authenticator auth.Authenticator
var cache store.Cache

//retorna un handler de http con nuestros endpoints
func NewHTTPServer(ctx context.Context, endpoints Endpoints) http.Handler {
	setupGoGuardian()
	r := mux.NewRouter()
	miscR := r.PathPrefix("/other").Subrouter()
	//r.Use(commonMiddleware)

	miscR.Use(commonMiddleware)

	miscR.Methods("GET").Path("/user/{id}").Handler(httptransport.NewServer(
		endpoints.GetUser,
		decodeEmailReq,
		encodeResponse,
	))

	miscR.Methods("GET").Path("/users").Handler(httptransport.NewServer(
		endpoints.GetUsers,
		decodeUsersReq,
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

	r.Methods("PUT").Path("/validate").Handler(httptransport.NewServer(
		endpoints.ValidateToken,
		decodeValidateTokenReq,
		encodeResponse,
	))

	loginR := r.Methods(http.MethodPost).Subrouter()
	loginR.Use(PostMiddleware)

	loginR.Methods("POST").Path("/newuser").Handler(httptransport.NewServer(
		endpoints.CreateUser,
		decodeUserReq,
		encodeResponse,
	))

	loginR.Methods("PUT").Path("/auth").Handler(httptransport.NewServer(
		endpoints.ValidateUser,
		decodeValidateReq,
		encodeResponse,
	))

	return r

}

// middleware que se asegura de que usemos JSON
func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing common Middleware")
		fmt.Println("Executing common Middleware")
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
func PostMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing Auth Middleware")
		fmt.Println("Executing Auth Middleware")
		user, err := authenticator.Authenticate(r)
		if err != nil {
			code := http.StatusUnauthorized
			http.Error(w, http.StatusText(code), code)
			return
		}
		log.Printf("User %s Authenticated\n", user.UserName())
		next.ServeHTTP(w, r)
	})
}

func setupGoGuardian() {
	cfg := &ldap.Config{
		//cn=admin,dc=arqsoft,dc=unal,dc=edu,dc=co
		BaseDN:       "dc=unstream",                       //"dc=example,dc=com",
		BindDN:       "cn=admin,dc=arqsoft,dc=unal,dc=co", //"cn=read-only-admin,dc=example,dc=com",
		Port:         "389",
		Host:         "localhost", //"ldap.forumsys.com",//us-ldap //host.docker.internal
		BindPassword: "admin",
		Filter:       "(uid=%s)",
	}
	authenticator = auth.New()
	cache = store.NewFIFO(context.Background(), time.Minute*10)
	strategy := ldap.NewCached(cfg, cache)
	authenticator.EnableStrategy(ldap.StrategyKey, strategy)
}

/*
func main() {
	setupGoGuardian()
	router := mux.NewRouter()
	router.HandleFunc("/v1/book/{id}", middleware(http.HandlerFunc(getBookAuthor))).Methods("GET")
	log.Println("server started and listening on http://127.0.0.1:8080")
	http.ListenAndServe("127.0.0.1:8080", router)
}

func getBookAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	books := map[string]string{
		"1449311601": "Ryan Boyd",
		"148425094X": "Yvonne Wilson",
		"1484220498": "Prabath Siriwarden",
	}
	body := fmt.Sprintf("Author: %s \n", books[id])
	w.Write([]byte(body))
}

func setupGoGuardian() {
	cfg := &ldap.Config{
		BaseDN:       "dc=unstream,dc=com",          //"dc=example,dc=com",
		BindDN:       "cn=admin,dc=unstream,dc=com", //"cn=read-only-admin,dc=example,dc=com",
		Port:         "389",
		Host:         "localhost", //"ldap.forumsys.com",//us-ldap //host.docker.internal
		BindPassword: "admin",
		Filter:       "(uid=%s)",
	}
	authenticator = auth.New()
	cache = store.NewFIFO(context.Background(), time.Minute*10)
	strategy := ldap.NewCached(cfg, cache)
	authenticator.EnableStrategy(ldap.StrategyKey, strategy)
}

func middleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing Auth Middleware")
		user, err := authenticator.Authenticate(r)
		if err != nil {
			code := http.StatusUnauthorized
			http.Error(w, http.StatusText(code), code)
			return
		}
		log.Printf("User %s Authenticated\n", user.UserName())
		next.ServeHTTP(w, r)
	})
}
*/
