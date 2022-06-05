package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	mockdb "github.com/mdanialr/sns_backend/internal/database/mock"
	database "github.com/mdanialr/sns_backend/internal/database/sql"
	"github.com/mdanialr/sns_backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateSend(t *testing.T) {
	const testName = "Should fail when sending empty content-type and return code 400 then has expected error message"
	t.Run(testName, func(t *testing.T) {
		conf := service.Config{}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sns := mockdb.NewMockSNS(ctrl)
		app := fiber.New()
		app.Post("/", CreateSend(&conf, sns))

		req := httptest.NewRequest(fiber.MethodPost, "/", nil)
		req.Header.Add("content-type", "")
		res, _ := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
		var r JsonResponse
		_ = json.NewDecoder(res.Body).Decode(&r)
		assert.Contains(t, r.Msg, "failed to parse json payload")
	})

	t.Run("Should fail when not sending required `file` field in multipart request", func(t *testing.T) {
		conf := service.Config{}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sns := mockdb.NewMockSNS(ctrl)
		app := fiber.New()
		app.Post("/", CreateSend(&conf, sns))

		body := new(bytes.Buffer)
		wr := multipart.NewWriter(body)

		require.NoError(t, wr.Close(), "failed to close multipart writer")

		req := httptest.NewRequest(fiber.MethodPost, "/", body)
		req.Header.Add("content-type", wr.FormDataContentType())
		res, _ := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
		var r JsonResponse
		_ = json.NewDecoder(res.Body).Decode(&r)
		assert.Contains(t, r.Msg, "failed to get file instance from multipart request")
	})

	t.Run("Should fail when using invalid or inaccessible upload dir to save uploaded file", func(t *testing.T) {
		conf := service.Config{UploadDir: "/fake/upload/dir/"}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sns := mockdb.NewMockSNS(ctrl)
		app := fiber.New()
		app.Post("/", CreateSend(&conf, sns))

		body := new(bytes.Buffer)
		wr := multipart.NewWriter(body)
		fl, err := os.Open(fakeFilePath)
		require.NoError(t, err, "failed opening fake file path")
		defer fl.Close()

		form, err := wr.CreateFormFile("file", path.Base(fakeFilePath))
		require.NoError(t, err)

		if _, err := io.Copy(form, fl); err != nil {
			require.NoError(t, err, "failed copy fake file to multipart")
		}

		require.NoError(t, wr.Close(), "failed to close multipart writer")

		req := httptest.NewRequest(fiber.MethodPost, "/", body)
		req.Header.Add("content-type", wr.FormDataContentType())
		res, _ := app.Test(req)

		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
		var r JsonResponse
		_ = json.NewDecoder(res.Body).Decode(&r)
		assert.Contains(t, r.Msg, "failed to save uploaded file to local disk")
	})

	t.Run("Should fail when database failed to create provided data, this error should has nothing to do"+
		" with our code but should be from database themself therefor return code 500", func(t *testing.T) {
		conf := service.Config{UploadDir: "/tmp/"}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sns := mockdb.NewMockSNS(ctrl)
		sns.EXPECT().
			CreateSend(gomock.Any(), gomock.AssignableToTypeOf(database.CreateSendParams{})).
			Times(1).
			Return(database.Send{}, sql.ErrNoRows)

		app := fiber.New()
		app.Post("/", CreateSend(&conf, sns))

		body := new(bytes.Buffer)
		wr := multipart.NewWriter(body)
		fl, err := os.Open(fakeFilePath)
		require.NoError(t, err, "failed opening fake file path")
		defer fl.Close()

		form, err := wr.CreateFormFile("file", path.Base(fakeFilePath))
		require.NoError(t, err)

		if _, err := io.Copy(form, fl); err != nil {
			require.NoError(t, err, "failed copy fake file to multipart")
		}

		require.NoError(t, wr.Close(), "failed to close multipart writer")

		req := httptest.NewRequest(fiber.MethodPost, "/", body)
		req.Header.Add("content-type", wr.FormDataContentType())
		res, _ := app.Test(req)

		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
		var r JsonResponse
		_ = json.NewDecoder(res.Body).Decode(&r)
		assert.Contains(t, r.Msg, "failed to create new send with the given payload")
	})

	t.Run("Should pass when all database operations not returning any errors and all required fields in payload"+
		" is provided", func(t *testing.T) {
		conf := service.Config{UploadDir: "/tmp/"}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sns := mockdb.NewMockSNS(ctrl)
		sns.EXPECT().
			CreateSend(gomock.Any(), gomock.AssignableToTypeOf(database.CreateSendParams{})).
			Times(1).
			Return(database.Send{}, nil)

		app := fiber.New()
		app.Post("/", CreateSend(&conf, sns))

		body := new(bytes.Buffer)
		wr := multipart.NewWriter(body)
		fl, err := os.Open(fakeFilePath)
		require.NoError(t, err, "failed opening fake file path")
		defer fl.Close()

		form, err := wr.CreateFormFile("file", path.Base(fakeFilePath))
		require.NoError(t, err)

		if _, err := io.Copy(form, fl); err != nil {
			require.NoError(t, err, "failed copy fake file to multipart")
		}

		require.NoError(t, wr.Close(), "failed to close multipart writer")

		req := httptest.NewRequest(fiber.MethodPost, "/", body)
		req.Header.Add("content-type", wr.FormDataContentType())
		res, _ := app.Test(req)

		assert.Equal(t, fiber.StatusCreated, res.StatusCode)
		var r JsonResponse
		_ = json.NewDecoder(res.Body).Decode(&r)
		assert.Contains(t, r.Msg, "successfully created")
	})

	t.Cleanup(func() {
		os.Remove(fakeFilePath)
	})
}
