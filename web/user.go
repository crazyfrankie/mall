package web

import (
	"errors"
	"net/http"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"

	"mall/domain"
	"mall/service"
)

type UserHandler struct {
	svc               service.UserService
	emailRegexpExp    *regexp.Regexp
	passwordRegexpExp *regexp.Regexp
	phoneRegexExp     *regexp.Regexp
}

func NewUserHandler(svc service.UserService) *UserHandler {
	const (
		passwordRegexPattern = `^(?=.*[a-zA-Z])(?=.*\d)(?=.*[$@$!%*#?&])[a-zA-Z\d$@$!%*#?&]{8,}$`
		phoneRegexPattern    = `^1[3-9]\d{9}$`
	)
	passwordRegexExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	phoneRegexExp := regexp.MustCompile(phoneRegexPattern, regexp.None)

	return &UserHandler{
		svc:               svc,
		passwordRegexpExp: passwordRegexExp,
		phoneRegexExp:     phoneRegexExp,
	}
}

func (ctl *UserHandler) RegisterRoute(r *gin.Engine) {
	userGroup := r.Group("user")
	{
		userGroup.POST("signup", ctl.SignUp())
	}
}

func (ctl *UserHandler) SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		type Req struct {
			Name            string `json:"name"`
			Password        string `json:"password"`
			ConfirmPassword string `json:"confirmPassword"`
		}
		var req Req
		if err := c.Bind(&req); err != nil {
			return
		}

		// 两次密码不一致
		if req.Password != req.ConfirmPassword {
			c.JSON(http.StatusBadRequest, GetResponse(WithStatus(http.StatusBadRequest), WithMsg("Two passwords do not match")))
			return
		}

		// 密码格式不正确
		ok, err := ctl.passwordRegexpExp.MatchString(req.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, GetResponse(WithStatus(http.StatusBadRequest), WithMsg("system error")))
			return
		}
		if !ok {
			c.JSON(http.StatusBadRequest, GetResponse(WithStatus(http.StatusBadRequest), WithMsg("password format is incorrect")))
			return
		}

		err = ctl.svc.SignUp(c.Request.Context(), domain.User{
			Name:     req.Name,
			Password: req.Password,
		})
		switch {
		case errors.Is(err, service.ErrUserDuplicateName):
			c.JSON(http.StatusBadRequest, GetResponse(WithStatus(http.StatusBadRequest), WithMsg("duplicate user name")))
			return
		case err != nil:
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		default:
			c.JSON(http.StatusOK, GetResponse(WithStatus(http.StatusOK), WithMsg("sign up successfully")))
		}
	}
}

func (ctl *UserHandler) LogIn() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
