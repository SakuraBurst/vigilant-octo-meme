package tests

import (
	"encoding/json"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/domain/models"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"strconv"
	"testing"
	"time"
)

type BannerResponse struct {
	Id int `json:"id"`
}
type TokenResponse struct {
	Token string `json:"token"`
}

func TestBannerCrud_HappyPath(t *testing.T) {
	apiClient := resty.New()
	apiClient.SetBaseURL("http://localhost:8080")

	adminTokenResponse, err := apiClient.R().SetQueryParam("isAdmin", "true").Get("/user_token")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, adminTokenResponse.StatusCode())
	require.NotEmpty(t, adminTokenResponse.String())
	tokenResponse := TokenResponse{}
	err = json.Unmarshal(adminTokenResponse.Body(), &tokenResponse)
	adminToken := tokenResponse.Token
	require.NoError(t, err)

	bannerContentBytes, err := gofakeit.JSON(&gofakeit.JSONOptions{
		Type:     "object",
		RowCount: 3,
		Fields: []gofakeit.Field{
			{Name: "id", Function: "autoincrement"},
			{Name: "first_name", Function: "firstname"},
			{Name: "last_name", Function: "lastname"},
		},
	})
	require.NoError(t, err)
	bannerContent := map[string]interface{}{}
	err = json.Unmarshal(bannerContentBytes, &bannerContent)
	require.NoError(t, err)
	tagID := gofakeit.IntN(100)
	tagIDString := strconv.Itoa(tagID)
	featureID := gofakeit.IntN(100)
	featureIDString := strconv.Itoa(featureID)
	banner := models.Banner{
		Tags:     []int{tagID},
		Feature:  featureID,
		Content:  bannerContent,
		IsActive: true,
	}

	createBannerTime := time.Now()

	createBannerResponse, err := apiClient.R().SetBody(banner).SetHeader("token", adminToken).Post("/banner")
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, createBannerResponse.StatusCode())
	assert.NotEmpty(t, createBannerResponse.String())
	bannerResponse := BannerResponse{}
	err = json.Unmarshal(createBannerResponse.Body(), &bannerResponse)
	require.NoError(t, err)
	assert.NotZero(t, bannerResponse.Id)

	getAllBannersResponse, err := apiClient.R().SetHeader("token", adminToken).SetQueryParam("feature_id", featureIDString).SetQueryParam("tag_id", tagIDString).Get("/banner")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, getAllBannersResponse.StatusCode())
	assert.NotEmpty(t, getAllBannersResponse.String())
	var banners []models.Banner
	err = json.Unmarshal(getAllBannersResponse.Body(), &banners)
	require.NoError(t, err)
	assert.Len(t, banners, 1)
	assert.Equal(t, tagID, banners[0].Tags[0])
	assert.Equal(t, featureID, banners[0].Feature)
	assert.Equal(t, bannerContent, banners[0].Content)
	assert.True(t, banners[0].IsActive)
	assert.InDelta(t, createBannerTime.Unix(), banners[0].CreatedAt.Unix(), 5)
	assert.InDelta(t, createBannerTime.Unix(), banners[0].UpdatedAt.Unix(), 5)

	getUserBannerResponse, err := apiClient.R().SetHeader("token", adminToken).SetQueryParam("feature_id", featureIDString).SetQueryParam("tag_id", tagIDString).SetQueryParam("use_last_revision", "false").Get("/user_banner")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, getUserBannerResponse.StatusCode())
	assert.NotEmpty(t, getUserBannerResponse.String())

	var userBanner map[string]interface{}
	err = json.Unmarshal(getUserBannerResponse.Body(), &userBanner)
	require.NoError(t, err)
	assert.Equal(t, bannerContent, userBanner)
	newBannerContentBytes, err := gofakeit.JSON(&gofakeit.JSONOptions{
		Type:     "object",
		RowCount: 3,
		Fields: []gofakeit.Field{
			{Name: "id", Function: "autoincrement"},
			{Name: "first_name", Function: "firstname"},
			{Name: "last_name", Function: "lastname"},
		},
	})
	require.NoError(t, err)
	newBannerContent := map[string]interface{}{}
	err = json.Unmarshal(newBannerContentBytes, &newBannerContent)
	require.NoError(t, err)
	banner.Content = newBannerContent

	updateTime := time.Now()

	updateBannerResponse, err := apiClient.R().SetBody(banner).SetHeader("token", adminToken).Patch("/banner/" + strconv.Itoa(bannerResponse.Id))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, updateBannerResponse.StatusCode())
	assert.Equal(t, "OK", updateBannerResponse.String())

	getAllBannersResponse, err = apiClient.R().SetHeader("token", adminToken).SetQueryParam("feature_id", featureIDString).SetQueryParam("tag_id", tagIDString).Get("/banner")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, getAllBannersResponse.StatusCode())
	assert.NotEmpty(t, getAllBannersResponse.String())
	err = json.Unmarshal(getAllBannersResponse.Body(), &banners)
	require.NoError(t, err)
	assert.Len(t, banners, 1)
	assert.Equal(t, tagID, banners[0].Tags[0])
	assert.Equal(t, featureID, banners[0].Feature)
	assert.Equal(t, newBannerContent, banners[0].Content)
	assert.True(t, banners[0].IsActive)
	assert.InDelta(t, createBannerTime.Unix(), banners[0].CreatedAt.Unix(), 5)
	assert.InDelta(t, updateTime.Unix(), banners[0].UpdatedAt.Unix(), 5)

	getUserBannerResponse, err = apiClient.R().SetHeader("token", adminToken).SetQueryParam("feature_id", featureIDString).SetQueryParam("tag_id", tagIDString).SetQueryParam("use_last_revision", "false").Get("/user_banner")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, getUserBannerResponse.StatusCode())
	assert.NotEmpty(t, getUserBannerResponse.String())
	err = json.Unmarshal(getUserBannerResponse.Body(), &userBanner)
	require.NoError(t, err)
	assert.Equal(t, bannerContent, userBanner)

	getUserBannerResponse, err = apiClient.R().SetHeader("token", adminToken).SetQueryParam("feature_id", featureIDString).SetQueryParam("tag_id", tagIDString).SetQueryParam("use_last_revision", "true").Get("/user_banner")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, getUserBannerResponse.StatusCode())
	assert.NotEmpty(t, getUserBannerResponse.String())
	err = json.Unmarshal(getUserBannerResponse.Body(), &userBanner)
	require.NoError(t, err)
	assert.Equal(t, newBannerContent, userBanner)

	deleteBannerResponse, err := apiClient.R().SetHeader("token", adminToken).Delete("/banner/" + strconv.Itoa(bannerResponse.Id))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, deleteBannerResponse.StatusCode())
	assert.Equal(t, "OK", deleteBannerResponse.String())

	getAllBannersResponse, err = apiClient.R().SetHeader("token", adminToken).SetQueryParam("feature_id", featureIDString).SetQueryParam("tag_id", tagIDString).Get("/banner")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, getAllBannersResponse.StatusCode())
	assert.NotEmpty(t, getAllBannersResponse.String())
	err = json.Unmarshal(getAllBannersResponse.Body(), &banners)
	require.NoError(t, err)
	assert.Len(t, banners, 0)

}
