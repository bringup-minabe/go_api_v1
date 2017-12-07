package main

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"runtime"
	"strconv"
)

type User struct {
	gorm.Model
	Uuid        string `json:"uuid"`
	Username    string `json:"username"`
	Password    string `json:"-"`
	UserProfile UserProfile
}

type Res struct {
	Data     []User
	Paginate Paginate
}

/**
 * GetUsers
 */
func GetUsers() echo.HandlerFunc {
	return func(c echo.Context) error {

		var db, _ = DbConnect()

		var paginate Paginate

		// set paginate params
		SetPaginateParams(&paginate, c)

		/**
		 * Select
		 * 並列処理で通常検索とカウント用検索をおこなう
		 */

		// CPU数取得
		cpus := runtime.NumCPU()

		//Find
		uf := make(chan []User, cpus)
		go func() {
			var users []User
			db.
				Limit(paginate.Limit).
				Offset(paginate.Offset).
				Order("id desc").
				Find(&users)

			// Related
			for i := range users {
				db.Model(users[i]).Related(&users[i].UserProfile)
			}

			uf <- users
		}()

		//Count
		uc := make(chan int, cpus)
		go func() {
			count := 0
			db.
				Model(User{}).
				Count(&count)
			uc <- count
		}()

		users := <-uf
		users_count := <-uc
		paginate.Count = users_count

		// set paginate params
		SetPaginateParams(&paginate, c)

		/**
		 * set Responce
		 */
		var res Res
		res.Data = users
		res.Paginate = paginate

		return c.JSON(http.StatusOK, res)
	}
}

/**
 * GetUser
 */
func GetUser() echo.HandlerFunc {
	return func(c echo.Context) error {

		var db, _ = DbConnect()

		id, _ := strconv.Atoi(c.Param("id"))

		/**
		 * Select
		 */
		user := User{}
		if err := db.Where("id = ?", id).First(&user).Error; err != nil {
			return echo.NewHTTPError(http.StatusNotFound)
		} else {
			return c.JSON(http.StatusOK, user)
		}
	}
}

/**
 * AddUser
 */
func AddUser() echo.HandlerFunc {
	return func(c echo.Context) error {

		var db, _ = DbConnect()

		user := User{}

		user.Username = c.FormValue("username")
		user.Uuid = uuid.NewV4().String()

		//パスワード
		password := c.FormValue("password")
		hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		user.Password = string(hash)

		// add user profile
		user_profile := UserProfile{
			LastName:  c.FormValue("last_name"),
			FirstName: c.FormValue("first_name"),
		}
		user.UserProfile = user_profile

		// db.NewRecord(user)
		db.Create(&user)

		return c.JSON(http.StatusOK, user)
	}
}

/**
 * EditUser
 */
func EditUser() echo.HandlerFunc {
	return func(c echo.Context) error {

		var db, _ = DbConnect()

		id, _ := strconv.Atoi(c.Param("id"))

		user := User{}
		if err := db.Where("id = ?", id).First(&user).Error; err != nil {
			return echo.NewHTTPError(http.StatusNotFound)
		}

		user.Username = c.FormValue("username")

		//パスワード
		password := c.FormValue("password")
		if password != "" {
			hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			user.Password = string(hash)
		}

		// edit user profile
		user_profile := UserProfile{
			LastName:  c.FormValue("last_name"),
			FirstName: c.FormValue("first_name"),
		}
		user.UserProfile = user_profile

		db.Save(&user)

		return c.JSON(http.StatusOK, user)
	}
}

/**
 * DeleteUser
 */
func DeleteUser() echo.HandlerFunc {
	return func(c echo.Context) error {

		var db, _ = DbConnect()

		id, _ := strconv.Atoi(c.Param("id"))

		user := User{}
		if err := db.Where("id = ?", id).First(&user).Error; err != nil {
			return echo.NewHTTPError(http.StatusNotFound)
		}

		db.Delete(&user)

		return c.JSON(http.StatusOK, user)
	}
}
