package jwt

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

const (
	authorization = "Authorization"
)

var mySigningKey = []byte("secret")

type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type JwtHandler struct {
}

func NewJwt() *JwtHandler {
	return &JwtHandler{}
}

func newClaims(username string) (string, error) {
	var claims = MyClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(1)).Unix(),
			Issuer:    "test",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(mySigningKey)
	return tokenString, err
}

type Token struct {
	Token string `json:"token"`
}

func (h *JwtHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.Method == "POST" && r.URL.Path == "/login" {
		r.ParseForm()
		var username = r.Form.Get("username")
		next(w, r)
		if w.Header().Get("Content-Type") == "application/json" {
			if tokenString, err := newClaims(username); err != nil {
				fmt.Fprint(w, err)
			} else {
				json, err := json.Marshal(Token{tokenString})
				if err != nil {
					fmt.Fprint(w, err)
				} else {
					fmt.Fprint(w, string(json))
				}
			}
		}
	} else {
		next(w, r)
	}
}

func ValidatorJWT(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return mySigningKey, nil
	})

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Unauthorized: %v", err)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		r.ParseForm()
		r.Form.Set("username", claims["username"].(string))
		r.PostForm.Set("username", claims["username"].(string))
		fmt.Println(claims)
		next(w, r)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Unauthorized: %v", err)
	}
}
