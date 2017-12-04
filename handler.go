package main

import (
	"github.com/labstack/echo"
	"net/http"
)

func MainPage() echo.HandlerFunc {
	return func(c echo.Context) error {

		jsonMap := map[string]string{
			"foo":  "bar",
			"hoge": "fuga",
		}

		return c.JSON(http.StatusOK, jsonMap)
	}
}
