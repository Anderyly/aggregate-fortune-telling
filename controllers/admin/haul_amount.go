/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package admin

import (
	"aggregate-fortune-telling/ay"
	"aggregate-fortune-telling/models"
	"github.com/gin-gonic/gin"
	"strings"
)

type HaulAmountController struct {
}

type haulAmountListForm struct {
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	Key      string `form:"key"`
	Type     string `form:"type"`
}

// List 列表
func (con HaulAmountController) List(c *gin.Context) {
	var data haulAmountListForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	if Auth() == false {
		ay.Json{}.Msg(c, 401, "请登入", gin.H{})
		return
	}

	var list []models.HaulAmount

	var count int64
	res := ay.Db.Table("sm_haul_amount")

	row := res

	row.Count(&count)

	res.Order("sort asc").
		Limit(data.PageSize).
		Offset((data.Page - 1) * data.PageSize).
		Find(&list)

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"list":  list,
		"total": count,
	})
}

//type orderDetailForm struct {
//	Id int `form:"id"`
//}

// Detail 用户详情
func (con HaulAmountController) Detail(c *gin.Context) {
	var data orderDetailForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	if Auth() == false {
		ay.Json{}.Msg(c, 401, "请登入", gin.H{})
		return
	}

	var user models.HaulAmount

	ay.Db.First(&user, data.Id)

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"info": user,
	})
}

type haulAmountOptionForm struct {
	Id     int     `form:"id"`
	Amount float64 `form:"amount"`
	Sort   int     `form:"sort"`
}

// Option 添加 编辑
func (con HaulAmountController) Option(c *gin.Context) {
	var data haulAmountOptionForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	if Auth() == false {
		ay.Json{}.Msg(c, 401, "请登入", gin.H{})
		return
	}

	var res models.HaulAmount
	ay.Db.First(&res, data.Id)

	if data.Id != 0 {

		res.Amount = data.Amount
		res.Sort = data.Sort

		if err := ay.Db.Save(&res).Error; err != nil {
			ay.Json{}.Msg(c, 400, "修改失败", gin.H{})
		} else {
			ay.Json{}.Msg(c, 200, "修改成功", gin.H{})
		}
	} else {
		ay.Db.Create(&models.HaulAmount{
			Amount: data.Amount,
			Sort:   data.Sort,
		})
		ay.Json{}.Msg(c, 200, "创建成功", gin.H{})

	}

}

func (con HaulAmountController) Delete(c *gin.Context) {
	var data orderDeleteForm
	if err := c.ShouldBind(&data); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	if Auth() == false {
		ay.Json{}.Msg(c, 401, "请登入", gin.H{})
		return
	}

	idArr := strings.Split(data.Id, ",")

	for _, v := range idArr {
		var order models.HaulAmount
		ay.Db.Delete(&order, v)
	}

	ay.Json{}.Msg(c, 200, "删除成功", gin.H{})
}
