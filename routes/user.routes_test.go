package routes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/s1Sharp/s1-tts-restapi/internal/config"
	"github.com/s1Sharp/s1-tts-restapi/internal/models"
	mockservice "github.com/s1Sharp/s1-tts-restapi/service/mock"
	"github.com/s1Sharp/s1-tts-restapi/utils"
	"github.com/stretchr/testify/require"
)

func TestGetMe(t *testing.T) {
	user := randomUser()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userService := mockservice.NewMockUserService(ctrl)

	// stub
	userService.EXPECT().FindUserById(gomock.Eq(user.ID.String())).Times(1).Return(&user, nil)

	// start test server and send request
	cfg := config.NewTestConfig()
	ctx := context.TODO()

	s := NewTestServer(ctx, cfg, TestServerInternal{userService: userService})

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	newUserAccessToken, _ := utils.CreateToken(cfg.AccessTokenExpired, user.ID.String(), cfg.AccessPrivateKey)
	request.AddCookie(&http.Cookie{
		Name:     "access_token",
		Value:    newUserAccessToken,
		MaxAge:   cfg.AccessTokenMaxAge,
		Path:     "localhost",
		Secure:   false,
		HttpOnly: true,
	})
	require.NoError(t, err)

	router := s.GetRouter()
	router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)
}

func randomUser() models.DBUserResponse {
	return models.DBUserResponse{}
}
