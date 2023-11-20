package app

import (
	"basicproject/models"
	u "basicproject/utils"
	"context"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		
		requestPath := r.URL.Path //current request path
		if requestPath == "/login" {
			next.ServeHTTP(w, r)
			return
		}

		response := make(map[string]interface{})

		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

		if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
			response = u.Message(false, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			response = u.Message(false, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		tokenPart := splitted[1] //Grab the token part, what we are truly interested in
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil { //Malformed token, returns with http code 403 as usual
			response = u.Message(false, "Malformed authentication token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		if !token.Valid { //Token is invalid, maybe not signed on this server
			response = u.Message(false, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}
		var Account models.Account = models.GetAccountFromDB(tk.UserId)

		//current request path
		AdminEndPoints := []string{"/login", "/admin/student", "/admin/faculty", "/admin/slot", "/admin/course"} //List of endpoints that doesn't require auth

		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)

		// next.ServeHTTP(w, r)
		// response = u.BodyBuilder(true)
		// response["data"] = "hehe"
		// w.Header().Add("Content-Type", "application/json")
		// u.Respond(w, response)
		// return

		if Account.Admin {
			for _, val := range AdminEndPoints {
				if val == requestPath {
					next.ServeHTTP(w, r)
					return
				}
			}
			response = u.BodyBuilder(false)
			response["data"] = "You are not allowed on this endpoint"
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}
		var flag bool = false
		if !Account.Admin {
			for _, val := range AdminEndPoints {
				if val == requestPath && val != "/login" {
					flag = true
					break
				}
			}
		}

		if !flag {
			next.ServeHTTP(w, r)
			return
		}
		response = u.BodyBuilder(false)
		response["data"] = "You are not allowed on this endpoint"
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		u.Respond(w, response)
		return
	})
}
