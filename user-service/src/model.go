package main

type user struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password"`
}

func (u *user) toString() string {
	return "user{Username=" + u.Username + ", Password=" + u.Password + "}"
}

type passwords struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

func (p *passwords) toString() string {
	return "passwords{OldPassWord=" + p.OldPassword + ", NewPassword=" + p.NewPassword + "}"
}

func (p *passwords) isSame() bool {
	return p.OldPassword == p.NewPassword
}
