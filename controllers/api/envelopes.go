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
	"time"
)

type EnvelopesController struct {
}

type EnvelopesSendForm struct {
	Amount    float64 `form:"amount" binding:"required" label:"金额"`
	To        string  `form:"to" binding:"required" label:"对象"`
	Remark    string  `form:"remark"`
	PayType   int     `form:"pay_type" binding:"required" label:"支付方式"`
	ReturnUrl string  `form:"return_url" binding:"required" label:"同步地址"`
}

// Send 用户给大师发红包
func (con EnvelopesController) Send(c *gin.Context) {
	var getForm EnvelopesSendForm
	if err := c.ShouldBind(&getForm); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	var user models.User
	ay.Db.First(&user, "id = ?", GetToken(Token))

	if user.Id == 0 {
		ay.Json{}.Msg(c, 401, "Token错误", gin.H{})
		return
	}

	var master models.User
	ay.Db.Where("phone = ?", getForm.To).First(&master)

	if master.Id == 0 {
		ay.Json{}.Msg(c, 400, "大师不存在", gin.H{})
		return
	}

	oid := ay.MakeOrder(time.Now())

	order := &models.Order{
		Oid:        oid,
		Type:       7,
		Ip:         c.ClientIP(),
		Des:        "发送红包" + strconv.FormatFloat(getForm.Amount, 'g', -1, 64) + "元",
		Amount:     getForm.Amount,
		Uid:        user.Id,
		Status:     -1,
		Appid:      Appid,
		OutTradeNo: oid,
		Op:         2,
		OldAmount:  getForm.Amount,
		ToUid:      master.Id,
		Remark:     getForm.Remark,
		PayType:    9,
		ReturnUrl:  getForm.ReturnUrl + "&oid=" + oid,
	}

	if err := ay.Db.Create(&order).Error; err != nil {
		ay.Json{}.Msg(c, 400, "数据错误", gin.H{})
		return
	}

	if getForm.PayType == 3 {
		if user.Amount < getForm.Amount {
			ay.Json{}.Msg(c, 406, "账户余额不足", gin.H{})
			return
		}

		user.Amount = user.Amount - getForm.Amount
		if err := ay.Db.Save(&user).Error; err != nil {
			ay.Json{}.Msg(c, 400, "请联系管理员", gin.H{})
			return
		}

		order.Status = 0
		ay.Db.Save(&order)

		//// 上级消费
		//models.UserInviteConsumptionModel{}.Set(user.Id, user.Pid, getForm.Amount, oid)
		ay.Json{}.Msg(c, 200, "发送成功", gin.H{
			"oid":    oid,
			"remark": getForm.Remark,
		})
	} else {
		if Appid == 1 {
			if getForm.PayType == 1 {
				ay.Json{}.Msg(c, 200, "success", gin.H{
					"url":    ay.Yaml.GetString("domain") + "/pay/alipay?oid=" + order.Oid,
					"oid":    oid,
					"remark": getForm.Remark,
				})
				return
			} else {
				ay.Json{}.Msg(c, 200, "success", gin.H{
					"url":    ay.Yaml.GetString("domain") + "/pay/wechat?oid=" + order.Oid,
					"oid":    oid,
					"remark": getForm.Remark,
				})
				return
			}
		} else if Appid == 2 {
			// 百度支付
			is, msg, rj := BaiDuController{}.Baidu(order.Oid)

			if is {
				ay.Json{}.Msg(c, 200, "success", gin.H{
					"info": rj,
				})
				return
			} else {
				ay.Json{}.Msg(c, 400, msg, gin.H{})
				return
			}
		}

	}
}

type EnvelopesDetailForm struct {
	Oid float64 `form:"oid" binding:"required" label:"订单号"`
}

// Detail 红包详情
func (con EnvelopesController) Detail(c *gin.Context) {
	var getForm EnvelopesDetailForm
	if err := c.ShouldBind(&getForm); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	var user models.User
	ay.Db.First(&user, "id = ?", GetToken(Token))

	if user.Id == 0 {
		ay.Json{}.Msg(c, 401, "Token错误", gin.H{})
		return
	}

	var order models.Order
	ay.Db.Where("oid = ?", getForm.Oid).First(&order)

	if order.Id == 0 {
		ay.Json{}.Msg(c, 400, "订单不存在", gin.H{})
		return
	}
	ay.Json{}.Msg(c, 200, "success", gin.H{
		"avatar":   ay.Yaml.GetString("domain") + user.Avatar,
		"remark":   order.Remark,
		"amount":   order.Amount,
		"status":   order.Status,
		"nickname": user.NickName,
	})

}

type EnvelopesRewardForm struct {
	Amount    float64 `form:"amount" binding:"required" label:"金额"`
	To        string  `form:"to" binding:"required" label:"对象"`
	PayType   int     `form:"pay_type" binding:"required" label:"支付方式"`
	ReturnUrl string  `form:"return_url" binding:"required" label:"同步地址"`
}

// Reward 打赏
func (con EnvelopesController) Reward(c *gin.Context) {
	var getForm EnvelopesRewardForm
	if err := c.ShouldBind(&getForm); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	var user models.User
	ay.Db.First(&user, "id = ?", GetToken(Token))

	if user.Id == 0 {
		ay.Json{}.Msg(c, 401, "Token错误", gin.H{})
		return
	}

	var master models.User
	ay.Db.Where("phone = ?", getForm.To).First(&master)

	if master.Id == 0 {
		ay.Json{}.Msg(c, 400, "大师不存在", gin.H{})
		return
	}

	oid := ay.MakeOrder(time.Now())

	order := &models.Order{
		Oid:        oid,
		Type:       5,
		Ip:         c.ClientIP(),
		Des:        "打赏" + strconv.FormatFloat(getForm.Amount, 'g', -1, 64) + "元",
		Amount:     getForm.Amount,
		Uid:        user.Id,
		Appid:      Appid,
		OutTradeNo: oid,
		Op:         2,
		OldAmount:  getForm.Amount,
		ToUid:      master.Id,
		Status:     0,
		ReturnUrl:  getForm.ReturnUrl,
		PayType:    9,
	}

	if err := ay.Db.Create(order).Error; err != nil {
		ay.Json{}.Msg(c, 400, "数据错误", gin.H{})
		return
	}

	if getForm.PayType == 3 {
		if user.Amount < getForm.Amount {
			ay.Json{}.Msg(c, 406, "账户余额不足", gin.H{})
			return
		}
		// 扣除用户余额
		user.Amount -= getForm.Amount
		if err := ay.Db.Save(&user).Error; err != nil {
			ay.Json{}.Msg(c, 400, "数据错误", gin.H{})
			return
		}
		// 增加大师余额
		master.Amount += getForm.Amount
		if err := ay.Db.Save(&master).Error; err == nil {
			// 上级消费
			models.UserInviteConsumptionModel{}.Set(user.Id, user.Pid, getForm.Amount, oid)
			ay.Json{}.Msg(c, 200, "打赏成功", gin.H{})
		} else {
			ay.Json{}.Msg(c, 400, "数据错误", gin.H{})
		}
	} else {
		if Appid == 1 {
			if getForm.PayType == 1 {
				ay.Json{}.Msg(c, 200, "success", gin.H{
					"url": ay.Yaml.GetString("domain") + "/pay/alipay?oid=" + order.Oid,
				})
				return
			} else {
				ay.Json{}.Msg(c, 200, "success", gin.H{
					"url": ay.Yaml.GetString("domain") + "/pay/wechat?oid=" + order.Oid,
				})
				return
			}
		} else {
			//百度支付
			is, msg, rj := BaiDuController{}.Baidu(order.Oid)

			if is {
				ay.Json{}.Msg(c, 200, "success", gin.H{
					"info": rj,
				})
				return
			} else {
				ay.Json{}.Msg(c, 400, msg, gin.H{})
				return
			}
		}

	}

}

type envelopesReceiveForm struct {
	Oid string `form:"oid" binding:"required" label:"订单号"`
}

// Receive 大师接收红包
func (con EnvelopesController) Receive(c *gin.Context) {
	var getForm envelopesReceiveForm
	if err := c.ShouldBind(&getForm); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	isMaster, msg, master := AuthMaster()
	if !isMaster {
		ay.Json{}.Msg(c, 401, msg, gin.H{})
		return
	}

	var order models.Order
	ay.Db.Where("oid = ? and type = 7", getForm.Oid).First(&order)

	if order.Id == 0 {
		ay.Json{}.Msg(c, 400, "红包不存在", gin.H{})
		return
	}

	if order.ToUid != master.Id {
		ay.Json{}.Msg(c, 400, "这不是您的红包", gin.H{})
		return
	}

	if order.Status == 1 {
		ay.Json{}.Msg(c, 400, "该红包已领取过", gin.H{})
		return
	}

	tx := ay.Db.Begin()

	order.Status = 1

	if err := tx.Save(&order).Error; err == nil {
		master.Amount += order.Amount
		if err := tx.Save(&master).Error; err == nil {
			tx.Commit()
			ay.Json{}.Msg(c, 200, "成功领取红包", gin.H{})
			return
		} else {
			tx.Rollback()
			ay.Json{}.Msg(c, 400, "请联系管理员", gin.H{})
			return
		}
	} else {
		tx.Rollback()
		ay.Json{}.Msg(c, 400, "请联系管理员", gin.H{})
		return
	}

}
