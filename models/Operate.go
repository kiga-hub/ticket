package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

// init 初始化
func init() {
	orm.RegisterModel(new(Operate))
}

//OperateTBName .
func OperateTBName() string {
	table := "operate"
	return table
}

// TableName 设置User表名
func (a *Operate) TableName() string {
	return OperateTBName()
}

// OperateQueryParam 用于查询的类
type OperateQueryParam struct {
	BaseQueryParam
	SearchLike string //模糊查询
	IdLike     string
	CatalogLike string
}

//Operate 实体类 .
type Operate struct {
	Id                 string `orm:"size(50);pk"`
	UserId             string `orm:"size(50)"`
	OperatecardtaskUrl string `orm:"size(255)"`

	GetoperateTime   time.Time `orm:"auto_now_add;type(datetime)"`
	OperatestartTime time.Time `orm:"auto_now_add;type(datetime)"`
	OperateendTime   time.Time `orm:"auto_now;type(datetime)"`
	Command string `orm:"size(255)"`
	Operatesigner  string `orm:"size(50)"`
	Operatemanager string `orm:"size(50)"`
	Catalog        string `orm:"size(50)"`
}

//OperateOne .
func OperateOne(id string) (*Operate, error) {
	o := orm.NewOrm()
	m := Operate{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

//OperateByUserID .
func OperateByUserID(userid string) ([]*Operate, int64) {
	query := orm.NewOrm().QueryTable(OperateTBName())
	data := make([]*Operate, 0)
	query = query.Filter("UserId", userid)
	total, _ := query.Count()
	query.All(&data)
	return data, total
}

//OperatePageList .
func OperatePageList(params *OperateQueryParam, curUser *User) ([]*Operate, int64) {
	query := orm.NewOrm().QueryTable(OperateTBName()).Distinct()
	data := make([]*Operate, 0)

	//默认排序
	sortorder := "Id"
	switch params.Sort {
	case "Id":
		sortorder = "Id"
	default:
		sortorder = "OperatestartTime"
	}
	if params.Order != "asc" {
		sortorder = "-" + sortorder
	}
	// if params.SearchLike != "" {
	// 	cond := orm.NewCondition()
	// 	query = query.SetCond(
	// 		cond.AndCond(
	// 			cond.Or(
	// 				"UserId__icontains", params.SearchLike,
	// 			).Or(
	// 				"Catalog__icontains", params.SearchLike,
	// 			).Or(
	// 				"Operatesigner__icontains", params.SearchLike,
	// 			).Or(
	// 				"Operatemanager__icontains", params.SearchLike,
	// 			),
	// 		),
	// 	)
	// }
	if params.IdLike != "" {
		query = query.Filter("UserId", params.IdLike)
	}
	if params.CatalogLike != "" {
		query = query.Filter("Catalog", params.CatalogLike)
	}
	if params.SearchLike != "" {
		query = query.Filter("OperatecardtaskUrl", params.SearchLike)
	}

	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	return data, total

}
