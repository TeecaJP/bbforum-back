package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

func (s *Server) HandleSignIn(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		s.handleSignInPOST(w, r)
		return
	}
	respondHTTPErr(w, r, http.StatusNotFound)
}

func (s *Server) handleSignInPOST(w http.ResponseWriter, r *http.Request) {
	var account Account
	err := decodeBody(r, &account)
	if err != nil {
		respondErr(w, r, http.StatusBadRequest, "", err)
		return
	}
	var user User
	if account.Email != "" {
		user, err = s.findUserByEmail(account.Email)
		if err != nil {
			log.Println(err)
			respondErr(w, r, http.StatusBadRequest, "", err)
			return
		}
	} else {
		user, err = s.findUserByID(account.ID)
		if err != nil {
			log.Println(err)
			respondErr(w, r, http.StatusBadRequest, "", err)
			return
		}
	}
	mac := hmac.New(sha256.New, []byte(os.Getenv("AWS_COGNITO_APP_CLIENT_SECRET")))
	mac.Write([]byte(*aws.String(user.ID) + os.Getenv("AWS_COGNITO_APP_CLIENT_ID")))
	secretHash := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	params := &cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow: aws.String(cognitoidentityprovider.AuthFlowTypeAdminNoSrpAuth),
		AuthParameters: map[string]*string{
			"USERNAME":    aws.String(user.ID),
			"PASSWORD":    aws.String(account.Password),
			"SECRET_HASH": &secretHash,
		},
		ClientId:   aws.String(os.Getenv("AWS_COGNITO_APP_CLIENT_ID")),
		UserPoolId: aws.String(os.Getenv("AWS_COGNITO_USER_POOL_ID")),
	}
	client := cognitoidentityprovider.New(
		session.Must(session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})),
	)

	res, err := client.AdminInitiateAuth(params)
	if err != nil {
		log.Println(err)
		respondErr(w, r, http.StatusBadRequest, "", err.Error())
		return
	}
	if res == nil || res.AuthenticationResult == nil || res.AuthenticationResult.IdToken == nil {
		log.Println(err)
		respondErr(w, r, http.StatusBadRequest, "", err.Error())
		return
	}

	SetTokenToCookie(w, *res.AuthenticationResult.IdToken)
	respond(w, r, http.StatusOK, user)
}
