/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package api

import (
	"aggregate-fortune-telling/ay"
	"aggregate-fortune-telling/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

type HomeController struct {
}

func (con HomeController) Channel(c *gin.Context) {
	type getForm struct {
		Type      int    `form:"type" binding:"required" label:"通道"`
		Channel   int64  `form:"channel" binding:"required" label:"通道"`
		RequestId string `form:"request_id" binding:"required" label:"信息"`
		AdId      string `form:"ad_id" binding:"required" label:"信息"`
	}

	var data getForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	var res models.AdvertLog
	ay.Db.Where("cid = ? and request_id = ? and ad_id = ?", data.Channel, data.RequestId, data.AdId).First(&res)
	if res.Id == 0 {
		ay.Db.Create(&models.AdvertLog{
			Type:      data.Type,
			Oid:       "",
			Cid:       data.Channel,
			Amount:    0,
			RequestId: data.RequestId,
			AdId:      data.AdId,
			Ip:        c.ClientIP(),
			Status:    0,
		})
	}
	ay.Json{}.Msg(c, 200, "通道记录成功", gin.H{})

}

func (con HomeController) Adv(c *gin.Context) {
	// 广告
	t, _ := strconv.Atoi(c.Query("type"))
	adv := models.AdvModel{}.GetType(t)

	for k, v := range adv {
		adv[k].Image = ay.Yaml.GetString("domain") + v.Image
	}

	if adv != nil {
		ay.Json{}.Msg(c, 200, "success", gin.H{
			"list": adv,
		})
	} else {
		ay.Json{}.Msg(c, 200, "success", gin.H{
			"list": []string{},
		})
	}

}

func (con HomeController) Home(c *gin.Context) {

	// 测算量
	var count int64
	var order models.Order
	ay.Db.Model(&order).Where("type = 1 OR type = 2").Count(&count)

	// 广告
	adv := models.AdvModel{}.GetType(1)

	for k, v := range adv {
		adv[k].Image = ay.Yaml.GetString("domain") + v.Image
	}

	// 热门咨询
	recommend := models.ConsultModel{}.GetType(1)
	hot := models.ConsultModel{}.GetType(2)

	var banner []models.Banner
	ay.Db.Order("sort asc").Find(&banner)

	for k, v := range banner {
		banner[k].Image = ay.Yaml.GetString("domain") + v.Image
	}

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"adv":   adv,
		"count": count,
		"consult": gin.H{
			"recommend": recommend,
			"hot":       hot,
		},
		"banner": banner,
	})
}

func (con HomeController) Config(c *gin.Context) {
	config := models.ConfigModel{}.GetId(1)
	// 广告
	adv := models.AdvModel{}.GetType(2)
	for k, v := range adv {
		adv[k].Image = ay.Yaml.GetString("domain") + v.Image
	}

	var shareText []string

	s := models.ShareTextModel{}.GetAll()
	for _, v := range s {
		shareText = append(shareText, v.Content)
	}
	ay.Json{}.Msg(c, 200, "success", gin.H{
		"kf_link":     config.Kf,
		"master_link": config.MasterLink,
		"adv":         adv,
		"invite_rate": config.InviteRate,
		"share_text":  shareText,
	})
}

// GetMasterPhone 获取所有大师手机
func (con HomeController) GetMasterPhone(c *gin.Context) {

	key := c.Query("key")

	if key != "anderyly" {
		ay.Json{}.Msg(c, 400, "鉴权失败", gin.H{})
		return
	}

	type r struct {
		Id       int64  `json:"-"`
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
		Phone    string `json:"phone"`
		Token    string `json:"token"`
	}

	var list []r
	ay.Db.Model(&models.User{}).Where("type = 1").Select("id,nickname,avatar,phone").Find(&list)

	for k, v := range list {
		list[k].Avatar = ay.Yaml.GetString("domain") + v.Avatar
		list[k].Token = ay.AuthCode(strconv.Itoa(int(v.Id)), "ENCODE", "", 0)
	}

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"list": list,
	})
}
