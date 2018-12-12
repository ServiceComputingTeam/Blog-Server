package jwt

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	if r.Method == "PUT" && r.URL.Path == "/user/login" {
		var username = r.URL.Query().Get("username")
		next(w, r)
		if strings.Contains(w.Header().Get("Content-Type"), "application/json") {
			if tokenString, err := newClaims(username); err != nil {
				fmt.Fprint(w, err)
			} else {
				json, err := json.Marshal(Token{tokenString})
				if err != nil {
					fmt.Fprint(w, err)
				} else {
					w.Write(json)
					fmt.Println("Token: ", string(json))
				}
			}
		}
	} else {
		next(w, r)
	}
}

func ValidatorJWT(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.URL.Path == "/user" || r.URL.Path == "/user/login" {
		next(w, r)
		return
	}

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return mySigningKey, nil
	})
	if err != nil {
		next(w, r)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		r.Header.Set("username", claims["username"].(string))
		fmt.Println("username:", claims["username"].(string))
		next(w, r)
	} else {
		next(w, r)
	}
}
