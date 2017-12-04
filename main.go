package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {

	//Database Conect
	var Db, _ = DbConnect()

	// close Db
	defer Db.Close()

	// create echo instance
	e := echo.New()

	// API Path
	api_path := "api/v1"

	/**
	 * set Use
	 */
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost/",
			"https://localhost/",
		},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	/**
	 * routing
	 */

	//top
	e.GET("/", MainPage())

	//users
	e.GET(api_path+"/users/", GetUsers())
	e.GET(api_path+"/users/:id", GetUser())
	e.POST(api_path+"/users/", AddUser())

	/**
	 * start
	 */
	e.Logger.Fatal(e.Start(":1323"))

}
