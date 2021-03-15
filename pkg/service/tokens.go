package service

import (
	"fmt"
	"time"

	"github.com/KazakMaxim/auth-jwt-go/models"
	"github.com/KazakMaxim/auth-jwt-go/pkg/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type TokensService struct {
	repos *repository.Repository
}

func NewTokensService(repos *repository.Repository) *TokensService {
	return &TokensService{repos: repos}
}

type tokenClaims struct {
	jwt.StandardClaims
	UserGuid  string
	Guid      string
	TypeToken string
}

func (s *TokensService) GetTokensByGuid(userGuid string) ([]string, error) {
	user := s.repos.FindUserByGuid(userGuid)
	if !user {
		return nil, fmt.Errorf("User is not found")
	}
	getTokens, err := s.GenerateTokens(userGuid)
	if err != nil {
		return nil, err
	}

	return getTokens, nil
}

func (s *TokensService) GenerateTokens(userGuid string) ([]string, error) {

	guidTokenPaire := uuid.New().String()

	generateAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(viper.GetDuration("tokens.access.ttl") * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userGuid,
		guidTokenPaire,
		viper.GetString("tokens.access.type"),
	})

	accessToken, accessErr := generateAccessToken.SignedString([]byte(viper.GetString("tokens.key")))
	if accessErr != nil {
		return nil, accessErr
	}

	generateRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(viper.GetDuration("tokens.refresh.ttl") * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userGuid,
		guidTokenPaire,
		viper.GetString("tokens.refresh.type"),
	})

	refreshToken, refreshErr := generateRefreshToken.SignedString([]byte(viper.GetString("tokens.key")))
	if refreshErr != nil {
		return nil, refreshErr
	}

	bcryptToken, bcryptErr := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if bcryptErr != nil {
		return nil, bcryptErr
	}

	saveRefresh := models.RefreshTable{
		User_guid:     userGuid,
		Refresh_token: string(bcryptToken),
	}

	deleteErr := s.repos.DeleteRefersh(userGuid)

	if deleteErr != nil {
		return nil, deleteErr
	}

	if saveRefreshErr := s.repos.SaveRefresh(saveRefresh); saveRefreshErr != nil {
		return nil, saveRefreshErr
	}

	tokens := []string{userGuid, accessToken, refreshToken}

	return tokens, nil
}

func (s *TokensService) NewTokens(tokens models.Tokens) ([]string, error) {

	accessToken, _ := jwt.Parse(tokens.Access, nil)
	if accessToken == nil {
		return nil, fmt.Errorf("Access token is nil")
	}

	claimsOldAccessToken, ok := accessToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Claims access is not result")
	}

	if claimsOldAccessToken["TypeToken"] != viper.GetString("tokens.access.type") {
		return nil, fmt.Errorf("Invalid access token")
	}

	refreshToken, refreshErr := jwt.Parse(tokens.Refresh, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpectef signing metod")
		}

		return []byte(viper.GetString("tokens.key")), nil
	})
	if refreshErr != nil {
		return nil, refreshErr
	}

	claimsOldRefreshToken, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Claims refresh is not result")
	}

	findRefresh := s.repos.FindRefersh(tokens.Refresh, claimsOldRefreshToken["UserGuid"].(string))
	if findRefresh != nil {
		return nil, fmt.Errorf("Refresh token is not found")
	}

	if claimsOldRefreshToken["TypeToken"] != viper.GetString("tokens.refresh.type") ||
		claimsOldAccessToken["Guid"] != claimsOldRefreshToken["Guid"] {
		return nil, fmt.Errorf("Invalid tokens")
	}

	newTokens, newTokensErr := s.GenerateTokens(claimsOldRefreshToken["UserGuid"].(string))
	if newTokensErr != nil {
		return nil, newTokensErr
	}

	return newTokens, nil
}
