package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/go-playground/validator/v10"
	"mall/internal/auth/jwt"
	"mall/internal/user/domain"
	"mall/internal/user/service"
	"mall/pkg/logger"
)

type UserHandler struct {
	userSvc *service.UserService
	codeSvc *service.CodeService
	jwtHdl  *jwt.TokenHandler
	l       logger.Logger
}

func NewUserHandler(userSvc *service.UserService, codeSvc *service.CodeService, jwtHdl *jwt.TokenHandler, l logger.Logger) *UserHandler {
	return &UserHandler{
		userSvc: userSvc,
		codeSvc: codeSvc,
		jwtHdl:  jwtHdl,
		l:       l,
	}
}

func (ctl *UserHandler) RegisterRoute(r *server.Hertz) {
	userGroup := r.Group("api/user")
	{
		userGroup.POST("login", ctl.NameLogin())
		userGroup.POST("send-code", ctl.SendVerificationCode())
		userGroup.POST("verify-code", ctl.VerificationCode())
		userGroup.POST("bind/password", ctl.UpdatePassword())
		userGroup.POST("bind/name", ctl.UpdateName())
		userGroup.POST("bind/birthday", ctl.UpdateBirthday())
		userGroup.GET("refresh", ctl.KeepAive())
	}
}

func (ctl *UserHandler) SendVerificationCode() app.HandlerFunc {
	return WrapReq(func(ctx context.Context, c *app.RequestContext, req struct {
		Phone string `json:"phone" validate:"required,len=11"`
		Biz   string `json:"biz"`
	}) (Response, error) {
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			ctl.l.Error("发送验证码:校验失败")
			return Response{}, err
		}

		err := ctl.codeSvc.Send(ctx, req.Biz, req.Phone)
		if err != nil {
			return Response{}, err
		}

		// 如果发送成功，返回成功的响应
		ctl.l.Info("发送验证码成功")
		return GetResponse(WithStatus(http.StatusOK), WithMsg("send successfully")), nil
	}, func(c *app.RequestContext, err error) (Response, bool) {
		switch {
		case errors.Is(err, service.ErrSendTooMany):
			ctl.l.Error("发送验证码:发送过于频繁")
			return GetResponse(WithStatus(http.StatusTooManyRequests), WithMsg("send too many")), true
		default:
			// 处理未处理的错误
			ctl.l.Error("发送验证码:系统错误", logger.Error("error", err))
			return GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")), false
		}
	})
}

func (ctl *UserHandler) VerificationCode() app.HandlerFunc {
	return WrapReq(func(ctx context.Context, c *app.RequestContext, req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
		Biz   string `json:"biz"`
	}) (Response, error) {
		// 校验验证码
		_, err := ctl.codeSvc.Verify(ctx, req.Biz, req.Phone, req.Code)
		if err != nil {
			ctl.l.Error(fmt.Sprintf("%s:验证码校验失败", req.Biz), logger.String("phone", req.Phone), logger.Error("error", err))
			return Response{}, NewBusinessError("failed to verify code", err)
		}

		// 查找或创建用户
		user, err := ctl.userSvc.FindOrCreateUser(ctx, req.Phone)
		if err != nil {
			ctl.l.Error(fmt.Sprintf("%s:查找或创建用户失败", req.Biz), logger.String("phone", req.Phone), logger.Error("error", err))
			return Response{}, NewBusinessError("failed to find or create user", err)
		}

		// 设置 Session
		ssid, err := ctl.userSvc.SetSession(ctx, req.Phone)
		if err != nil {
			ctl.l.Error(fmt.Sprintf("%s:设置 Session 失败", req.Biz), logger.String("phone", req.Phone), logger.Error("error", err))
			return Response{}, NewBusinessError("failed to set session: %w", err)
		}

		// 生成 JWT
		err = ctl.jwtHdl.GenerateToken(c, ssid, user.IsMerchant)
		if err != nil {
			ctl.l.Error(fmt.Sprintf("%s:生成 JWT 失败", req.Biz), logger.String("phone", req.Phone), logger.Error("error", err))
			return Response{}, NewBusinessError("failed to generate JWT: %w", err)
		}

		// 打印成功日志
		maskedPhone := req.Phone[:3] + "****" + req.Phone[len(req.Phone)-4:]
		ctl.l.Info(fmt.Sprintf("%s:用户处理成功", req.Biz), logger.String("phone", maskedPhone))

		// 返回成功响应
		return GetResponse(WithStatus(http.StatusOK), WithMsg(fmt.Sprintf("%s successfully", req.Biz)), WithData(map[string]interface{}{
			"id":    user.Id,
			"phone": user.Phone,
			"name":  user.Name,
		})), nil
	}, func(c *app.RequestContext, err error) (Response, bool) {
		// 根据错误类型处理特定错误
		var busErr *BusinessError
		switch {
		case errors.Is(err, service.ErrVerifyTooMany):
			ctl.l.Error(fmt.Sprintf("%s:校验验证码:校验次数过多", c.Param("biz")), logger.Error("error", err))
			return GetResponse(WithStatus(http.StatusTooManyRequests), WithMsg("verify too many")), true
		case errors.As(err, &busErr):
			ctl.l.Error(busErr.Message, logger.Error("error", busErr.Err))
			return GetResponse(WithStatus(http.StatusInternalServerError), WithMsg(busErr.Message)), true
		default:
			ctl.l.Error("校验验证码:系统错误", logger.Error("error", err))
			return GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")), false
		}
	})
}

func (ctl *UserHandler) UpdatePassword() app.HandlerFunc {
	return WrapReq(func(ctx context.Context, c *app.RequestContext, req struct {
		UserId          uint64 `json:"user_id"`
		Password        string `json:"password" validate:"required,min=8,containsany=abcdefghijklmnopqrstuvwxyz,containsany=0123456789,containsany=$@$!%*#?&"`
		ConfirmPassword string `json:"confirmPassword" validate:"eqfield=Password"`
	}) (Response, error) {
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			return Response{}, err
		}

		err := ctl.userSvc.UpdatePassword(ctx, domain.User{
			Id:       req.UserId,
			Password: req.Password,
		})
		if err != nil {
			return Response{}, err
		}

		ctl.l.Info("更新用户密码成功")
		return GetResponse(WithStatus(http.StatusOK), WithMsg("bind user's password successfully")), nil
	}, func(c *app.RequestContext, err error) (Response, bool) {
		switch {
		case errors.As(err, &validator.ValidationErrors{}):
			ctl.l.Error("更新用户密码:验证失败", logger.Error("error", err))
			return GetResponse(WithStatus(http.StatusBadRequest), WithMsg("Validation failed")), true
		default:
			ctl.l.Error("更新用户密码:系统错误", logger.Error("error", err))
			return GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")), false
		}
	})
}

func (ctl *UserHandler) UpdateBirthday() app.HandlerFunc {
	return WrapReq(func(ctx context.Context, c *app.RequestContext, req struct {
		UserId   uint64    `json:"user_id"`
		Birthday time.Time `json:"birthday"`
	}) (Response, error) {
		err := ctl.userSvc.UpdateBirthday(ctx, domain.User{
			Id:       req.UserId,
			Birthday: req.Birthday,
		})
		if err != nil {
			return Response{}, err
		}

		ctl.l.Info("更新用户生日成功")
		return GetResponse(WithStatus(http.StatusOK), WithMsg("update user's birthday successfully")), nil
	}, func(c *app.RequestContext, err error) (Response, bool) {
		ctl.l.Error("更新用户生日:系统错误", logger.Error("error", err))
		return GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")), false
	})
}

func (ctl *UserHandler) UpdateName() app.HandlerFunc {
	return WrapReq(func(ctx context.Context, c *app.RequestContext, req struct {
		UserId uint64 `json:"user_id"`
		Name   string `json:"name" validate:"required,min=3,max=20"`
	}) (Response, error) {
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			return Response{}, err
		}

		err := ctl.userSvc.UpdateName(ctx, domain.User{
			Id:   req.UserId,
			Name: req.Name,
		})
		if err != nil {
			return Response{}, err
		}

		ctl.l.Info("更新用户名成功")
		return GetResponse(WithStatus(http.StatusOK), WithMsg("update user's name successfully")), nil
	}, func(c *app.RequestContext, err error) (Response, bool) {
		switch {
		case errors.As(err, &validator.ValidationErrors{}):
			ctl.l.Error("更新用户名:校验失败", logger.Error("error", err))
			return GetResponse(WithStatus(http.StatusBadRequest), WithMsg("Validation failed")), true
		default:
			ctl.l.Error("系统错误", logger.Error("error", err))
			return GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")), false
		}
	})
}

func (ctl *UserHandler) NameLogin() app.HandlerFunc {
	return WrapReq(func(ctx context.Context, c *app.RequestContext, req struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}) (Response, error) {
		user, err := ctl.userSvc.NameLogin(ctx, domain.User{
			Name:     req.Name,
			Password: req.Password,
		})
		if err != nil {
			return Response{}, err
		}

		// 设置 Session
		ssid, err := ctl.userSvc.SetSession(ctx, user.Phone)
		if err != nil {
			ctl.l.Error("用户名登录:设置 Session 失败", logger.String("phone", user.Phone), logger.Error("error", err))
			return Response{}, NewBusinessError("failed to set session", err)
		}

		// 生成 JWT
		err = ctl.jwtHdl.GenerateToken(c, ssid, user.IsMerchant)
		if err != nil {
			ctl.l.Error("用户名登录:生成 JWT 失败", logger.String("phone", user.Phone), logger.Error("error", err))
			return Response{}, NewBusinessError("failed to generate JWT", err)
		}

		ctl.l.Info("用户登录成功")
		return GetResponse(WithStatus(http.StatusOK), WithMsg("name login successfully")), nil
	}, func(c *app.RequestContext, err error) (Response, bool) {
		// 根据错误类型记录日志
		var busErr *BusinessError
		switch {
		case errors.Is(err, service.ErrRecordNotFound):
			ctl.l.Error("用户名登录:用户不存在", logger.Error("error", err))
			return GetResponse(WithStatus(http.StatusNotFound), WithMsg("user not found")), true
		case errors.Is(err, service.ErrInvalidUserOrPassword):
			ctl.l.Error("用户名登录:用户名或密码错误", logger.Error("error", err))
			return GetResponse(WithStatus(http.StatusUnauthorized), WithMsg("username or password error")), true
		case errors.As(err, &busErr):
			ctl.l.Error(busErr.Message, logger.Error("error", busErr.Err))
			return GetResponse(WithStatus(http.StatusInternalServerError), WithMsg(busErr.Message)), true
		default:
			ctl.l.Error("系统错误", logger.Error("error", err))
			return GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")), false
		}
	})
}

func (ctl *UserHandler) Logout() app.HandlerFunc {
	return WrapReq(func(ctx context.Context, c *app.RequestContext, req struct{}) (Response, error) {
		// 提取和解析 token
		token := ctl.jwtHdl.ExtractToken(c)
		claim, err := ctl.jwtHdl.ParseToken(token)
		if err != nil {
			ctl.l.Error("退出登录:解析 Token 错误", logger.Error("error", err))
			return Response{}, NewBusinessError("logout failed parse token", err)
		}

		// 删除 session
		err = ctl.userSvc.DeleteSession(ctx, claim.IsMerchant, claim.Id)
		if err != nil {
			ctl.l.Error("退出登录:删除 Session 错误", logger.Error("error", err))
			return Response{}, NewBusinessError("logout failed delete session", err)
		}

		// 记录成功日志
		ctl.l.Info("用户退出登录成功", logger.String("user_id", claim.Id))
		return GetResponse(WithStatus(http.StatusOK), WithMsg("log out successfully")), nil
	}, func(c *app.RequestContext, err error) (Response, bool) {
		// 根据错误类型记录日志
		var busErr *BusinessError
		if errors.As(err, &busErr) {
			ctl.l.Error(busErr.Message, logger.Error("error", busErr.Err))
			return GetResponse(WithStatus(http.StatusInternalServerError), WithMsg(busErr.Message)), true
		}

		ctl.l.Error("系统错误", logger.Error("error", err))
		return GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")), false
	})
}

func (ctl *UserHandler) KeepAive() app.HandlerFunc {
	return WrapReq(func(ctx context.Context, c *app.RequestContext, req struct{}) (Response, error) {
		tokenHeader := ctl.jwtHdl.ExtractToken(c)

		claims, err := ctl.jwtHdl.ParseToken(tokenHeader)
		if err != nil {
			ctl.l.Error("维持登录状态:解析 Token 错误")
			return Response{}, NewBusinessError("keep alive:parse token failed", err)
		}

		// 刷新 Session 有效期
		err = ctl.userSvc.ExtendSessionExpiration(ctx, claims.IsMerchant, claims.SessionId)
		if err != nil {
			ctl.l.Error("维持登录状态:刷新 session 错误")
			return Response{}, NewBusinessError("keep alive:extend session failed", err)
		}

		ctl.l.Info("维持登录状态成功")
		return GetResponse(WithStatus(http.StatusOK), WithMsg("session refreshed")), nil
	}, func(c *app.RequestContext, err error) (Response, bool) {
		var busErr *BusinessError
		if errors.As(err, &busErr) {
			ctl.l.Error(busErr.Message, logger.Error("error", busErr.Err))
			return GetResponse(WithStatus(http.StatusInternalServerError), WithMsg(busErr.Message)), true
		}

		ctl.l.Error("维持登录状态:系统错误")
		return GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")), false
	})
}

func (ctl *UserHandler) BindAddress() app.HandlerFunc {
	return WrapReq(func(ctx context.Context, c *app.RequestContext, req struct {
		UserId    uint64 `json:"userId"`
		Street    string `json:"street"`
		City      string `json:"city"`
		State     string `json:"state"`
		ZipCode   string `json:"zipCode"`
		Country   string `json:"country"`
		IsDefault bool   `json:"isDefault"`
	}) (Response, error) {
		err := ctl.userSvc.BindAddress(ctx, domain.Address{
			UserId:    req.UserId,
			Street:    req.Street,
			State:     req.State,
			City:      req.City,
			ZipCode:   req.ZipCode,
			Country:   req.Country,
			IsDefault: req.IsDefault,
		})
		if err != nil {
			return Response{}, err
		}

		ctl.l.Info("添加地址成功")
		return GetResponse(WithStatus(http.StatusOK), WithMsg("add addr successfully")), nil
	}, func(c *app.RequestContext, err error) (Response, bool) {
		ctl.l.Error("添加地址:系统错误")
		return GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system error")), false
	})
}

func (ctl *UserHandler) AcquireAllAddr() app.HandlerFunc {
	return WrapReq(func(ctx context.Context, c *app.RequestContext, req struct{}) (Response, error) {
		id := c.Param("id")
		addresses, err := ctl.userSvc.AcquireAllAddr(ctx, id)
		if err != nil {
			ctl.l.Error("查询所有地址:系统错误")
			return Response{}, err
		}

		ctl.l.Info("查询所有地址成功")
		return GetResponse(WithStatus(http.StatusOK), WithMsg("acquire all addresses successfully"), WithData(addresses)), nil
	}, func(c *app.RequestContext, err error) (Response, bool) {
		return GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("查询所有地址:系统错误")), false
	})
}

func (ctl *UserHandler) DeleteAddress() app.HandlerFunc {
	return WrapReq(func(ctx context.Context, c *app.RequestContext, req struct{}) (Response, error) {
		id := c.Param("id")

		err := ctl.userSvc.DeleteAddress(ctx, id)
		if err != nil {
			ctl.l.Error("删除地址:系统错误")
			return Response{}, err
		}

		ctl.l.Info("删除地址成功")
		return GetResponse(WithStatus(http.StatusOK), WithMsg("delete addr successfully")), nil
	}, func(c *app.RequestContext, err error) (Response, bool) {
		return GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("删除地址:系统错误")), false
	})
}

func (ctl *UserHandler) UpdateAddress() app.HandlerFunc {
	return WrapReq(func(ctx context.Context, c *app.RequestContext, req struct {
		UserId    uint64 `json:"userId"`
		Street    string `json:"street"`
		City      string `json:"city"`
		State     string `json:"state"`
		ZipCode   string `json:"zipCode"`
		Country   string `json:"country"`
		IsDefault bool   `json:"isDefault"`
	}) (Response, error) {
		addr, err := ctl.userSvc.UpdateAddress(ctx, domain.Address{
			UserId:    req.UserId,
			Street:    req.Street,
			State:     req.State,
			City:      req.City,
			ZipCode:   req.ZipCode,
			Country:   req.Country,
			IsDefault: req.IsDefault,
		})
		if err != nil {
			ctl.l.Error("更新地址:系统错误")
			return Response{}, err
		}

		ctl.l.Info("更新地址成功")
		return GetResponse(WithStatus(http.StatusOK), WithMsg("update address successfully"), WithData(addr)), nil
	}, func(c *app.RequestContext, err error) (Response, bool) {
		ctl.l.Error("更新地址:系统错误")
		return GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("system erro")), false
	})
}
