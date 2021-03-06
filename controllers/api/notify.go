/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package api

import (
	"aggregate-fortune-telling/ay"
	"aggregate-fortune-telling/controllers/advert"
	"aggregate-fortune-telling/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/go-pay/gopay/wechat"
	"log"
	"net/http"
	"strconv"
	"time"
)

type NotifyController struct {
}

func (con NotifyController) AliPay(c *gin.Context) {

	var pay models.Pay
	ay.Db.First(&pay, "id = ?", 7)

	notifyReq, err := alipay.ParseNotifyToBodyMap(c.Request)
	con.CheckErr(err)

	j, err := json.Marshal(notifyReq)
	con.CheckErr(err)

	ay.CreateMutiDir("log/alipay")
	ay.WriteFile("log/alipay/"+notifyReq.Get("out_trade_no")+".txt", []byte(j), 0755)

	ok, err := alipay.VerifySign(pay.PayKey, notifyReq)

	if !ok {
		log.Println(err)
		c.String(http.StatusOK, "%s", "fail")
		return
	}

	var order models.Order
	ay.Db.First(&order, "out_trade_no = ?", notifyReq.Get("out_trade_no"))

	// 查询订单失败
	if order.Id == 0 {
		c.String(http.StatusOK, "%s", "fail")
		return
	}

	// 订单已处理过
	if order.Status == 1 {
		c.String(http.StatusOK, "%s", "success")
		return
	}

	res := 0

	switch order.Type {
	case 9:
		// 增加余额
		res = con.AddUserAmount(order.Uid, order.Amount)
	case 1:
		res = 1
	case 3:
		res = 1
	case 5:
		// 增加大师余额
		res = con.AddMasterAmount(order.Uid, order.ToUid, order.Amount, order.Oid)
	case 7:
		res = 1
	case 6:
		// 增加大师余额
		res = con.AddMasterAmount(order.ToUid, order.Uid, order.Amount, order.Oid)
	}

	if res == 1 {
		// 优惠卷设置过期
		var coupon models.Coupon
		ay.Db.First(&coupon, "id = ?", order.Coupon)
		if coupon.Id != 0 {
			coupon.Status = 1
			coupon.UsedAt = time.Now().Format("2006-01-02 15:04:05")
			ay.Db.Save(&coupon)
		}
		if order.Type == 3 || order.Type == 7 {
			order.Status = 0
		} else {
			order.Status = 1
		}
		order.PayType = 1
		order.TradeNo = notifyReq.Get("trade_no")
		order.PayTime = time.Now().Format("2006-01-02 15:04:05")
		ay.Db.Save(&order)

		if order.Type != 9 && order.Type != 7 && order.Type != 6 {
			// 上级消费
			var user models.User
			ay.Db.First(&user, order.Uid)
			models.UserInviteConsumptionModel{}.Set(user.Id, user.Pid, order.Amount, order.Oid)
		}

		var advertLog models.AdvertLog
		ay.Db.Where("oid = ? and status = 0", order.Oid).First(&advertLog)
		if advertLog.Id != 0 {
			if advertLog.Type == 1 {
				advert.Vivo{}.Up(advertLog.Cid, strconv.FormatFloat(order.Amount, 'g', -1, 64), advertLog.RequestId, advertLog.AdId)
				advertLog.Status = 1
				ay.Db.Save(&advertLog)
			} else {
				advertLog.Status = 1
				ay.Db.Save(&advertLog)
			}
		}

		c.String(http.StatusOK, "%s", "success")
	} else {
		c.String(http.StatusOK, "%s", "fail")
	}

}

func (con NotifyController) WeChat(c *gin.Context) {

	var pay models.Pay
	ay.Db.First(&pay, "id = ?", 6)

	notifyReq, err := wechat.ParseNotifyToBodyMap(c.Request)
	con.CheckErr(err)
	log.Println(notifyReq)

	j, err := json.Marshal(notifyReq)
	con.CheckErr(err)

	ay.CreateMutiDir("log/wechat")
	ay.WriteFile("log/wechat/"+notifyReq.Get("out_trade_no")+".txt", []byte(j), 0755)

	// 验签操作
	ok, err := wechat.VerifySign(pay.VKey, wechat.SignType_MD5, notifyReq)

	if !ok {
		log.Println(err)
		c.String(http.StatusOK, "%s", "fail")
		return
	}

	var order models.Order
	ay.Db.First(&order, "out_trade_no = ?", notifyReq.Get("out_trade_no"))

	// 查询订单失败
	if order.Id == 0 {
		rsp := new(wechat.NotifyResponse) // 回复微信的数据
		rsp.ReturnCode = gopay.FAIL
		rsp.ReturnMsg = gopay.FAIL
		c.String(http.StatusOK, "%s", rsp.ToXmlString())
		return
	}

	// 订单已处理过
	if order.Status == 1 {
		rsp := new(wechat.NotifyResponse) // 回复微信的数据
		rsp.ReturnCode = gopay.SUCCESS
		rsp.ReturnMsg = gopay.OK
		c.String(http.StatusOK, "%s", rsp.ToXmlString())
		return
	}

	res := 0

	switch order.Type {
	case 9:
		// 增加余额
		res = con.AddUserAmount(order.Uid, order.Amount)
	case 1:
		res = 1
	case 3:
		res = 1
	case 5:
		// 增加大师余额
		res = con.AddMasterAmount(order.Uid, order.ToUid, order.Amount, order.Oid)
	case 7:
		res = 1
	case 6:
		// 增加大师余额
		res = con.AddMasterAmount(order.ToUid, order.Uid, order.Amount, order.Oid)
	}

	if res == 1 {
		// 优惠卷设置过期
		var coupon models.Coupon
		ay.Db.First(&coupon, "id = ?", order.Coupon)
		if coupon.Id != 0 {
			coupon.Status = 1
			coupon.UsedAt = time.Now().Format("2006-01-02 15:04:05")
			ay.Db.Save(&coupon)
		}

		if order.Type == 3 || order.Type == 7 {
			order.Status = 0
		} else {
			order.Status = 1
		}

		order.PayType = 2
		order.TradeNo = notifyReq.Get("transaction_id")
		order.PayTime = time.Now().Format("2006-01-02 15:04:05")
		ay.Db.Save(&order)

		if order.Type != 9 && order.Type != 7 && order.Type != 6 {
			// 上级消费
			var user models.User
			ay.Db.First(&user, order.Uid)
			models.UserInviteConsumptionModel{}.Set(user.Id, user.Pid, order.Amount, order.Oid)
		}

		var advertLog models.AdvertLog
		ay.Db.Where("oid = ? and status = 0", order.Oid).First(&advertLog)
		if advertLog.Id != 0 {
			if advertLog.Type == 1 {
				advert.Vivo{}.Up(advertLog.Cid, strconv.FormatFloat(order.Amount, 'g', -1, 64), advertLog.RequestId, advertLog.AdId)
				advertLog.Status = 1
				ay.Db.Save(&advertLog)
			} else {
				advertLog.Status = 1
				ay.Db.Save(&advertLog)
			}

		}

		rsp := new(wechat.NotifyResponse) // 回复微信的数据
		rsp.ReturnCode = gopay.SUCCESS
		rsp.ReturnMsg = gopay.OK
		c.String(http.StatusOK, "%s", rsp.ToXmlString())
	} else {
		rsp := new(wechat.NotifyResponse) // 回复微信的数据
		rsp.ReturnCode = gopay.FAIL
		rsp.ReturnMsg = gopay.FAIL
		c.String(http.StatusOK, "%s", rsp.ToXmlString())
	}

}

// AddUserAmount 增加用户余额
func (con NotifyController) AddUserAmount(uid int64, amount float64) int {
	// 查询用户
	var user models.User
	ay.Db.First(&user, "id = ?", uid)

	if user.Id == 0 {
		log.Println("用户未找到")
		return 0
	}

	user.Amount = user.Amount + amount
	ay.Db.Save(&user)

	return 1

}

func (con NotifyController) CheckErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func (con NotifyController) AddMasterAmount(uid int64, toUid int64, amount float64, oid string) int {
	var master models.User
	ay.Db.Where("id = ?", toUid).First(&master)

	if master.Id == 0 {
		return 0
	}
	// 增加大师余额
	master.Amount += amount
	if err := ay.Db.Save(&master).Error; err == nil {
		return 1
	} else {
		return 0
	}
}

type baiduNotifyController struct {
	UserID         string `form:"userId"`
	OrderID        string `form:"orderId"`
	UnitPrice      int64  `form:"unitPrice"`
	Count          int64  `form:"count"`
	TotalMoney     int64  `form:"totalMoney"`
	PayMoney       int64  `form:"payMoney"`
	PromoMoney     int64  `form:"promoMoney"`
	HbMoney        int64  `form:"hbMoney"`
	HbBalanceMoney int64  `form:"hbBalanceMoney"`
	GiftCardMoney  int64  `form:"giftCardMoney"`
	DealID         string `form:"dealId"`
	PayTime        int64  `form:"payTime"`
	PromoDetail    string `form:"promoDetail"`
	PayType        int64  `form:"payType"`
	PartnerID      int64  `form:"partnerId"`
	Status         int64  `form:"status"`
	TpOrderID      string `form:"tpOrderId"`
	ReturnData     string `form:"returnData"`
	RsaSign        string `form:"rsaSign"`
}

func (con NotifyController) Baidu(c *gin.Context) {
	var data baiduNotifyController
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	log.Println(data.TpOrderID)
	j, err := json.Marshal(data)
	con.CheckErr(err)
	//
	ay.CreateMutiDir("log/baidu")
	ay.WriteFile("log/baidu/"+data.TpOrderID+".txt", []byte(j), 0755)

	type r struct {
		Errno int    `json:"errno"`
		Msg   string `json:"msg"`
		Data  struct {
			IsConsumed int `json:"isConsumed"`
		} `json:"data"`
	}

	var rj r

	var order models.Order
	ay.Db.First(&order, "out_trade_no = ?", data.TpOrderID)

	// 查询订单失败
	if order.Id == 0 {
		rj.Errno = 0
		rj.Msg = "fail"
		rj.Data.IsConsumed = 2
		s, _ := json.Marshal(rj)
		c.String(http.StatusOK, "%s", s)
		//c.String(http.StatusOK, "%s", "fail")
		return
	}

	// 订单已处理过
	if order.Status == 1 {
		rj.Errno = 0
		rj.Msg = "success"
		rj.Data.IsConsumed = 2
		s, _ := json.Marshal(rj)
		c.String(http.StatusOK, "%s", s)
		//c.String(http.StatusOK, "%s", "success")
		return
	}

	res := 0

	switch order.Type {
	case 9:
		// 增加余额
		res = con.AddUserAmount(order.Uid, order.Amount)
	case 1:
		res = 1
	case 3:
		res = 1
	case 5:
		// 增加大师余额
		res = con.AddMasterAmount(order.Uid, order.ToUid, order.Amount, order.Oid)
	case 7:
		res = 1
	case 6:
		// 增加大师余额
		res = con.AddMasterAmount(order.ToUid, order.Uid, order.Amount, order.Oid)
	}

	if res == 1 {
		// 优惠卷设置过期
		var coupon models.Coupon
		ay.Db.First(&coupon, "id = ?", order.Coupon)
		if coupon.Id != 0 {
			coupon.Status = 1
			coupon.UsedAt = time.Now().Format("2006-01-02 15:04:05")
			ay.Db.Save(&coupon)
		}
		if order.Type == 3 || order.Type == 7 {
			order.Status = 0
		} else {
			order.Status = 1
		}
		order.PayType = 1
		order.TradeNo = data.OrderID
		order.PayTime = time.Now().Format("2006-01-02 15:04:05")
		ay.Db.Save(&order)

		if order.Type != 9 && order.Type != 7 && order.Type != 6 {
			// 上级消费
			var user models.User
			ay.Db.First(&user, order.Uid)
			models.UserInviteConsumptionModel{}.Set(user.Id, user.Pid, order.Amount, order.Oid)
		}

		rj.Errno = 0
		rj.Msg = "success"
		rj.Data.IsConsumed = 2
		s, _ := json.Marshal(rj)
		c.String(http.StatusOK, "%s", s)
	} else {
		rj.Errno = 0
		rj.Msg = "fail"
		rj.Data.IsConsumed = 2
		s, _ := json.Marshal(rj)
		c.String(http.StatusOK, "%s", s)
	}

}
