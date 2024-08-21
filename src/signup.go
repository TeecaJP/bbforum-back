package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/google/uuid"
)

type Account struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		s.handleSignUpPost(w, r)
		return
	}
	respondHTTPErr(w, r, http.StatusNotFound)
}

func (s *Server) handleSignUpPost(w http.ResponseWriter, r *http.Request) {
	var account Account
	if err := decodeBody(r, &account); err != nil {
		respondErr(w, r, http.StatusBadRequest, "", err)
		return
	}
	user, err := s.mysqlUserCreate(User{Email: account.Email})
	if err != nil {
		respondErr(w, r, http.StatusInternalServerError, "", err)
		return
	}
	newUserData := &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId: aws.String(os.Getenv("AWS_COGNITO_USER_POOL_ID")),
		Username:   aws.String(user.ID),
	}
	client := cognitoidentityprovider.New(
		session.Must(session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})),
	)
	_, err = client.AdminCreateUser(newUserData)
	if err != nil {
		respondErr(w, r, http.StatusInternalServerError, "cognito problem: account create", err)
		return
	}
	input := &cognitoidentityprovider.AdminSetUserPasswordInput{
		UserPoolId: aws.String(os.Getenv("AWS_COGNITO_USER_POOL_ID")),
		Username:   aws.String(user.ID),
		Password:   aws.String(account.Password),
		Permanent:  aws.Bool(true),
	}
	_, err = client.AdminSetUserPassword(input)
	if err != nil {
		respondErr(w, r, http.StatusInternalServerError, "cognito problem: password change", err)
		return
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

	initResult, err := client.AdminInitiateAuth(params)
	if err != nil {
		respondErr(w, r, http.StatusBadRequest, "", err.Error())
		return
	}
	if initResult == nil || initResult.AuthenticationResult == nil || initResult.AuthenticationResult.IdToken == nil {
		respondErr(w, r, http.StatusBadRequest, "", err.Error())
		return
	}

	SetTokenToCookie(w, *initResult.AuthenticationResult.IdToken)
	respond(w, r, http.StatusOK, user)
}

func (s *Server) mysqlUserCreate(user User) (User, error) {
	_, err := s.findUserByEmail(user.Email)
	if err != sql.ErrNoRows && err != nil {
		return User{}, err
	}
	u, err := uuid.NewRandom()
	if err != nil {
		return User{}, err
	}
	user.ID = u.String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Type = "person"

	query := `INSERT INTO users (id, email, name, image, type, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err = s.db.Exec(query, user.ID, user.Email, user.Name, user.Image, user.Type, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		log.Println(err)
		return User{}, err
	}

	result, err := s.findUserByEmail(user.Email)
	if err != nil {
		return User{}, err
	}

	return result, nil
}
