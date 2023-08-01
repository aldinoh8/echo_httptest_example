package test

import (
	"encoding/json"
	"example/config"
	"example/model"
	"example/routes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	app *echo.Echo
)

func TestMain(m *testing.M) {
	db = config.InitDatabase("host=localhost user=postgres password=postgres dbname=sandbox_test_db port=5432 sslmode=disable TimeZone=Asia/Jakarta")
	app = routes.InitApp(db)

	db.Exec(`
		INSERT INTO books (name, description)
		VALUES ('test1', 'description1'),
		('test2', 'description2'),
		('test3', 'description3');
	`)

	m.Run()

	db.Exec("TRUNCATE books RESTART IDENTITY")
}

func TestFindAllBooks(t *testing.T) {
	t.Run("Should retreive all books", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/books", nil)

		app.ServeHTTP(w, req)
		response := w.Result()
		body, _ := io.ReadAll(response.Body)
		var responseBody []model.Book
		json.Unmarshal(body, &responseBody)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.NotEmpty(t, responseBody)
		assert.Greater(t, len(responseBody), 0)
		assert.NotEmpty(t, responseBody[0])
		assert.NotEmpty(t, responseBody[0].ID)
		assert.NotEmpty(t, responseBody[0].Name)
		assert.NotEmpty(t, responseBody[0].Description)
	})
}

func TestFindDetailBook(t *testing.T) {
	t.Run("Should retreive book", func(t *testing.T) {
		w := httptest.NewRecorder()
		// assume id 1 always exist because we already seed books on MainFunction
		req, _ := http.NewRequest(http.MethodGet, "/books/1", nil)

		app.ServeHTTP(w, req)
		response := w.Result()
		body, _ := io.ReadAll(response.Body)
		var responseBody model.Book
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.NotEmpty(t, responseBody)
		assert.NotEmpty(t, responseBody)
		assert.NotEmpty(t, responseBody.ID)
		assert.NotEmpty(t, responseBody.Name)
		assert.NotEmpty(t, responseBody.Description)
	})

	t.Run("Should response not found error", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/books/10000", nil)

		app.ServeHTTP(w, req)
		response := w.Result()
		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]string
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusNotFound, response.StatusCode)
		assert.NotEmpty(t, responseBody)
		assert.Equal(t, "Book Not Found", responseBody["Message"])
	})
}

func TestCreateBooks(t *testing.T) {
	t.Run("Should success create books", func(t *testing.T) {
		requestBody := strings.NewReader(`{
			"name": "Example 1",
    	"description": "Lorem Ipsum"
		}`)

		var beforeCount int64
		db.Find(&model.Book{}).Count(&beforeCount)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/books", requestBody)
		req.Header.Set("Content-Type", "application/json")

		app.ServeHTTP(w, req)
		response := w.Result()
		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.Equal(t, "success create book", responseBody["message"])

		var afterCount int64
		db.Find(&model.Book{}).Count(&afterCount)

		assert.Greater(t, afterCount, beforeCount)
	})
}

func TestDelete(t *testing.T) {
	t.Run("should delete book", func(t *testing.T) {
		var beforeCount int64
		db.Find(&model.Book{}).Count(&beforeCount)

		w := httptest.NewRecorder()
		// assume id 1 always exist because we already seed books on MainFunction
		req, _ := http.NewRequest(http.MethodDelete, "/books/1", nil)
		app.ServeHTTP(w, req)

		response := w.Result()
		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]string
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.NotEmpty(t, responseBody)
		assert.Equal(t, "Success delete book", responseBody["message"])

		var afterCount int64
		db.Find(&model.Book{}).Count(&afterCount)

		assert.Greater(t, beforeCount, afterCount)
	})

	t.Run("should response not found", func(t *testing.T) {
		var beforeCount int64
		db.Find(&model.Book{}).Count(&beforeCount)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/books/10000", nil)
		app.ServeHTTP(w, req)

		response := w.Result()
		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]string
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusNotFound, response.StatusCode)
		assert.NotEmpty(t, responseBody)
		assert.Equal(t, "Book Not Found", responseBody["Message"])

		var afterCount int64
		db.Find(&model.Book{}).Count(&afterCount)

		assert.Equal(t, beforeCount, afterCount)
	})
}
