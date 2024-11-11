package models

import (
	"context"
	"fiber_boilerplate/db"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthModel struct{}

type TokenPayload struct {
	AccessToken  string
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
	AccessToken  string `json:"access_token"`
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
	payload.AccessToken, err = accessToken.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
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

func (a AuthModel) ExtractToken(r *fiber.Ctx) string {
	bearerToken := r.Get("Authorization")
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func (a AuthModel) VerifyToken(r *fiber.Ctx) (*jwt.Token, error) {
	tokenString := a.ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, jwt.ErrInvalidKey
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (a AuthModel) ExtractTokenMetadata(r *fiber.Ctx) (*AccessDetails, error) {
	token, err := a.VerifyToken(r)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}

		return &AccessDetails{
			AccessUUID: accessUUID,
			UserID:     userId,
		}, nil
	}
	return nil, err
}

func (a AuthModel) TokenValid(r *fiber.Ctx) error {
	token, err := a.VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

func (a AuthModel) GetAuth(accessDetail *AccessDetails) (int64, error) {
	ctx := context.Background()
	userId, err := db.GetRedis().Get(ctx, accessDetail.AccessUUID).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseInt(userId, 10, 64)
	return userID, nil
}

func (a AuthModel) DestroyAuth(uuid string) (int64, error) {
	ctx := context.Background()
	destroyed, err := db.GetRedis().Del(ctx, uuid).Result()
	if err != nil {
		return 0, err
	}
	return destroyed, nil
}
