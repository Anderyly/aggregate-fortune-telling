/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package models

type AskLogModel struct {
}

type AskLog struct {
	Id      int64  `json:"id"`
	Type    int    `json:"type"`
	Content string `json:"content"`
}

func (AskLog) TableName() string {
	return "sm_ask_log"
}
