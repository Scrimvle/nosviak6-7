package functions

import (
	"Nosviak4/source/database"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"sync"

	"github.com/dgrijalva/jwt-go"
)

// Token represents a collection of data also known as a token
type Token struct {
	Token       *jwt.Token
	Cookie      *http.Cookie
	Signature   []byte
	TokenClaims jwt.MapClaims
}

// Jwt represents the key value for unlocking tokens
var Jwt = *database.NewSalt(32)

// Tokens represents a collection of data also known as tokens
var Tokens map[string]*Token = make(map[string]*Token)

// mux prevents concurrent map operations for tokens
var mux sync.Mutex

// NewToken will generate a brand new Token
func NewToken(claims jwt.MapClaims) (*Token, error) {
	if _, ok := claims["user"]; !ok {
		return nil, errors.New("attach a user")
	}

	token := new(Token)
	token.Signature = *database.NewSalt(32)
	token.Token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	claims["signature"] = hex.EncodeToString(sha256.New().Sum(token.Signature))

	claimed, err := token.SignJwt(Jwt)
	if err != nil {
		return nil, err
	}

	token.Cookie = &http.Cookie{
		Name: "auth_token",
		Value: claimed,
	}

	mux.Lock()
	defer mux.Unlock()
	Tokens[claims["signature"].(string)] = token
	return token, nil
}


// salt represents the key we use within the JWT token
func (token *Token) SignJwt(salt []byte) (string, error) {
	return token.Token.SignedString(salt)
}

// GetUser will return the owner of the token
func (token *Token) GetUser() (*database.User, error) {
	return database.DB.GetUser(token.Token.Claims.(jwt.MapClaims)["user"].(string))
}

// ExtractToken will attempt to extract the token from the request 
func ExtractToken(r *http.Request) (*Token, error) {
	cookie, err := r.Cookie("auth_token")
	if err != nil || cookie == nil {
		return nil, err
	}

	token, err := jwt.Parse(cookie.Value, signFunc)
	if err != nil || token == nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims == nil {
		return nil, errors.New("unknown token")
	}

	collected, ok := Tokens[claims["signature"].(string)]
	if !ok || collected == nil || cookie.Value != collected.Cookie.Value {
		return nil, errors.New("unknown token")
	}

	return collected, nil
} 

// signFunc will be used within the jwt parsing
func signFunc(recv *jwt.Token) (interface{}, error) {
	if _, ok := recv.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("unexpected signing method")
	}

	return Jwt, nil
}