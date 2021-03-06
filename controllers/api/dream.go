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
)

type DreamController struct {
}

// Main 首页
func (con DreamController) Main(c *gin.Context) {
	res := models.DreamTypeModel{}.GetAllType()
	recommend := models.DreamModel{}.GetList("0", 20)

	for k, v := range recommend {
		recommend[k].Message = ay.Summary(v.Message, 50)
	}

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"list":      res,
		"recommend": recommend,
		//"notice":    notice,
	})
}

type GetDreamForm struct {
	Title string `form:"title" binding:"required" label:"标题"`
}

// Search 搜索
func (con DreamController) Search(c *gin.Context) {
	var getForm GetDreamForm
	if err := c.ShouldBind(&getForm); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}
	res := models.DreamModel{}.GetList(getForm.Title, 10)

	for k, v := range res {
		res[k].Message = ay.Summary(v.Message, 50)
	}

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"list": res,
	})
}

type GetDreamDetailForm struct {
	Id int `form:"id" binding:"required"`
}

// Detail 详情
func (con DreamController) Detail(c *gin.Context) {
	var getForm GetDreamDetailForm
	if err := c.ShouldBind(&getForm); err != nil {
		ay.Json{}.Msg(c, 400, ay.Validator{}.Translate(err), gin.H{})
		return
	}

	res := models.DreamModel{}.GetDetail(getForm.Id)

	total := models.DreamTotalModel{}.GetTotal(getForm.Id)

	ay.Json{}.Msg(c, 200, "success", gin.H{
		"info":  res,
		"total": total.Num,
		//"notice": notice,
	})
}
