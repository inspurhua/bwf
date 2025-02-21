package dao

import (
	"errors"
	"huage.tech/mini/app/bean"
	"huage.tech/mini/app/config"
	"time"
)

func Login(account, password string) (u bean.User, err error) {
	err = db.Model(&bean.User{}).Where("account=$1 and password=md5($2) and status=1", account,
		config.JwtSecret+password).First(&u).Error
	//非超级管理员时判断此人的tenant状态是有效的
	if err != nil && u.TenantId > 0 {
		t, _ := TenantRead(u.TenantId)
		if t.Status == 2 {
			err = errors.New("该租户已经被禁止使用")
			return
		}
	}
	return
}

func UserList(tenantId int64, account string, roleId int64, orgId int64, offset, limit int64) (r []bean.UserResponse, count int, err error) {
	t, _ := TenantRead(tenantId)

	sql1 := "select count(u.id) from " + config.Prefix + "user u left join " + config.Prefix + "org o on u.org_id=o.id" +
		" where u.tenant_id = ? and u.account like ? and u.role_id != ?"
	param1 := []interface{}{tenantId, "%" + account + "%", t.RoleAdmin}

	sql := "select u.id,u.account,u.status,u.role_id,u.org_id,u.update_at,r.name as role,o.name as org from " + config.Prefix + "user u " +
		" left join " + config.Prefix + "role r on u.role_id = r.id " +
		" left join " + config.Prefix + "org o on u.org_id = o.id " +
		" where u.tenant_id = ? and u.account like ? and u.role_id != ?"
	param := []interface{}{tenantId, "%" + account + "%", t.RoleAdmin}

	if roleId > 0 {
		sql1 += "and u.role_id = ? "
		param1 = append(param1, roleId)

		sql += "and u.role_id = ? "
		param = append(param, roleId)
	}
	if orgId > 0 {
		sql1 += " and u.org_id = ? "
		param1 = append(param1, orgId)

		sql += " and u.org_id = ? "
		param = append(param, orgId)
	}
	err = db.Raw(sql1, param1...).Row().Scan(&count)
	if err != nil {
		return
	}
	sql += " order by u.update_at desc offset ? limit ?"
	param = append(param, offset, limit)
	err = db.New().Raw(sql, param...).Scan(&r).Error

	return
}

func UserCreate(User bean.User) (result bean.User, err error) {
	now := time.Now()
	User.CreateAt = now
	User.UpdateAt = now
	result = User

	err = db.Create(&result).Error
	result.Password = ""
	return
}

func UserDelete(tenantId, id int64) (err error) {
	err = db.
		Where("tenant_id=?", tenantId).
		Where("id=?", id).
		Delete(&bean.User{}).Error
	return
}

func UserRead(tenantId, id int64) (result bean.User, err error) {
	err = db.
		Where("tenant_id=?", tenantId).
		Where("id=?", id).First(&result).Error
	result.Password = ""
	return
}

func UserUpdate(User bean.User) (result bean.User, err error) {
	User.UpdateAt = time.Now()
	err = db.Model(&result).Update(User).Error
	result.Password = ""
	return
}
