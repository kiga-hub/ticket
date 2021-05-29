package models

import (
	_ "Two-Card/enums"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "strings"
)

// init 初始化
func init() {
	orm.RegisterModel(new(Voice))
}

func VoiceTBName() string {
	table := "voice"
	return table
}

// TableName 设置User表名
func (a *Voice) TableName() string {
	return VoiceTBName()
}

// UserQueryParam 用于查询的类
type VoiceQueryParam struct {
	BaseQueryParam
	SearchLike string //模糊查询
	IdLike     string
	TypeLike   int
}

//Voice 实体类 .
type Voice struct {
	Id         string `orm:"size(50);pk"`
	UserId     string `orm:"size(50)"`
	WorkId     string `orm:"size(50)"`
	OperateId  string `orm:"size(50)"`
	VoiceUrl   string `orm:"size(255)"`
	VideoUrl   string `orm:"size(50)"`
	CardType   int
	CardOption string `orm:"size(50)"`
}

//VoiceOne .
func VoiceOne(id string) (*Voice, error) {
	o := orm.NewOrm()
	m := Voice{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

//VoicesByUserID .
func VoicesByUserID(userid string) []*Voice {
	query := orm.NewOrm().QueryTable(VoiceTBName())
	data := make([]*Voice, 0)
	query = query.Filter("User__Id", userid)
	query.OrderBy("Xmin").All(&data)
	return data
}

func OperateByUserIDAndOption(userid, option string) *Voice {
	query := orm.NewOrm().QueryTable(OperateTBName())
	data := Voice{}
	query = query.Filter("UserId", userid)
	query = query.Filter("CardOption", option)
	//total ,_:= query.Count()
	query.All(&data)
	return &data
}

func WorkVoicePageList(params *VoiceQueryParam) ([]*Voice, int64) {
	query := orm.NewOrm().QueryTable(VoiceTBName()).Distinct()
	data := make([]*Voice, 0)

	//默认排序
	sortorder := "Id"
	switch params.Sort {
	case "Id":
		sortorder = "Id"
	default:
		sortorder = "WorkId"
	}
	if params.Order != "asc" {
		sortorder = "-" + sortorder
	}
	fmt.Println("ewqewewqweewe", params.TypeLike)
	if params.SearchLike != "" {
		cond := orm.NewCondition()
		query = query.SetCond(
			cond.AndCond(
				cond.Or(
					"UserId__icontains", params.SearchLike,
				).Or(
					"WorkId__icontains", params.SearchLike,
				).Or(
					"CardOption__icontains", params.SearchLike,
				),
			),
		)
	}
	if params.TypeLike != -1 {
		query = query.Filter("Cardtype", params.TypeLike)
	}
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	return data, total
}
