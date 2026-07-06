package ctrl

import "cloudflared-tunnel/internal/module/user/svc"

type Ctrl struct {
	UserSvc svc.UserSvc
}

func NewCtrl(us svc.UserSvc) *Ctrl {
	return &Ctrl{UserSvc: us}
}
