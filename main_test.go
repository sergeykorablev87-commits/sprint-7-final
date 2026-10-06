package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	city := "moscow"
	requests := []struct {
		count int
		want  int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, min(100, len(cafeList[city]))},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		request := fmt.Sprintf("/cafe?count=%d&city=%s", v.count, city)
		req := httptest.NewRequest("GET", request, nil)
		handler.ServeHTTP(response, req)
		require.Equal(t, http.StatusOK, response.Code)
		body := strings.TrimSpace(response.Body.String())
		if body == "" {
			assert.Equal(t, v.want, 0)
		} else {
			assert.Len(t, strings.Split(body, ","), v.want)
		}
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	city := "moscow"
	requests := []struct {
		search    string
		wantCount int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		request := fmt.Sprintf("/cafe?city=%s&search=%s", city, v.search)
		req := httptest.NewRequest("GET", request, nil)
		handler.ServeHTTP(response, req)
		require.Equal(t, http.StatusOK, response.Code)
		body := strings.TrimSpace(response.Body.String())
		if body == "" {
			assert.Equal(t, v.wantCount, 0)
		} else {
			cafes := strings.Split(body, ",")
			assert.Len(t, cafes, v.wantCount)
			for _, cafe := range cafes {
				assert.True(t, strings.Contains(strings.ToLower(cafe), strings.ToLower(v.search)))
			}
		}
	}
}
