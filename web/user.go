package web

import (
	"errors"
	"log"
	"mall/domain"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"mall/middleware/jwt"
	"mall/service"
)

type UserHandler struct {
	userSvc *service.UserService
	codeSvc *service.CodeService
	jwtHdl  *jwt.TokenHandler
}

func NewUserHandler(userSvc *service.UserService, codeSvc *service.CodeService, jwtHdl *jwt.TokenHandler, sessHdl *jwt.RedisSession) *UserHandler {
	return &UserHandler{
		userSvc: userSvc,
		codeSvc: codeSvc,
		jwtHdl:  jwtHdl,
	}
}

func (ctl *UserHandler) RegisterRoute(r *gin.Engine) {
	userGroup := r.Group("api/user")
	{
		userGroup.POST("signup", ctl.PreSignupCheck())
		userGroup.POST("login", ctl.NameLogin())
		userGroup.POST("send-code", ctl.SendVerificationCode())
		userGroup.POST("signup/verify-code", ctl.SignupVerifyCode())
		userGroup.POST("login/verify-code", ctl.LoginVerifyCode())
		userGroup.POST("bind", ctl.BindInfo())
	}
}

func (ctl *UserHandler) PreSignupCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		type Req struct {
			Phone string `json:"phone"`
		}
		var req Req
		if err := c.Bind(&req); err != nil {
			return
		}

		log.Println("Checking phone:", req.Phone)
		err := ctl.userSvc.CheckPhone(c.Request.Context(), req.Phone)
		if err == nil {
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}

		c.JSON(http.StatusAccepted, GetResponse(WithStatus(http.StatusAccepted), WithMsg("phone need verification")))
	}
}

func (ctl *UserHandler) SendVerificationCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		type Req struct {
			Phone string `json:"phone" validate:"required,len=11"`
			Biz   string `json:"biz"`
		}

		var req Req
		if err := c.Bind(&req); err != nil {
			return
		}

		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			c.JSON(http.StatusBadRequest, GetResponse(WithStatus(http.StatusBadRequest), WithMsg("Validation failed:"+err.Error()), WithErr(err.Error())))
			return
		}

		err := ctl.codeSvc.Send(c.Request.Context(), req.Biz, req.Phone)
		switch {
		case errors.Is(err, service.ErrSendTooMany):
			c.JSON(http.StatusBadRequest, GetResponse(WithStatus(http.StatusBadRequest), WithMsg("send too many")))
			return
		case err != nil:
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		default:
			c.JSON(http.StatusOK, GetResponse(WithStatus(http.StatusOK), WithMsg("send successfully")))
		}
	}
}

func (ctl *UserHandler) SignupVerifyCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		type Req struct {
			Phone string `json:"phone" validate:"required,len=11"`
			Code  string `json:"code"`
		}

		var req Req
		if err := c.Bind(&req); err != nil {
			return
		}

		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			c.JSON(http.StatusBadRequest, GetResponse(WithStatus(http.StatusBadRequest), WithMsg("Validation failed"+err.Error()), WithErr(err.Error())))
		}

		_, err := ctl.codeSvc.Verify(c.Request.Context(), "signup", req.Phone, req.Code)
		switch {
		case errors.Is(err, service.ErrVerifyTooMany):
			c.JSON(http.StatusBadRequest, GetResponse(WithStatus(http.StatusBadRequest), WithMsg("verify too many")))
			return
		case err != nil:
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}

		err = ctl.userSvc.CreateUser(c.Request.Context(), req.Phone)
		switch {
		case errors.Is(err, service.ErrUserDuplicatePhone):
			c.JSON(http.StatusBadRequest, GetResponse(WithStatus(http.StatusBadRequest), WithMsg("duplicate phone")))
			return
		case err != nil:
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}

		c.JSON(http.StatusOK, GetResponse(WithStatus(http.StatusOK), WithMsg("signup successfully")))
	}
}

func (ctl *UserHandler) LoginVerifyCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		type Req struct {
			Phone string `json:"phone" validate:"required,len=11"`
			Code  string `json:"code"`
		}

		var req Req
		if err := c.Bind(&req); err != nil {
			return
		}

		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			c.JSON(http.StatusBadRequest, GetResponse(WithStatus(http.StatusBadRequest), WithMsg("Validation failed"+err.Error()), WithErr(err.Error())))
		}

		_, err := ctl.codeSvc.Verify(c.Request.Context(), "login", req.Phone, req.Code)
		switch {
		case errors.Is(err, service.ErrVerifyTooMany):
			c.JSON(http.StatusBadRequest, GetResponse(WithStatus(http.StatusBadRequest), WithMsg("verify too many")))
			return
		case err != nil:
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}

		// 设置 JWT
		var ssid string
		ssid, err = ctl.userSvc.SetSession(c.Request.Context(), req.Phone)
		if err != nil {
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}
		err = ctl.jwtHdl.GenerateToken(c, ssid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}

		c.JSON(http.StatusOK, GetResponse(WithStatus(http.StatusOK), WithMsg("login successfully")))
	}
}

func (ctl *UserHandler) BindInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		type Req struct {
			Phone           string    `json:"phone"`
			Name            string    `json:"name" validate:"required,min=3,max=20"`
			Password        string    `json:"password" validate:"required,min=8,containsany=abcdefghijklmnopqrstuvwxyz,containsany=0123456789,containsany=$@$!%*#?&"`
			ConfirmPassword string    `json:"confirmPassword" validate:"eqfield=Password"`
			Birthday        time.Time `json:"birthday"`
		}
		var req Req
		if err := c.Bind(&req); err != nil {
			return
		}

		// 使用 validator 进行字段验证
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			c.JSON(http.StatusBadRequest, GetResponse(WithStatus(http.StatusBadRequest), WithMsg("Validation failed: "+err.Error()), WithErr(err.Error())))
			return
		}

		err := ctl.userSvc.BindInfo(c.Request.Context(), domain.User{
			Phone:    req.Phone,
			Name:     req.Name,
			Password: req.Password,
			Birthday: req.Birthday,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}

		c.JSON(http.StatusOK, GetResponse(WithStatus(http.StatusOK), WithMsg("bind successfully")))
	}
}

func (ctl *UserHandler) NameLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		type Req struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}

		var req Req
		if err := c.Bind(&req); err != nil {
			return
		}

		phone, err := ctl.userSvc.NameLogin(c.Request.Context(), domain.User{
			Name:     req.Name,
			Password: req.Password,
		})

		switch {
		case errors.Is(err, service.ErrRecordNotFound):
			c.JSON(http.StatusBadRequest, GetResponse(WithStatus(http.StatusBadRequest), WithMsg("user not found")))
			return
		case errors.Is(err, service.ErrInvalidUserOrPassword):
			c.JSON(http.StatusBadRequest, GetResponse(WithStatus(http.StatusBadRequest), WithMsg("username or password error")))
			return
		case err != nil:
			c.JSON(http.StatusBadRequest, GetResponse(WithStatus(http.StatusBadRequest), WithMsg("system error")))
			return
		}

		var ssid string
		ssid, err = ctl.userSvc.SetSession(c.Request.Context(), phone)
		err = ctl.jwtHdl.GenerateToken(c, ssid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}

		c.JSON(http.StatusOK, GetResponse(WithStatus(http.StatusOK), WithMsg("login successfully")))
	}
}

func (ctl *UserHandler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := ctl.jwtHdl.ExtractToken(c)

		claim, err := ctl.jwtHdl.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusBadRequest, GetResponse(WithStatus(http.StatusBadRequest), WithMsg(err.Error()), WithErr(err.Error())))
			return
		}

		err = ctl.userSvc.DeleteSession(c.Request.Context(), claim.SessionId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}

		c.JSON(http.StatusOK, GetResponse(WithStatus(http.StatusOK), WithMsg("log out successfully")))
	}
}

func (ctl *UserHandler) EditInfo() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (ctl *UserHandler) Heartbeat() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
