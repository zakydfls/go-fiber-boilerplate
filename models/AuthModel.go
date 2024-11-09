package models

import (
	"context"
	"fiber_boilerplate/db"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthModel struct{}

type TokenPayload struct {
	AccesToken   string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpiresAt  int64
	RtExpiresAt  int64
}

type AccessDetails struct {
	AccessUUID string
	UserID     int64
}

type Token struct {
	AccesToken   string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (a AuthModel) GenerateToken(userId int64) (*TokenPayload, error) {
	payload := &TokenPayload{}

	payload.AtExpiresAt = time.Now().Add(time.Hour * 24).Unix()
	payload.AccessUUID = uuid.New().String()

	payload.RtExpiresAt = time.Now().Add(time.Hour * 24 * 7).Unix()
	payload.RefreshUUID = uuid.New().String()

	var err error

	accessTokenClaim := jwt.MapClaims{
		"iss": os.Getenv("APP_NAME"),
	}
	accessTokenClaim["access_uuid"] = payload.AccessUUID
	accessTokenClaim["is_signed"] = true
	accessTokenClaim["user_id"] = userId
	accessTokenClaim["exp"] = payload.AtExpiresAt

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaim)
	payload.AccesToken, err = accessToken.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	refreshTokenClaim := jwt.MapClaims{
		"iss": os.Getenv("APP_NAME"),
	}
	refreshTokenClaim["refresh_uuid"] = payload.RefreshUUID
	refreshTokenClaim["user_id"] = userId
	refreshTokenClaim["exp"] = payload.RtExpiresAt
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaim)
	payload.RefreshToken, err = refreshToken.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (a AuthModel) StartAuth(userId int64, payload *TokenPayload) error {
	accesToken := time.Unix(payload.AtExpiresAt, 0)
	refreshToken := time.Unix(payload.RtExpiresAt, 0)
	now := time.Now()

	ctx := context.Background()
	errAccess := db.GetRedis().Set(ctx, payload.AccessUUID, strconv.Itoa(int(userId)), accesToken.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := db.GetRedis().Set(ctx, payload.RefreshUUID, strconv.Itoa(int(userId)), refreshToken.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}

	return nil
}

func (a AuthModel) ExtractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}
