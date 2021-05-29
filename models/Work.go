package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

// init 初始化
func init() {
	orm.RegisterModel(new(Work))
}

//WorkTBName .
func WorkTBName() string {
	table := "work"
	return table
}

// TableName 设置User表名
func (a *Work) TableName() string {
	return WorkTBName()
}

// WorkQueryParam 用于查询的类
type WorkQueryParam struct {
	BaseQueryParam
	SearchLike string //模糊查询
	IdLike     string
	CatalogLike string
}

//Work 实体类 .
type Work struct {
	Id              string `orm:"size(50);pk"`
	UserId          string `orm:"size(50)"`
	WorkcardtaskUrl string `orm:"size(255)"`

	GetworkTime   time.Time `orm:"auto_now_add;type(datetime)"`
	WorkstartTime time.Time `orm:"auto_now_add;type(datetime)"`
	WorkendTime   time.Time `orm:"auto_now;type(datetime)"`

	Worksigner  string `orm:"size(50)"`
	Workmanager string `orm:"size(50)"`
	Catalog     string `orm:"size(50)"`
}

//WorkOne .
func WorkOne(id string) (*Work, error) {
	o := orm.NewOrm()
	m := Work{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

//WorkOneByID .
func WorkOneByID(id string) (*Work, error) {
	o := orm.NewOrm()
	m := Work{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

//WorkByUserID .
func WorkByUserID(userid string) ([]*Work, int64) {
	query := orm.NewOrm().QueryTable(WorkTBName())
	data := make([]*Work, 0)
	query = query.Filter("UserId", userid)

	total, _ := query.Count()
	query.All(&data)
	return data, total
}

//WorkPageList .
func WorkPageList(params *WorkQueryParam, curUser *User) ([]*Work, int64) {
	query := orm.NewOrm().QueryTable(WorkTBName()).Distinct()
	data := make([]*Work, 0)

	//默认排序
	sortorder := "Id"
	switch params.Sort {
	case "Id":
		sortorder = "Id"
	default:
		sortorder = "WorkstartTime"
	}
	if params.Order != "asc" {
		sortorder = "-" + sortorder
	}
	// if params.SearchLike != "" {
	// 	cond := orm.NewCondition()
	// 	query = query.SetCond(
	// 		cond.AndCond(
	// 			cond.Or(
	// 				"UserId__contains", params.IdLike,
	// 			).Or(
	// 				"Catalog__contains", params.SearchLike,
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
		query = query.Filter("WorkcardtaskUrl", params.SearchLike)
	}

	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	return data, total

}
