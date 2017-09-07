package handler

import (
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/enum"
	"github.com/beewit/pay/global"
	"github.com/labstack/echo"
	"strconv"
)

func CreateMemberTypeOrder(c echo.Context) error {
	mtIdStr := c.FormValue("mtId")
	mtcIdStr := c.FormValue("mtc")
	if mtIdStr != "" && utils.IsValidNumber(mtIdStr) {
		return utils.Error(c, "请正确选择会员类型", nil)
	}
	if mtcIdStr != "" && utils.IsValidNumber(mtcIdStr) {
		return utils.Error(c, "请正确选择会员套餐", nil)
	}
	mtId, _ := strconv.ParseInt(mtIdStr, 10, 64)
	mtcId, _ := strconv.ParseInt(mtcIdStr, 10, 64)
	mt := getMemberType(mtId)
	mtc := getMemberTypeCharge(mtcId)
	if mt == nil {
		return utils.Error(c, "选择的会员类型不存在", nil)
	}
	if mtc == nil {
		return utils.Error(c, "选择的会员套餐不存在", nil)
	}
	return utils.Success(c, "", nil)
}

func GetMemberTypeAndCharge(c echo.Context) error {
	m := make(map[string]interface{})
	sql := "SELECT * FROM member_type WHERE status=? ORDER BY `order` DESC,ct_time DESC"
	rows, _ := global.DB.Query(sql, enum.NORMAL)
	if len(rows) < 1 {
		return utils.NullData(c)
	}
	mt := rows
	m["memberType"] = mt
	sql = "SELECT * FROM member_type_charge WHERE status=? ORDER BY `order` DESC,ct_time DESC"
	rows2, _ := global.DB.Query(sql, enum.NORMAL)
	if len(rows2) > 0 {
		mtc := rows2
		m["memberTypeCharge"] = mtc
	}
	return utils.Success(c, "", m)
}

func getMemberType(id int64) map[string]interface{} {
	sql := "SELECT * FROM member_type WHERE status=? AND id=? LIMIT 1"
	rows, _ := global.DB.Query(sql, enum.NORMAL, id)
	if len(rows) < 1 {
		return nil
	}
	return rows[0]
}

func getMemberTypeCharge(id int64) map[string]interface{} {
	sql := "SELECT * FROM member_type_charge WHERE status=? AND id=? LIMIT 1"
	rows, _ := global.DB.Query(sql, enum.NORMAL, id)
	if len(rows) < 1 {
		return nil
	}
	return rows[0]
}
