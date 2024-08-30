package tests

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"grpc-sso/internal/domain/models"
	"grpc-sso/internal/grpc/proto/sso"
	"grpc-sso/tests/suite"
	"testing"
	"time"
)

const (
	appID     = 1
	appSecret = "test-secret"

	passDefaultLen = 10
)

func TestRegisterLogin_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePassword()

	registerResponse, err := st.AuthClient.Register(ctx,
		&sso.RegisterRequest{
			Email:    email,
			Password: pass,
		})

	require.NoError(t, err)
	require.NotNil(t, registerResponse)
	assert.NotEqual(t, registerResponse.GetUserId(), models.EmptyUserID)

	loginResponse, err := st.AuthClient.Login(ctx,
		&sso.LoginRequest{
			Email:    email,
			Password: pass,
			AppId:    appID,
		})

	require.NoError(t, err)
	require.NotNil(t, loginResponse)

	loginTime := time.Now()

	token := loginResponse.GetToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, registerResponse.GetUserId(), int64(claims["user_id"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))

	const deltaSeconds = 1

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["expires"].(float64), deltaSeconds)
}

func TestRegisterLogin_DoubleRegistration(t *testing.T) {
	ctx, st := suite.New(t)
	email := gofakeit.Email()
	pass := randomFakePassword()

	registerRequest := sso.RegisterRequest{
		Email:    email,
		Password: pass,
	}

	// First registration
	registerResponse, err := st.AuthClient.Register(ctx,
		&registerRequest)

	require.NoError(t, err)
	require.NotNil(t, registerResponse)
	assert.NotEqual(t, registerResponse.GetUserId(), models.EmptyUserID)

	// Try to register same user
	registerResponse, err = st.AuthClient.Register(ctx,
		&registerRequest)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "user already exists")
	assert.Nil(t, registerResponse)
}

func TestRegisterIsAdmin_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)
	email := gofakeit.Email()
	pass := randomFakePassword()

	registerRequest := sso.RegisterRequest{
		Email:    email,
		Password: pass,
	}

	registerResponse, err := st.AuthClient.Register(ctx,
		&registerRequest)

	require.NoError(t, err)
	require.NotNil(t, registerResponse)

	userID := registerResponse.GetUserId()
	assert.NotEqual(t, userID, models.EmptyUserID)

	isAdminRequest := sso.IsAdminRequest{UserId: userID}

	isAdminResponse, err := st.AuthClient.IsAdmin(ctx, &isAdminRequest)
	require.NoError(t, err)
	assert.NotNil(t, isAdminResponse)
	assert.Equal(t, isAdminResponse.GetIsAdmin(), false)
}

func TestIsAdmin_UserNotFound(t *testing.T) {
	ctx, st := suite.New(t)

	userID := int64(-1)

	isAdminRequest := sso.IsAdminRequest{UserId: userID}

	isAdminResponse, err := st.AuthClient.IsAdmin(ctx, &isAdminRequest)

	require.Error(t, err)
	assert.ErrorContains(t, err, "user not found")
	assert.Nil(t, isAdminResponse)
}

func TestRegister_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		expectedErr string
	}{
		{
			name:        "Register with empty email",
			email:       "",
			password:    randomFakePassword(),
			expectedErr: "email is required",
		},
		{
			name:        "Register with empty password",
			email:       gofakeit.Email(),
			password:    "",
			expectedErr: "password is required",
		},
		{
			name:        "Register with empty email and empty password",
			email:       "",
			password:    "",
			expectedErr: "email is required",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registerRequest := sso.RegisterRequest{
				Email:    test.email,
				Password: test.password,
			}

			registerResponse, err := st.AuthClient.Register(ctx,
				&registerRequest)

			require.Error(t, err)
			require.ErrorContains(t, err, test.expectedErr)
			require.Nil(t, registerResponse)
		})
	}
}

func randomFakePassword() string {
	return gofakeit.Password(
		true,
		true,
		true,
		true,
		false,
		passDefaultLen)
}
