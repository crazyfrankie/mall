package user

import "mall/internal/user/web"

type Handler = web.UserHandler // 暴露出去给 ioc 使用
