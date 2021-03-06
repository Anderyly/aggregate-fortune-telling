package api

import (
	"aggregate-fortune-telling/ay"
	"aggregate-fortune-telling/models"
	"fmt"
	"github.com/gin-gonic/gin"
	qrcode "github.com/skip2/go-qrcode"
	"strconv"
	"time"
)

var (
	Token string
	Appid int64
)

type CommonController struct {
}

func (con CommonController) A() {

}

func AuthMaster() (bool, string, models.User) {
	var user models.User
	ay.Db.First(&user, "id = ?", GetToken(Token))

	if user.Id == 0 {
		return false, "大师不存在", user
	}

	if user.Type != 1 {
		return false, "大师不存在", user
	}

	return true, "", user
}

func GetRequestIP(c *gin.Context) string {
	reqIP := c.ClientIP()
	if reqIP == "::1" {
		reqIP = "127.0.0.1"
	}
	return reqIP
}

func GetToken(token string) string {
	uid := ay.AuthCode(token, "DECODE", "", 0)
	return uid
}

func DateFormat(y, m, d, h, i, s int) string {
	vm := strconv.Itoa(m)
	if len(vm) == 1 {
		vm = "0" + vm
	}

	vd := strconv.Itoa(d)
	if len(vd) == 1 {
		vd = "0" + vd
	}

	vh := strconv.Itoa(h)
	if len(vh) == 1 {
		vh = "0" + vh
	}

	vi := strconv.Itoa(i)
	if len(vi) == 1 {
		vi = "0" + vi
	}

	vs := strconv.Itoa(s)
	if len(vs) == 1 {
		vs = "0" + vs
	}

	return strconv.Itoa(y) + "-" + vm + "-" + vd + " " + vh + ":" + vi + ":" + vs
}

func (con CommonController) MakeQrCode(text string) string {
	name := ay.MD5(fmt.Sprintf("%s%s", text, time.Now().String())) + ".png"
	fileDir := fmt.Sprintf("static/qrcode/%d-%d/", time.Now().Year(), time.Now().Month())

	err := ay.CreateMutiDir(fileDir)
	if err != nil {
		return ""
	}

	err = qrcode.WriteFile(text, qrcode.Medium, 152, fileDir+name)
	if err != nil {
		return ""
	}
	return fileDir + name

}
