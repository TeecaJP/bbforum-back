package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/lestrrat-go/jwx/jwk"
)

func (s *Server) VerifyUser(userid string, w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie("auth")
	if err != nil {
		respondErr(w, r, http.StatusBadRequest, "please set token to cookie")
		return false
	}

	type AWSCognitoClaims struct {
		Name string `json:"cognito:username"`
		jwt.StandardClaims
	}
	var JWTResult AWSCognitoClaims

	jwt.ParseWithClaims(cookie.Value, &JWTResult, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})

	user, err := s.findUserByID(userid)
	if err != nil {
		respondErr(w, r, http.StatusBadRequest, fmt.Sprintf("unknown user: %s", userid))
		return false
	}
	return JWTResult.Name == user.ID
}

func (s *Server) GetCurrentUser(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie("auth")
	if err != nil {
		respondErr(w, r, http.StatusBadRequest, "please set token to cookie")
		return "", err
	}

	type AWSCognitoClaims struct {
		Name string `json:"cognito:username"`
		jwt.StandardClaims
	}
	var JWTResult AWSCognitoClaims

	jwt.ParseWithClaims(cookie.Value, &JWTResult, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})

	return JWTResult.Name, nil
}

func NeedToken(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth")
		if err != nil {
			respondErr(w, r, http.StatusBadRequest, "please set token to cookie", err)
			return
		}
		publicKeysURL := "https://cognito-idp." + os.Getenv("AWS_REGION") + ".amazonaws.com/" + os.Getenv("AWS_COGNITO_USER_POOL_ID") + "/.well-known/jwks.json"

		publicKeySet, err := jwk.Fetch(context.TODO(), publicKeysURL)
		if err != nil {
			log.Printf("failed to parse key: %s", err)
		}

		type AWSCognitoClaims struct {
			Name string `json:"cognito:username"`
			jwt.StandardClaims
		}

		token, err := jwt.ParseWithClaims(cookie.Value, &AWSCognitoClaims{}, func(token *jwt.Token) (interface{}, error) {

			_, ok := token.Method.(*jwt.SigningMethodRSA)

			if !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			kid, ok := token.Header["kid"].(string)
			if !ok {
				return nil, errors.New("kid header not found")
			}

			_, ok = token.Claims.(*AWSCognitoClaims)
			if !ok {
				return nil, errors.New("there is problem to get claims")
			}

			keys, ok := publicKeySet.LookupKeyID(kid)
			if !ok {
				return nil, fmt.Errorf("key %v not found", kid)
			}

			var tokenKey interface{}
			if err := keys.Raw(&tokenKey); err != nil {
				return nil, errors.New("failed to create token key")
			}

			return tokenKey, nil
		})

		if err != nil {
			respondErr(w, r, http.StatusUnauthorized, "token problem")
			return
		}

		if !token.Valid {
			respondErr(w, r, http.StatusUnauthorized, "token is invalid")
			return
		}

		fn(w, r)
	}
}

func SetTokenToCookie(w http.ResponseWriter, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth",
		Value:    value,
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
		SameSite: 4,
		Expires:  time.Now().Add(365 * 24 * time.Hour),
	})
}
