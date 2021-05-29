/*
 * @Author: your name
 * @Date: 2020-03-23 16:47:39
 * @LastEditTime: 2020-03-23 17:59:19
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \Two-Card\utils\sort.go
 */
package utils

import (
	"sort"
)

type TaskPool struct {
	Id               string
	Name             string
	Type             string
	DeliveryDatetime string
	CreateDatetime   string
	Owner            string
	Count            int
	Tags             []string
}

type TaskPools []*TaskPool

func (s TaskPools) Len() int {
	return len(s)
}

func (s TaskPools) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type Desc struct{ TaskPools }

func (s Desc) Less(i, j int) bool {
	return s.TaskPools[i].Count > s.TaskPools[j].Count
}

type Asc struct{ TaskPools }

func (s Asc) Less(i, j int) bool {
	return s.TaskPools[i].Count < s.TaskPools[j].Count
}

func TaskPoolSort(rows []interface{}, s string) []interface{} {
	var pools TaskPools
	for _, row := range rows {
		data := row.(map[string]interface{})
		pools = append(pools, &TaskPool{
			Id:               data["Id"].(string),
			Name:             data["Name"].(string),
			Type:             data["Type"].(string),
			DeliveryDatetime: data["DeliveryDatetime"].(string),
			CreateDatetime:   data["CreateDatetime"].(string),
			Owner:            data["Owner"].(string),
			Count:            data["Count"].(int),
			Tags:             data["Tags"].([]string),
		})
	}

	if s == "asc" {
		sort.Sort(Asc{pools})
	} else {
		sort.Sort(Desc{pools})
	}

	var result []interface{}
	for _, dt := range pools {
		result = append(result, map[string]interface{}{
			"Id":               dt.Id,
			"Name":             dt.Name,
			"Type":             dt.Type,
			"DeliveryDatetime": dt.DeliveryDatetime,
			"CreateDatetime":   dt.CreateDatetime,
			"Owner":            dt.Owner,
			"Count":            dt.Count,
			"Tags":             dt.Tags,
		})
	}
	return result
}
