package web

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"mall/internal/user/domain"
	"mall/internal/user/middleware/jwt"
	"mall/internal/user/service"
)

type UserHandler struct {
	userSvc *service.UserService
	codeSvc *service.CodeService
	jwtHdl  *jwt.TokenHandler
}

func NewUserHandler(userSvc *service.UserService, codeSvc *service.CodeService, jwtHdl *jwt.TokenHandler) *UserHandler {
	return &UserHandler{
		userSvc: userSvc,
		codeSvc: codeSvc,
		jwtHdl:  jwtHdl,
	}
}

func (ctl *UserHandler) RegisterRoute(r *gin.Engine) {
	userGroup := r.Group("api/user")
	{
		userGroup.POST("signup", ctl.PreCheckPhone())
		userGroup.POST("login", ctl.NameLogin())
		userGroup.POST("send-code", ctl.SendVerificationCode())
		userGroup.POST("verify-code", ctl.VerificationCode())
		userGroup.POST("bind/password", ctl.UpdatePassword())
		userGroup.POST("bind/name", ctl.UpdateName())
		userGroup.POST("bind/birthday", ctl.UpdateBirthday())
		userGroup.GET("refresh", ctl.KeepAive())
	}
}

func (ctl *UserHandler) PreCheckPhone() gin.HandlerFunc {
	return func(c *gin.Context) {
		type Req struct {
			Phone string `json:"phone"`
		}
		var req Req
		if err := c.Bind(&req); err != nil {
			zap.L().Error("预检手机号:绑定信息错误", zap.Error(err))
			return
		}

		err := ctl.userSvc.CheckPhone(c.Request.Context(), req.Phone)
		switch {
		case errors.Is(err, service.ErrRecordNotFound):
			zap.L().Info("预检手机号不存在:需要注册", zap.Error(err))
			c.JSON(http.StatusAccepted, GetResponse(WithStatus(http.StatusAccepted), WithMsg("phone need signup")))
			return
		case err != nil:
			zap.L().Error("预检手机号:系统错误", zap.Error(err))
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}

		zap.L().Info("手机号预检成功")
		c.JSON(http.StatusOK, GetResponse(WithStatus(http.StatusOK), WithMsg("phone need login")))
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
			zap.L().Error("发送验证码:绑定信息错误", zap.Error(err))
			return
		}

		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			zap.L().Error("发送验证码:手机号格式校验错误", zap.Error(err))
			c.JSON(http.StatusBadRequest, GetResponse(WithStatus(http.StatusBadRequest), WithMsg("Validation failed:"+err.Error())))
			return
		}

		err := ctl.codeSvc.Send(c.Request.Context(), req.Biz, req.Phone)
		switch {
		case errors.Is(err, service.ErrSendTooMany):
			zap.L().Error("发送验证码:发送过于频繁", zap.Error(err))
			c.JSON(http.StatusTooManyRequests, GetResponse(WithStatus(http.StatusTooManyRequests), WithMsg("send too many")))
			return
		case err != nil:
			zap.L().Error("发送验证码:系统错误", zap.Error(err))
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		default:
			zap.L().Info("发送验证码成功")
			c.JSON(http.StatusOK, GetResponse(WithStatus(http.StatusOK), WithMsg("send successfully")))
		}
	}
}

func (ctl *UserHandler) VerificationCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		type Req struct {
			Phone string `json:"phone" validate:"required,len=11"`
			Code  string `json:"code"`
			Biz   string `json:"biz"`
		}

		var req Req
		if err := c.Bind(&req); err != nil {
			zap.L().Error("校验验证码绑定信息错误", zap.Error(err))
			return
		}

		_, err := ctl.codeSvc.Verify(c.Request.Context(), req.Biz, req.Phone, req.Code)
		switch {
		case errors.Is(err, service.ErrVerifyTooMany):
			zap.L().Error(fmt.Sprintf("%s:校验验证码:校验次数过多", req.Biz), zap.Error(err))
			c.JSON(http.StatusTooManyRequests, GetResponse(WithStatus(http.StatusTooManyRequests), WithMsg("verify too many")))
			return
		case err != nil:
			zap.L().Error(fmt.Sprintf("%s:校验验证码:系统错误", req.Biz), zap.Error(err))
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}

		var user domain.User
		user, err = ctl.userSvc.FindOrCreateUser(c.Request.Context(), req.Phone)
		if err != nil {
			zap.L().Error(fmt.Sprintf("%s:查找或创建用户:系统错误", req.Biz), zap.Error(err))
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}

		// 设置 JWT
		var ssid string
		ssid, err = ctl.userSvc.SetSession(c.Request.Context(), req.Phone)
		if err != nil {
			zap.L().Error(fmt.Sprintf("手机号%s:创建 Session:系统错误", req.Biz), zap.Error(err))
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}
		err = ctl.jwtHdl.GenerateToken(c, ssid)
		if err != nil {
			zap.L().Error(fmt.Sprintf("手机号%s:设置 JWT:系统错误", req.Biz), zap.Error(err))
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}

		maskedPhone := req.Phone[:3] + "****" + req.Phone[len(req.Phone)-4:]
		zap.L().Info(fmt.Sprintf("%s:用户处理成功", req.Biz), zap.String("phone", maskedPhone))

		c.JSON(http.StatusOK, GetResponse(WithStatus(http.StatusOK), WithMsg(fmt.Sprintf("%s successfully", req.Biz)), WithData(map[string]interface{}{
			"id":    user.Id,
			"phone": user.Phone,
			"name":  user.Name,
		})))
	}
}

func (ctl *UserHandler) UpdatePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		type Req struct {
			UserId          uint64 `json:"user_id"`
			Password        string `json:"password" validate:"required,min=8,containsany=abcdefghijklmnopqrstuvwxyz,containsany=0123456789,containsany=$@$!%*#?&"`
			ConfirmPassword string `json:"confirmPassword" validate:"eqfield=Password"`
		}
		var req Req
		if err := c.Bind(&req); err != nil {
			zap.L().Error("绑定用户密码:绑定信息错误", zap.Error(err))
			return
		}

		// 使用 validator 进行字段验证
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			zap.L().Error("绑定用户信息:信息格式错误", zap.Error(err))
			c.JSON(http.StatusBadRequest, GetResponse(WithStatus(http.StatusBadRequest), WithMsg("Validation failed: "+err.Error())))
			return
		}

		err := ctl.userSvc.UpdatePassword(c.Request.Context(), domain.User{
			Id:       req.UserId,
			Password: req.Password,
		})
		if err != nil {
			zap.L().Error("绑定用户密码:系统错误", zap.Error(err))
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}

		zap.L().Info("绑定用户密码成功")

		c.JSON(http.StatusOK, GetResponse(WithStatus(http.StatusOK), WithMsg("bind user's password successfully")))
	}
}

func (ctl *UserHandler) UpdateBirthday() gin.HandlerFunc {
	return func(c *gin.Context) {
		type Req struct {
			UserId   uint64    `json:"user_id"`
			Birthday time.Time `json:"birthday"`
		}

		var req Req
		if err := c.Bind(&req); err != nil {
			zap.L().Error("绑定用户生日:绑定信息错误", zap.Error(err))
			return
		}

		// 使用 validator 进行字段验证
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			zap.L().Error("绑定用户生日:信息格式错误", zap.Error(err))
			c.JSON(http.StatusBadRequest, GetResponse(WithStatus(http.StatusBadRequest), WithMsg("Validation failed: "+err.Error())))
			return
		}

		err := ctl.userSvc.UpdateBirthday(c.Request.Context(), domain.User{
			Id:       req.UserId,
			Birthday: req.Birthday,
		})
		if err != nil {
			zap.L().Error("绑定用户生日:系统错误", zap.Error(err))
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}

		zap.L().Info("绑定用户生日成功")

		c.JSON(http.StatusOK, GetResponse(WithStatus(http.StatusOK), WithMsg("bind user's birthday successfully")))
	}
}

func (ctl *UserHandler) UpdateName() gin.HandlerFunc {
	return func(c *gin.Context) {
		type Req struct {
			UserId uint64 `json:"user_id"`
			Name   string `json:"name" validate:"required,min=3,max=20"`
		}

		var req Req
		if err := c.Bind(&req); err != nil {
			zap.L().Error("绑定用户名:绑定信息错误", zap.Error(err))
			return
		}

		// 使用 validator 进行字段验证
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			zap.L().Error("绑定用户名:信息格式错误", zap.Error(err))
			c.JSON(http.StatusBadRequest, GetResponse(WithStatus(http.StatusBadRequest), WithMsg("Validation failed: "+err.Error())))
			return
		}

		err := ctl.userSvc.UpdateName(c.Request.Context(), domain.User{
			Id:   req.UserId,
			Name: req.Name,
		})
		switch {
		case errors.Is(err, service.ErrUserDuplicateName):
			zap.L().Error("绑定用户名:用户名冲突", zap.Error(err))
			c.JSON(http.StatusConflict, GetResponse(WithStatus(http.StatusConflict), WithMsg("duplicate name")))
			return
		case err != nil:
			zap.L().Error("绑定用户名:系统错误", zap.Error(err))
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}

		zap.L().Info("绑定用户名成功")

		c.JSON(http.StatusOK, GetResponse(WithStatus(http.StatusOK), WithMsg("bind user's name successfully")))
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
			zap.L().Error("用户名登录:绑定信息错误", zap.Error(err))
			return
		}

		phone, err := ctl.userSvc.NameLogin(c.Request.Context(), domain.User{
			Name:     req.Name,
			Password: req.Password,
		})

		switch {
		case errors.Is(err, service.ErrRecordNotFound):
			zap.L().Error("用户名登录:用户不存在", zap.Error(err))
			c.JSON(http.StatusNotFound, GetResponse(WithStatus(http.StatusNotFound), WithMsg("user not found")))
			return
		case errors.Is(err, service.ErrInvalidUserOrPassword):
			zap.L().Error("用户名登录:用户名或密码不正确", zap.Error(err))
			c.JSON(http.StatusUnauthorized, GetResponse(WithStatus(http.StatusUnauthorized), WithMsg("username or password error")))
			return
		case err != nil:
			zap.L().Error("用户名登录:系统错误", zap.Error(err))
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}

		var ssid string
		ssid, err = ctl.userSvc.SetSession(c.Request.Context(), phone)
		if err != nil {
			zap.L().Error("登录:创建 Session:系统错误", zap.Error(err))
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}
		err = ctl.jwtHdl.GenerateToken(c, ssid)
		if err != nil {
			zap.L().Error("登录:设置 JWT:系统错误", zap.Error(err))
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}

		zap.L().Info("用户登录成功")
		c.JSON(http.StatusOK, GetResponse(WithStatus(http.StatusOK), WithMsg("login successfully")))
	}
}

func (ctl *UserHandler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := ctl.jwtHdl.ExtractToken(c)

		claim, err := ctl.jwtHdl.ParseToken(token)
		if err != nil {
			zap.L().Error("退出登录:解析 Token 错误", zap.Error(err))
			c.JSON(http.StatusBadRequest, GetResponse(WithStatus(http.StatusBadRequest), WithMsg(err.Error())))
			return
		}

		err = ctl.userSvc.DeleteSession(c.Request.Context(), claim.SessionId)
		if err != nil {
			zap.L().Error("退出登录:系统错误", zap.Error(err))
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}

		zap.L().Info("用户退出登录成功")
		c.JSON(http.StatusOK, GetResponse(WithStatus(http.StatusOK), WithMsg("log out successfully")))
	}
}

func (ctl *UserHandler) KeepAive() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHeader := ctl.jwtHdl.ExtractToken(c)

		claims, err := ctl.jwtHdl.ParseToken(tokenHeader)
		if err != nil {
			zap.L().Error("维持登录状态:解析 Token 错误")
			c.JSON(http.StatusBadRequest, GetResponse(WithStatus(http.StatusBadRequest), WithMsg("unauthorized")))
			return
		}

		// 刷新 Session 有效期
		err = ctl.userSvc.ExtendSessionExpiration(c.Request.Context(), claims.SessionId)
		if err != nil {
			zap.L().Error("维持登录状态:系统错误")
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}

		zap.L().Info("维持登录状态成功")
		c.JSON(http.StatusOK, GetResponse(WithStatus(http.StatusOK), WithMsg("session refreshed")))
	}
}

func (ctl *UserHandler) BindAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		type Req struct {
			UserId    uint64 `json:"userId"`
			Street    string `json:"street"`
			City      string `json:"city"`
			State     string `json:"state"`
			ZipCode   string `json:"zipCode"`
			Country   string `json:"country"`
			IsDefault bool   `json:"isDefault"`
		}

		var req Req
		if err := c.Bind(&req); err != nil {
			zap.L().Error("添加地址:绑定信息错误", zap.Error(err))
			return
		}

		err := ctl.userSvc.BindAddress(c.Request.Context(), domain.Address{
			UserId:    req.UserId,
			Street:    req.Street,
			State:     req.State,
			City:      req.City,
			ZipCode:   req.ZipCode,
			Country:   req.Country,
			IsDefault: req.IsDefault,
		})
		if err != nil {
			zap.L().Error("添加地址:系统错误")
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}

		zap.L().Info("添加地址成功")
		c.JSON(http.StatusOK, GetResponse(WithStatus(http.StatusOK), WithMsg("add addr successfully")))
	}
}

func (ctl *UserHandler) AcquireAllAddr() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		addresses, err := ctl.userSvc.AcquireAllAddr(c.Request.Context(), id)
		if err != nil {
			zap.L().Error("查询所有地址:系统错误")
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}

		zap.L().Info("查询所有地址成功")
		c.JSON(http.StatusOK, GetResponse(WithStatus(http.StatusOK), WithData(addresses)))
	}
}

func (ctl *UserHandler) DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		err := ctl.userSvc.DeleteAddress(c.Request.Context(), id)
		if err != nil {
			zap.L().Error("删除地址:系统错误")
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}

		zap.L().Info("删除地址成功")
		c.JSON(http.StatusOK, GetResponse(WithStatus(http.StatusOK), WithMsg("delete successfully")))
	}
}

func (ctl *UserHandler) UpdateAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		type Req struct {
			UserId    uint64 `json:"userId"`
			Street    string `json:"street"`
			City      string `json:"city"`
			State     string `json:"state"`
			ZipCode   string `json:"zipCode"`
			Country   string `json:"country"`
			IsDefault bool   `json:"isDefault"`
		}

		var req Req
		if err := c.Bind(&req); err != nil {
			zap.L().Error("更新地址:绑定信息错误", zap.Error(err))
			return
		}

		addr, err := ctl.userSvc.UpdateAddress(c.Request.Context(), domain.Address{
			UserId:    req.UserId,
			Street:    req.Street,
			State:     req.State,
			City:      req.City,
			ZipCode:   req.ZipCode,
			Country:   req.Country,
			IsDefault: req.IsDefault,
		})
		if err != nil {
			zap.L().Error("更新地址:系统错误", zap.Error(err))
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")))
			return
		}

		zap.L().Info("更新地址成功")
		c.JSON(http.StatusOK, GetResponse(WithStatus(http.StatusOK), WithData(addr)))
	}
}
