package controller

import (
	"encoding/json"
	"errors"
	"example/model"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	e *echo.Echo
)

func TestMain(m *testing.M) {
	e = echo.New()

	m.Run()
}

func TestGetAllBooks(t *testing.T) {
	t.Run("should response OK", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/books", nil)
		w := httptest.NewRecorder()
		c := e.NewContext(req, w)

		repo := BookMockRepository{}
		repo.Mock.On("FindAll").Return([]model.Book{
			{ID: 100, Name: "Book 100", Description: "Lorem 100"},
			{ID: 101, Name: "Book 101", Description: "Lorem 101"},
			{ID: 102, Name: "Book 102", Description: "Lorem 102"},
		}, nil)

		ctr := BookController{&repo}
		ctr.Index(c)

		response := w.Result()
		body, _ := io.ReadAll(response.Body)
		var responseBody []model.Book
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, 3, len(responseBody))
		assert.NotEmpty(t, responseBody[0])
		assert.NotEmpty(t, responseBody[0].ID)
		assert.NotEmpty(t, responseBody[0].Name)
		assert.NotEmpty(t, responseBody[0].Description)
	})
}

func TestFindDetailBook(t *testing.T) {
	t.Run("Should retreive book", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/books/1", nil)
		w := httptest.NewRecorder()
		c := e.NewContext(req, w)

		c.SetParamNames("id")
		c.SetParamValues("1")

		repo := BookMockRepository{}
		repo.Mock.On("FindById", 1).Return(model.Book{
			ID:          1,
			Name:        "Book 100",
			Description: "Lorem 100",
		}, nil)

		ctr := BookController{&repo}
		ctr.Detail(c)

		response := w.Result()
		body, _ := io.ReadAll(response.Body)
		var responseBody model.Book
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, responseBody)
		assert.Equal(t, responseBody.ID, 1)
		assert.Equal(t, responseBody.Name, "Book 100")
		assert.Equal(t, responseBody.Description, "Lorem 100")
	})

	t.Run("Should not retreive book", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/books/2", nil)
		w := httptest.NewRecorder()
		c := e.NewContext(req, w)

		c.SetParamNames("id")
		c.SetParamValues("2")

		repo := BookMockRepository{}
		repo.Mock.On("FindById", 2).Return(model.Book{}, errors.New("book notx found"))

		ctr := BookController{&repo}
		ctr.Detail(c)

		response := w.Result()
		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]string
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.NotEmpty(t, responseBody)
		assert.Equal(t, responseBody["message"], "book notx found")
	})
}

func TestCreate(t *testing.T) {
	t.Run("should create book", func(t *testing.T) {
		requestBody := strings.NewReader(`{
			"name": "Example 1",
    	"description": "Lorem Ipsum"
		}`)

		book := &model.Book{Name: "Example 1", Description: "Lorem Ipsum"}

		req := httptest.NewRequest(http.MethodGet, "/books", requestBody)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c := e.NewContext(req, w)

		repo := BookMockRepository{}
		repo.Mock.On("Create", book).Return(nil)

		ctr := BookController{&repo}
		ctr.Create(c)

		response := w.Result()
		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.NotEmpty(t, responseBody)
		assert.Equal(t, responseBody["message"], "success create book")
	})

	t.Run("should not create book", func(t *testing.T) {
		requestBody := strings.NewReader(`{
			"name": "Example 1",
    	"description": "Lorem Ipsum"
		}`)

		book := &model.Book{Name: "Example 1", Description: "Lorem Ipsum"}

		req := httptest.NewRequest(http.MethodGet, "/books", requestBody)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c := e.NewContext(req, w)

		repo := BookMockRepository{}
		repo.Mock.On("Create", book).Return(errors.New("failed create book"))

		ctr := BookController{&repo}
		ctr.Create(c)

		response := w.Result()
		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.NotEmpty(t, responseBody)
		assert.Equal(t, responseBody["message"], "failed create book")
	})
}
