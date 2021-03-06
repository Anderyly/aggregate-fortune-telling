/*
 * @Author anderyly
 * @email admin@aaayun.cc
 * @link http://blog.aaayun.cc/
 * @copyright Copyright (c) 2022
 */

package common

type XunKongModel struct {
}

var (
	xkarr = [][]string{
		{"甲子", "乙丑", "丙寅", "丁卯", "戊辰", "己巳", "庚午", "辛未", "壬申", "癸酉", "戌亥"},
		{"甲戌", "乙亥", "丙子", "丁丑", "戊寅", "己卯", "庚辰", "辛巳", "壬午", "癸未", "申酉"},
		{"甲申", "乙酉", "丙戌", "丁亥", "戊子", "己丑", "庚寅", "辛卯", "壬辰", "癸巳", "午未"},
		{"甲午", "乙未", "丙申", "丁酉", "戊戌", "己亥", "庚子", "辛丑", "壬寅", "癸卯", "辰巳"},
		{"甲辰", "乙巳", "丙午", "丁未", "戊申", "己酉", "庚戌", "辛亥", "壬子", "癸丑", "寅卯"},
		{"甲寅", "乙卯", "丙辰", "丁巳", "戊午", "己未", "庚申", "辛酉", "壬戌", "癸亥", "子丑"},
	}
)

func (con XunKongModel) Get(gz string) string {
	tag := 0
	xunk := ""
	for i := 0; i <= 5; i++ {
		for j := 0; j <= 10; j++ {
			if xkarr[i][j] == gz {
				xunk = xkarr[i][10]
				tag = 1
				break
			}
		}
		if tag == 1 {
			break
		}
	}
	return xunk
}
