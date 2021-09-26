package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)


var jwtKey=[]byte("secret_key") 
var users=map[string]string{
	"user1":"password1",
	"user2":"password2",
}

type Credentials struct {
	Username string `json:"username"`
	Password string	`json:"password"`
}
// claim creates payload in JWT when token is expired
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
// 1
func Login(w http.ResponseWriter,r *http.Request) {
	var credentials Credentials
	err :=json.NewDecoder(r.Body).Decode(&credentials)
	if err!=nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	expectedPassword,ok :=users[credentials.Username]
	// if username is not in map or password doesn't map the one in map
	if !ok || expectedPassword!=credentials.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// If pass verifying username % password
	// token will be expired after 5 minutes
	expirationTime :=time.Now().Add(time.Minute*5)
	// Create claim
	claims :=&Claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	// create token
	token :=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	tokenString,err :=token.SignedString(jwtKey)
	if err!=nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// If got tokenString sucessfully, set cookie
	http.SetCookie(w,&http.Cookie{
		Name: "token",
		Value: tokenString,
		Expires: expirationTime,
	})
}

func Home(w http.ResponseWriter,r *http.Request)  {
	cookie,err :=r.Cookie("token")
	if err!=nil {
		if err==http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Get tokenstring from requested cookie
	tokenStr :=cookie.Value
	claims :=&Claims{}
	tkn,err :=jwt.ParseWithClaims(tokenStr,claims,
	func(t *jwt.Token)(interface{},error){
		return jwtKey,nil
	})
	if err!=nil {
		if err==jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// if pass cookie check
	w.Write([]byte(fmt.Sprintf("Hello %s",claims.Username)))
}

func Refresh(w http.ResponseWriter,r *http.Request){
	// The same as Home:
	cookie,err :=r.Cookie("token")
	if err!=nil {
		if err==http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Get tokenstring from requested cookie
	tokenStr :=cookie.Value
	claims :=&Claims{}
	tkn,err :=jwt.ParseWithClaims(tokenStr,claims,
	func(t *jwt.Token)(interface{},error){
		return jwtKey,nil
	})
	if err!=nil {
		if err==jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Different from Home:
	// if token expires more than 30 seconds after, refresh token can't be executed
	// For test purpose, block temporaryly
	/* if time.Unix(claims.ExpiresAt,0).Sub(time.Now())>30*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return
	} */
	// create new token copied from Login
	expirationTime :=time.Now().Add(time.Minute*5)
	claims.ExpiresAt=expirationTime.Unix()
	// No need creating claim
	/* claims :=&Claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	} */
	// create token
	token :=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	tokenString,err :=token.SignedString(jwtKey)
	if err!=nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// If got tokenString sucessfully, set cookie
	http.SetCookie(w,&http.Cookie{
		// change token name
		Name: "refresh_token",
		Value: tokenString,
		Expires: expirationTime,
	})
}