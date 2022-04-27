/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package task

import (
	"gin/ay"
	"gin/models"
	"log"
)

func InviteAmount() {
	log.Println("开始计算邀请用户收益")
	var consumption []models.UserInviteConsumption
	ay.Db.Where("status = 0 AND now() >SUBDATE(created_at,interval -3 day)").Find(&consumption)
	for _, v := range consumption {
		var pUser models.User
		ay.Db.First(&pUser, v.Pid)
		pUser.InviteAmount += v.Amount
		if err := ay.Db.Save(&pUser).Error; err == nil {
			ay.Db.Model(models.UserInviteConsumption{}).Where("id = ?", v.Id).Update("status", 1)
		}
	}
}