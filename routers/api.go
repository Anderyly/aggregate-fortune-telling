/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package routers

import (
	"aggregate-fortune-telling/controllers/advert"
	"aggregate-fortune-telling/controllers/api"
	"github.com/gin-gonic/gin"
)

func ApiRouters(r *gin.RouterGroup) {

	apiGroup := r.Group("/api/")

	apiGroup.GET("home/main", api.HomeController{}.Home)
	apiGroup.GET("home/config", api.HomeController{}.Config)
	apiGroup.GET("home/adv", api.HomeController{}.Adv)
	apiGroup.GET("home/getAllMaster", api.HomeController{}.GetMasterPhone)

	apiGroup.POST("home/channel", api.HomeController{}.Channel)

	// 登入
	apiGroup.POST("user/login", api.LoginController{}.Login)
	apiGroup.POST("user/send", api.LoginController{}.Send)

	apiGroup.POST("login/baidu/bind", api.LoginController{}.BaiduBind)

	// 地区
	apiGroup.POST("area/get", api.AreaController{}.Get)

	// 收藏
	apiGroup.POST("collect/set", api.CollectController{}.Set)

	// 用户
	apiGroup.POST("user/edit", api.UserController{}.Edit)
	apiGroup.POST("user/info", api.UserController{}.Info)
	apiGroup.POST("user/upload", api.UserController{}.Upload)
	apiGroup.POST("user/coupon", api.UserController{}.Coupon)
	apiGroup.POST("user/collect", api.UserController{}.Collect)
	apiGroup.POST("user/history", api.UserController{}.History)
	apiGroup.POST("user/withdrawal", api.UserController{}.Withdrawal)
	apiGroup.POST("user/log", api.UserController{}.Log)
	apiGroup.POST("user/ask", api.UserController{}.Ask)
	apiGroup.POST("user/haul_bind", api.UserController{}.HaulBind)

	// 红包
	apiGroup.POST("envelopes/send", api.EnvelopesController{}.Send)
	apiGroup.POST("envelopes/reward", api.EnvelopesController{}.Reward)
	apiGroup.POST("envelopes/detail", api.EnvelopesController{}.Detail)
	apiGroup.POST("envelopes/receive", api.EnvelopesController{}.Receive)

	// 索要红包
	apiGroup.POST("force/send", api.ForceController{}.Send)
	apiGroup.POST("force/do", api.ForceController{}.Do)
	apiGroup.POST("force/detail", api.ForceController{}.Detail)
	//apiGroup.POST("envelopes/receive", api.ForceController{}.Receive)

	// 六爻
	apiGroup.GET("divination/init", api.HexapodController{}.Index)
	apiGroup.POST("divination/get", api.HexapodController{}.Get)

	// 周公解梦
	apiGroup.GET("dream/main", api.DreamController{}.Main)
	apiGroup.POST("dream/search", api.DreamController{}.Search)
	apiGroup.POST("dream/detail", api.DreamController{}.Detail)

	// 文章
	apiGroup.POST("news/type", api.NewsController{}.Type)
	apiGroup.POST("news/all", api.NewsController{}.All)
	apiGroup.POST("news/detail", api.NewsController{}.Detail)
	apiGroup.POST("news/recommend", api.NewsController{}.Recommend)

	// 八字
	apiGroup.POST("haul/submit", api.HaulController{}.Submit)
	apiGroup.POST("haul/detail", api.HaulController{}.Detail)
	apiGroup.GET("haul/main", api.HaulController{}.Main)
	apiGroup.GET("haul/notice", api.HaulController{}.Notice)
	apiGroup.POST("haul/coupon", api.HaulController{}.Coupon)

	// 百科/古籍
	apiGroup.POST("notice/search", api.NoticeController{}.Search)
	apiGroup.POST("ancient/classify", api.NoticeController{}.Classify)
	apiGroup.POST("ancient/detail", api.NoticeController{}.Detail)
	apiGroup.POST("baike/detail", api.NoticeController{}.BaiKe)

	// 日历
	apiGroup.POST("calender/get", api.CalenderController{}.Get)

	// 排盘
	apiGroup.POST("plate/submit", api.PlateController{}.Submit)
	apiGroup.POST("plate/detail", api.PlateController{}.Detail)
	apiGroup.POST("plate/info", api.PlateController{}.Info)
	apiGroup.POST("plate/year", api.PlateController{}.Year)
	apiGroup.POST("plate/month", api.PlateController{}.Month)

	// 大师
	apiGroup.GET("master/type", api.MasterController{}.Type)
	apiGroup.POST("master/list", api.MasterController{}.List)
	apiGroup.POST("master/detail", api.MasterController{}.Detail)
	apiGroup.POST("master/recommend", api.MasterController{}.Recommend)
	apiGroup.GET("master/recommend", api.MasterController{}.GetRecommend)

	// 支付
	apiGroup.POST("pay/recharge", api.RechargeController{}.Main)
	apiGroup.POST("pay/do", api.PayDoController{}.Do)

	// 异步
	apiGroup.Any("notify/alipay", api.NotifyController{}.AliPay)
	apiGroup.Any("notify/wechat", api.NotifyController{}.WeChat)
	apiGroup.Any("notify/baidu", api.NotifyController{}.Baidu)

	// 提问
	apiGroup.GET("ask/main", api.AskController{}.Main)
	apiGroup.POST("ask/submit", api.AskController{}.Submit)
	apiGroup.POST("ask/get", api.AskController{}.Get)
	apiGroup.POST("ask/detail", api.AskController{}.Detail)
	apiGroup.POST("ask/adopt", api.AskController{}.Adopt)

	// 邀请
	apiGroup.POST("invite/list", api.InviteController{}.List)
	apiGroup.POST("invite/withdraw", api.InviteController{}.Withdraw)
	apiGroup.POST("invite/do", api.InviteController{}.Do)
	apiGroup.POST("invite/share", api.InviteController{}.Share)

	apiGroup.Any("im/notify", api.ImController{}.Notify)

	apiGroup.Any("advert/vivo", advert.Vivo{}.GetCode)

}
