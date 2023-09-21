package echo

import (
	"StoryGoAPI/DTO"
	"StoryGoAPI/config"
	"StoryGoAPI/repository"
	"StoryGoAPI/service"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

var (
	restApi    *rest
	testConfig = &config.Config{
		DataBase: config.DataBase{
			//URI: "postgres://postgres:mozi_pass@127.0.0.1/story_test",
			//Driver: "pgx",
			URI:    "story:thisisPasswd@tcp(127.0.0.1:3306)/story_test?parseTime=True",
			Driver: "mysql",
		},
		Redis: config.Redis{
			Addr:             "localhost:6379",
			DB:               2,
			ExpTime:          time.Minute * 1,
			BlacklistExpTime: time.Minute * 6,
		},
		Secrets: config.Secrets{
			JwtSecret: "testSecret",
			ExpTime:   time.Minute * 5,
		},
	}
	testUser = DTO.RegisterRequest{
		Email:    "johndoe@gmail.com",
		Password: "StrongP@asw0rd",
		Name:     "John Doe",
	}
	testGuest = DTO.GuestRequest{
		VersionNumber:   1,
		OperatingSystem: "LINUX",
		DisplayDetails:  "1280x720",
		UserAgent:       "golang tests/v1",
	}
	createdStoryID int
	token          string
	guestToken     string
)

func setup() func() {
	cfg, _ := config.ReadConfig("./test_config.yml")
	if cfg != nil {
		testConfig = cfg // replace hardcoded config with config from file
	}

	repo, err := repository.NewRepo(testConfig)
	if err != nil {
		panic(err)
	}
	srv := service.NewService(testConfig, repo)
	restApi = NewRest(testConfig, srv).(*rest)

	restApi.echo.Use(restApi.handler.costumeContext)
	restApi.echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, URI=${uri}, code=${status}, latency=${latency_human}\n",
	}))
	restApi.echo.Pre(middleware.RemoveTrailingSlash())

	// set our costumeContext err handler
	restApi.echo.HTTPErrorHandler = restApi.ErrHandler

	restApi.echo.Debug = true
	restApi.routing()
	return func() {
		if err = repo.Psql.DropAllTables(); err != nil {
			log.Println(err)
		}
		if err = repo.Close(); err != nil {
			log.Println(err)
		}
	}
}

func makeRequest(method, url, token, guestToken string, body interface{}, queryParam url.Values) (*httptest.ResponseRecorder, error) {
	requestBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	if method == http.MethodPost || method == http.MethodPut || method == http.MethodDelete {
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}
	if token != "" {
		request.Header.Add("Authorization", "Bearer "+token)
	}

	if guestToken != "" {
		request.Header.Add("X-Guest-Token", guestToken)
	}

	if queryParam != nil {
		request.URL.RawQuery = queryParam.Encode()
	}

	writer := httptest.NewRecorder()

	restApi.echo.ServeHTTP(writer, request)

	return writer, nil
}

func bearerToken() (string, error) {
	recorder, _ := makeRequest(http.MethodPost, "/api/v1/user/login", "", "", testUser, nil)
	var response map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		return "", err
	}
	return response["token"], nil
}

func TestMain(m *testing.M) {
	//var code int
	//defer func() { os.Exit(code) }()
	tearDown := setup()
	defer tearDown()
	if restApi != nil {
		_ = m.Run()
	} else {
		log.Println("failed to init restapi...")
	}
}

func TestRegister(t *testing.T) {
	writer, err := makeRequest(http.MethodPost, "/api/v1/user/register", "", "", testUser, nil)
	if !assert.NoError(t, err, "does not expected an error") {
		return
	}
	assert.Equal(t, http.StatusCreated, writer.Code)

	var loginResp = &DTO.LoginResponse{}
	err = json.Unmarshal(writer.Body.Bytes(), loginResp)
	if !assert.NoError(t, err, "Failed to unmarshal writer body into DTO.LoginResponse") {
		return
	}

	assert.True(t, loginResp.Token != "")
	assert.Equal(t, loginResp.Email, testUser.Email)

	token = loginResp.Token
}
func TestLogin(t *testing.T) {
	writer, err := makeRequest(http.MethodPost, "/api/v1/user/login", "", "", testUser, nil)
	if !assert.NoError(t, err, "does not expected an error") {
		return
	}
	assert.Equal(t, http.StatusOK, writer.Code)

	var resp = &DTO.LoginResponse{}
	err = json.Unmarshal(writer.Body.Bytes(), resp)
	if !assert.NoError(t, err, "Failed to unmarshal writer body into DTO.LoginResponse") {

	}

	assert.True(t, resp.Token != "")
	assert.Equal(t, resp.Email, testUser.Email)
}
func TestCreateStory(t *testing.T) {
	now := time.Now()
	then := now.Add(time.Hour * 24)
	data := &DTO.StoryRequest{
		StoryID:         0,
		CreatorUserID:   0,
		FromTime:        now,
		ToTime:          then,
		StoryName:       "Wedding",
		BackgroundColor: "#FFFFFF",
		BackgroundImage: "https://example.com/background.jpg",
		IsShareable:     false,
		AttachedFile:    "",
		ExternalWebLink: "",
		Status:          "public",
	}
	writer, err := makeRequest(http.MethodPost, "/api/v1/user/new_story", token, "", data, nil)
	if !assert.NoError(t, err, "Expected no error while making the request") {
		return
	}
	assert.Equal(t, http.StatusCreated, writer.Code, "Expected HTTP status code to be 201 (Created)")

	// Parse the writer body into DTO.StoryResponse
	var resp = DTO.StoryResponse{}
	err = json.Unmarshal(writer.Body.Bytes(), &resp)
	if !assert.NoError(t, err, "Failed to unmarshal writer body into DTO.StoryResponse") {
		return
	}

	createdStoryID = resp.StoryID

	// Assertions on the DTO.StoryResponse fields
	assert.True(t, resp.StoryID > 0, "Expected StoryID to be greater than 0")
	assert.True(t, resp.CreatorUserID > 0, "Expected CreatorUserID to be greater than 0")
	assert.WithinDuration(t, data.FromTime, resp.FromTime, 0, "FromTime field mismatch")
	assert.WithinDuration(t, data.FromTime, resp.FromTime, 0, "FromTime field mismatch")
	assert.WithinDuration(t, data.ToTime, resp.ToTime, 0, "ToTime field mismatch")
	assert.WithinDuration(t, data.ToTime, resp.ToTime, 0, "ToTime field mismatch")
	assert.Equal(t, data.StoryName, resp.StoryName, "StoryName field mismatch")
	assert.Equal(t, data.BackgroundColor, resp.BackgroundColor, "BackgroundColor field mismatch")
	assert.Equal(t, data.BackgroundImage, resp.BackgroundImage, "BackgroundImage field mismatch")
	assert.Equal(t, data.IsShareable, resp.IsShareable, "IsShareable field mismatch")
	assert.Equal(t, data.AttachedFile, resp.AttachedFile, "AttachedFile field mismatch")
	assert.Equal(t, data.ExternalWebLink, resp.ExternalWebLink, "ExternalWebLink field mismatch")
}
func TestGetPostedStories(t *testing.T) {
	query := url.Values{
		"sort_by": {"created"},
		"status":  {"public"},
	}
	writer, err := makeRequest(http.MethodGet, "/api/v1/user/stories", token, "", nil, query)
	if !assert.NoError(t, err, "does not expected an error") {
		return
	}
	assert.Equal(t, http.StatusOK, writer.Code)

	stories := &DTO.StoriesResponse{}
	err = json.Unmarshal(writer.Body.Bytes(), &stories)
	if !assert.NoError(t, err, "Failed to unmarshal writer body into DTO.StoryResponse") {
		return
	}

	assert.Len(t, stories.Stories, 1)

	t.Run("TestUpdateStory", func(t *testing.T) {
		data := stories.Stories[0]
		data.StoryName = "TestTestTest"

		writer, err := makeRequest(http.MethodPut, fmt.Sprintf("/api/v1/user/edit_story/%d", data.StoryID), token, "", data, nil)
		if !assert.NoError(t, err, "does not expected an error") {
			return
		}
		assert.Equal(t, http.StatusOK, writer.Code)

		var resp DTO.StoryResponse
		err = json.Unmarshal(writer.Body.Bytes(), &resp)
		if !assert.NoError(t, err, "Failed to unmarshal writer body into DTO.StoryResponse") {
			return
		}

		assert.Equal(t, data.StoryName, resp.StoryName, "expected to update the story")
	})
}

func TestNewGuest(t *testing.T) {
	writer, err := makeRequest(http.MethodPost, "/api/v1/guest/new", "", "", testGuest, nil)
	if !assert.NoError(t, err, "no error expected") {
		return
	}
	assert.Equal(t, http.StatusCreated, writer.Code, "Expected HTTP status code to be 201 (Created)")

	resp := &DTO.GuestResponse{}
	err = json.Unmarshal(writer.Body.Bytes(), &resp)
	if !assert.NoError(t, err, "Failed to unmarshal writer body into DTO.GuestResponse") {
		return
	}
	assert.True(t, resp.Token != "", "expected a guest guestToken")

	guestToken = resp.Token
}
func TestVerifyGuest(t *testing.T) {
	data := testGuest
	data.Token = guestToken
	writer, err := makeRequest(http.MethodPost, "/api/v1/guest/verify", "", "", data, nil)
	if !assert.NoError(t, err, "no error expected") {
		return
	}
	assert.Equal(t, http.StatusOK, writer.Code, "Expected HTTP status code to be 200 (OK)")

	r := DTO.SuccessResponse{}
	err = json.Unmarshal(writer.Body.Bytes(), &r)
	if !assert.NoError(t, err, "Failed to unmarshal writer body into map") {
		return
	}
	assert.True(t, r.Success)
}
func TestScanStory(t *testing.T) {
	writer, err := makeRequest(http.MethodPost, fmt.Sprintf("/api/v1/guest/scan/%d", createdStoryID), "", guestToken, nil, nil)
	if !assert.NoError(t, err, "does not expected an error") {
		return
	}
	assert.Equal(t, http.StatusOK, writer.Code)

	r := DTO.SuccessResponse{}
	err = json.Unmarshal(writer.Body.Bytes(), &r)
	if !assert.NoError(t, err, "Failed to unmarshal writer body into map") {
		return
	}
	assert.True(t, r.Success)
}
func TestGuestStoryFeed(t *testing.T) {
	writer, err := makeRequest(http.MethodGet, "/api/v1/guest/stories", "", guestToken, nil, nil)
	if !assert.NoError(t, err, "does not expected an error") {
		return
	}
	assert.Equal(t, http.StatusOK, writer.Code)

	stories := &DTO.StoriesResponse{}
	err = json.Unmarshal(writer.Body.Bytes(), &stories)
	if !assert.NoError(t, err, "Failed to unmarshal writer body into DTO.StoryResponse") {
		return
	}

	assert.Len(t, stories.Stories, 1)
}
func TestStoryInfo(t *testing.T) {
	writer, err := makeRequest(http.MethodGet, fmt.Sprintf("/api/v1/story/%d", createdStoryID), "", "", nil, nil)
	if !assert.NoError(t, err, "does not expected an error") {
		return
	}
	assert.Equal(t, http.StatusOK, writer.Code)

	var resp = &DTO.StoryResponse{}
	err = json.Unmarshal(writer.Body.Bytes(), resp)
	if !assert.NoError(t, err, "Failed to unmarshal writer body into DTO.StoryResponse") {

	}

	assert.Equal(t, resp.StoryID, createdStoryID)
}

func TestDeleteGuest(t *testing.T) {
	writer, err := makeRequest(http.MethodDelete, "/api/v1/guest/delete", "", guestToken, nil, nil)
	if !assert.NoError(t, err, "no error expected") {
		return
	}
	assert.Equal(t, http.StatusOK, writer.Code, "Expected HTTP status code to be 200 (OK)")

	r := DTO.SuccessResponse{}
	err = json.Unmarshal(writer.Body.Bytes(), &r)
	if !assert.NoError(t, err, "Failed to unmarshal writer body into map") {
		return
	}
	assert.True(t, r.Success)
}
func TestDeleteStory(t *testing.T) {
	if createdStoryID == 0 {
		t.Errorf("story id is not set")
	}
	u := fmt.Sprintf("/api/v1/user/delete_story/%d", createdStoryID)

	writer, err := makeRequest(http.MethodDelete, u, token, "", nil, nil)
	if !assert.NoError(t, err, "does not expected an error") {
		return
	}
	assert.Equal(t, http.StatusOK, writer.Code)

	r := DTO.SuccessResponse{}
	err = json.Unmarshal(writer.Body.Bytes(), &r)
	if !assert.NoError(t, err, "Failed to unmarshal writer body into map") {
		return
	}
	assert.True(t, r.Success)
}
func TestDeleteAccount(t *testing.T) {
	token, err := bearerToken()
	if !assert.NoError(t, err) {
		return
	}
	writer, err := makeRequest(http.MethodDelete, "/api/v1/user/delete", token, "", testUser, nil)
	if !assert.NoError(t, err, "does not expected an error") {
		return
	}
	assert.Equal(t, http.StatusOK, writer.Code)

	r := DTO.SuccessResponse{}
	err = json.Unmarshal(writer.Body.Bytes(), &r)
	if !assert.NoError(t, err, "Failed to unmarshal writer body into map") {
		return
	}
	assert.True(t, r.Success)
}
