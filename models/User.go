package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

// init 初始化
func init() {
	orm.RegisterModel(new(User))
}

//UserTBName .
func UserTBName() string {
	table := "user"
	return table
}

// TableName 设置User表名
func (a *User) TableName() string {
	return UserTBName()
}

// UserQueryParam 用于查询的类
type UserQueryParam struct {
	BaseQueryParam
	SearchLike string //模糊查询
	NameLike   string
}

// User 实体类
type User struct {
	UserId  string `orm:"size(50);pk"`
	Name    string `orm:"size(50)"`
	Pwd     string
	Mobile  string `orm:"size(36)"`
	Address string `orm:"size(100)"`

	VoicePrint string
	Major      string
	Role       int
	Department string
	Gender     int
	IdCard     string
	VoiceUrl   string

	Post         string `orm:"size(50)"`
	Safetybelt   string `orm:"size(50)"`
	Safetyhelmet string `orm:"size(50)"`
	//Roles              []*Role         `orm:"rel(m2m);"` // 设置多对多
	CreateDatetime time.Time `orm:"auto_now_add;type(datetime)"`
	UpdateDatetime time.Time `orm:"auto_now;type(datetime)"`

	Token   string `orm:"-"`
	Expires int64  `orm:"-"`
}

// UserPageList 获取分页数据
func UserPageList(params *UserQueryParam, curUser *User) ([]*User, int64) {
	query := orm.NewOrm().QueryTable(UserTBName()).Distinct()
	data := make([]*User, 0)
	//默认排序
	sortorder := "UserId"
	switch params.Sort {
	case "UserId":
		sortorder = "UserId"
	default:
		sortorder = "UpdateDatetime"
	}
	if params.Order != "asc" {
		sortorder = "-" + sortorder
	}
	if params.SearchLike != "" {
		cond := orm.NewCondition()
		query = query.SetCond(
			cond.AndCond(
				cond.Or(
					"Name__icontains", params.SearchLike,
				).Or(
					"UserId__icontains", params.SearchLike,
				).Or(
					"Mobile__icontains", params.SearchLike,
				).Or(
					"VoicePrint__contains", params.SearchLike,
				),
			),
		)
	}

	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	return data, total
}

// UserOne 根据id获取单条
func UserOne(id string) (*User, error) {
	o := orm.NewOrm()
	m := User{UserId: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

//UserOneByCheck 获取单条
func UserOneByCheck(login, userpwd string) (*User, error) {
	m := User{}
	O := orm.NewOrm()
	query := O.QueryTable(UserTBName())
	cond := orm.NewCondition()
	query = query.SetCond(
		cond.AndCond(
			cond.Or(
				"UserId", login,
			).Or(
				"Name", login,
			),
		),
	)
	query = query.Filter("Pwd", userpwd)
	err := query.Filter("Name", login).One(&m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

//UserOneByMobile 根据手机号密码获取单条 .
func UserOneByMobile(mobile string) (*User, error) {
	m := User{}
	O := orm.NewOrm()
	query := O.QueryTable(UserTBName())
	query = query.Filter("Mobile", mobile)
	//query = query.Filter("IsDelete", false)
	err := query.One(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

//UserOneByName .
func UserOneByName(name string) (*User, error) {
	m := User{}
	O := orm.NewOrm()
	query := O.QueryTable(UserTBName())
	query = query.Filter("Name", name)

	err := query.One(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
