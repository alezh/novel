package basics

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"github.com/kataras/iris/mvc"
	"strings"
)

type AuthController struct {
	Ctx iris.Context
	Session *sessions.Session
	UserID string
}

var (
	PathLogin  = mvc.Response{Path: "/Admin/login"}
	PathLogout = mvc.Response{Path: "/Admin/logout"}
)



func (c *AuthController) BeginRequest(ctx iris.Context) {
	//if "/Admin/login" != ctx.GetCurrentRoute().Path(){
	//	c.UserID = c.Session.GetString(config.SessionIDKey)
	//	if c.UserID == ""{
	//		c.redirect()
	//	}
	//}
}

func (c *AuthController) EndRequest(ctx iris.Context) {}

func (c *AuthController) fireError(err error) mvc.View {
	return mvc.View{
		Code: iris.StatusBadRequest,
		Name: "shared/error.html",
		Data: iris.Map{"Title": "User Error", "Message": strings.ToUpper(err.Error())},
	}
}

func (c *AuthController) redirect() mvc.Response {
	return mvc.Response{Path: "/Admin/login"}
}

//func (c *AuthController) createOrUpdate(firstname, username, password string) (user Model, err error) {
//	username = strings.Trim(username, " ")
//	if username == "" || password == "" || firstname == "" {
//		return user, errors.New("empty firstname, username or/and password")
//	}
//
//	userToInsert := Model{
//		Firstname: firstname,
//		Username:  username,
//		password:  password,
//	} // password is hashed by the Source.
//
//	newUser, err := c.Source.InsertOrUpdate(userToInsert)
//	if err != nil {
//		return user, err
//	}
//
//	return newUser, nil
//}

func (c *AuthController) isLoggedIn() bool {
	// we don't search by session, we have the user id
	// already by the `BeginRequest` middleware.
	return c.UserID != ""
}

//func (c *AuthController) verify(username, password string) (user Model, err error) {
//
//	if username == "" || password == "" {
//		return user, errors.New("please fill both username and password fields")
//	}
//
//	u, found := c.Source.GetByUsername(username)
//	if !found {
//		// if user found with that username not found at all.
//		return user, errors.New("user with that username does not exist")
//	}
//
//	if ok, err := ValidatePassword(password, u.HashedPassword); err != nil || !ok {
//		// if user found but an error occurred or the password is not valid.
//		return user, errors.New("please try to login with valid credentials")
//	}
//
//	return u, nil
//}

// if logged in then destroy the session
// and redirect to the login page
// otherwise redirect to the registration page.
func (c *AuthController) logout() mvc.Response {
	if c.isLoggedIn() {
		c.Session.Destroy()
	}
	return PathLogin
}