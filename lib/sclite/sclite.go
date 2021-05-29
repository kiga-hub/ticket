package sclite

import (
	"time"
	"github.com/astaxie/beego"
	"path/filepath"
	"fmt"
	"os"
	"io"
	"bufio"
	"strings"
	"regexp"
	"strconv"
	"Two-Card/utils"
	"Two-Card/lib/extend"
	"errors"
)

type Sclite struct {
	Pra 	[]*SclitePra
	Snt     int 		//总句数
	Chr		int  		//总词数
	Corr	float64		// 正确率
	Sub		float64		// 代替错误率
	Del		float64		// 删除错误率
	Ins		float64 	// 插入错误率
	Err 	float64		// 总错误率
}

type SclitePra struct {
	Content	string
	Corr	int
	Sub		int
	Del		int
	Ins		int
	ValidOk	bool
	Ok 		bool
}

type AnalysisLine func([]byte) *Sclite

func IsExist(f string) bool {
    _, err := os.Stat(f)
    return err == nil || os.IsExist(err)
}

func CreateScliteFile(text, file string) error {
	f, err := os.Create(file) 
	if err != nil {
		return err
	}
	_, err = io.WriteString(f, fmt.Sprintf("%s\n", text)) 
	if err != nil {
		return err
	}
	return nil
}

func processLine(line []byte) *Sclite {
	if !strings.Contains(string(line), "Sum/Avg") {
		return nil
	}
	r, _ := regexp.Compile(`\d+\.?\d*`)
	res := r.FindAllString(string(line), 8)
	if len(res) != 8 {
		return nil
	}
	var s Sclite
	s.Snt, _ = strconv.Atoi(res[0])
	s.Chr, _ = strconv.Atoi(res[1])
	s.Corr, _ =  strconv.ParseFloat(res[2], 64)
	s.Sub, _ =  strconv.ParseFloat(res[3], 64)
	s.Del, _ =  strconv.ParseFloat(res[4], 64)
	s.Ins, _ =  strconv.ParseFloat(res[5], 64)
	s.Err, _ =  strconv.ParseFloat(res[6], 64)
	return &s
}
   
func ReadLine(filePth string, hookfn AnalysisLine) (*Sclite, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}
	defer f.Close()
   
	bfRd := bufio.NewReader(f)
	for {
		line, err := bfRd.ReadBytes('\n')
		s := hookfn(line)
		if s != nil {
			return s, nil
		}
		if err != nil {
			return nil, err
		}
	}
}

func AnalysisSclitePra(hyppath string) (*SclitePra, error) {
	prapath := fmt.Sprintf("%s.pra", hyppath)
	if !IsExist(prapath) {
		return nil, errors.New(fmt.Sprintf("%s not find", prapath))
	}

	f, err := os.Open(prapath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
   
	var result SclitePra
	var ref []string
	var hyp []string
	bfRd := bufio.NewReader(f)
	for {
		line, err := bfRd.ReadBytes('\n')
		if err != nil {
			break
		}
		data := strings.Replace(string(line), "\n", "", -1)
		// 解析Scores
		if strings.Contains(data, "Scores:") {
			r, _ := regexp.Compile(`\d+`)
			res := r.FindAllString(data, 4)
			if len(res) != 4 {
				return nil, errors.New(data)
			}
			result.Corr, _ = strconv.Atoi(res[0])
			result.Sub, _ = strconv.Atoi(res[1])
			result.Del, _ = strconv.Atoi(res[2])
			result.Ins, _ = strconv.Atoi(res[3])
		}
		// 解析ref
		if strings.Contains(data, "REF:") {
			ref_tmp := strings.Split(data[6:], " ")	
			for _, s := range ref_tmp {
				if s == "" {
					continue
				}
				ref = append(ref, s)
			}
		}
		// 解析hyp
		if strings.Contains(data, "HYP:") {
			hyp_tmp := strings.Split(data[6:], " ")	
			for _, s := range hyp_tmp {
				if s == "" {
					continue
				}
				hyp = append(hyp, s)
			}
		}
	}
	os.Remove(prapath)
	if len(ref) == 0 || len(ref) != len(hyp) {
		return nil, errors.New(fmt.Sprintf("%s(%d)|%s(%d)", strings.Join(ref, " "), len(ref), strings.Join(hyp, " "), len(hyp)))
	}
	var content []string
	for i, c := range ref {
		if i >= result.Corr + result.Sub + result.Del + result.Ins {
			break
		}
		if c == "***" {
			// 插入错误
			content = append(content, fmt.Sprintf("<i>%s</i>", hyp[i]))
		} else if hyp[i] == "***" {
			// 删除错误
			content = append(content, fmt.Sprintf("<b>%s</b>", c))
		} else if c != hyp[i] {
			// 修改错误
			content = append(content, fmt.Sprintf("<i>%s</i><b>%s</b>", hyp[i], c))
		} else {
			content = append(content, c)
		}
	}
	result.Content = strings.Join(content, "")
	result.Ok = true
	result.ValidOk = true
	return &result, nil
}


func AnalysisScliteSys(hyppath string) *Sclite {
	syspath := fmt.Sprintf("%s.sys", hyppath)
	if !IsExist(syspath) {
		return nil
	}

	s, err := ReadLine(syspath, processLine)
	os.Remove(syspath)
	if err != nil {
		return &Sclite{}
	}
	return s
}

func ScliteText(reftext, hyptext string) (*SclitePra, error ){
	if reftext == "" || hyptext == "" {
		return nil, errors.New("sclite null string")
	}
	random_suffix := utils.RandomString(10)
	time_suffix := time.Now().Format("20060102150405")
	DataRoot := beego.AppConfig.String("dataroot")
	ScliteDir := filepath.Join(DataRoot, "Sclite")
	if !IsExist(ScliteDir) {
		if err := os.MkdirAll(ScliteDir, 0775); err != nil {
			return nil, err
		}
	}
	refname := fmt.Sprintf("%s_%s.ref", time_suffix, random_suffix)
	refpath := filepath.Join(ScliteDir, refname)
	if err := CreateScliteFile(reftext, refpath); err != nil {
		return nil, err
	}

	hypname := fmt.Sprintf("%s_%s.hyp", time_suffix, random_suffix)
	hyppath := filepath.Join(ScliteDir, hypname)
	if err := CreateScliteFile(hyptext, hyppath); err != nil {
		return nil, err
	}

	utils.LogDebug(fmt.Sprintf("sclite[%s][%s]", reftext, hyptext))
	extend.ExecSclite(refpath, hyppath, ScliteDir)

	os.Remove(refpath)
	os.Remove(hyppath)
	os.Remove(fmt.Sprintf("%s.raw", hyppath))
	os.Remove(fmt.Sprintf("%s.sys", hyppath))
	return AnalysisSclitePra(hyppath)
}

func StringReplace(src string) string {
	// 去掉回车
	dst := strings.Replace(src, "\n", "", -1)
	// 去掉换行
	dst = strings.Replace(dst, "\r", "", -1)
	// 去掉TAB
	dst = strings.Replace(dst, "\t", "", -1)
	// 去掉空格
	dst = strings.Replace(dst, " ", "", -1)
	return dst
}

func ScliteTextGrid(refgrid, hypgrid interface{}, offset float64) *Sclite {
	var corr []*SclitePra
	findIndex := -1
	for _, ref := range refgrid.([]interface{}) {
		r := ref.(map[string]interface{})
		find := false
		pra := new(SclitePra)

		r_str := StringReplace(r["content"].(string))

		for i, hyp := range hypgrid.([]interface{}) {
			// 从上一个找到的分段向后找
			if i <= findIndex {
				continue
			}

			// 寻找对应比对段
			h := hyp.(map[string]interface{})
			if h["seconds"].(float64) < r["seconds"].(float64) - offset || h["seconds"].(float64) > r["seconds"].(float64) + offset {
				continue
			}

			find = true
			findIndex = i
			pra.Ok = true // 时间分隔线有效
			pra.ValidOk = true // 默认有效性正确

			h_str := StringReplace(h["content"].(string))

			if r["valid"].(string) != h["valid"].(string) {
				pra.ValidOk = false // 有效性不一致
				if r["valid"].(string) != "0" {
					pra.Del, _ = GetWordCount(r_str) // 全部删除错误
				} else {
					pra.Ins, _ = GetWordCount(h_str) // 全部插入错误
				}
				pra.Content = fmt.Sprintf("<i>%s</i><b>%s<b>", h_str, r_str)
			} else if r["valid"].(string)  == "0" {
				// 无效段比对
				if h_str != r_str {
					// pra.Sub = 1 // 全部修改错误
					pra.Content = fmt.Sprintf("<i>%s</i><b>%s<b>", h_str, r_str)
				} else {
					// pra.Corr = 1 // 全部正确
					pra.Content = r_str
				}
			} else if r_str == h_str  {
				pra.Corr, _ = GetWordCount(r_str) // 全部正确
				pra.Content = r_str
			} else if r_str == "" {
				pra.Ins, _ = GetWordCount(h_str) // 全部插入错误
				pra.Content = fmt.Sprintf("<i>%s</i>", h_str)
			} else if h_str == "" {
				pra.Del, _ = GetWordCount(r_str) // 全部删除错误
				pra.Content = fmt.Sprintf("<b>%s</b>", r_str)
			} else {
				// 有效段判断
				ppra, err := ScliteText(r_str, h_str)
				if err != nil {
					utils.LogError(err)
					pra.Sub, _ = GetWordCount(r_str) // 全部修改错误
					pra.Content = fmt.Sprintf("<i>%s</i><b>%s<b>", h_str, r_str)
				} else {
					pra = ppra
				}
			}
			break
		}
		if !find {
			if r["valid"].(string)  == "0" {
				// pra.Del = 1 // 全部删除错误
				pra.Content = fmt.Sprintf("<b>%s</b>", r_str)
			} else {
				pra.Del, _ = GetWordCount(r_str) // 全部删除错误
				pra.Content = fmt.Sprintf("<b>%s</b>", r_str)
			}
		}
		corr = append(corr, pra)
	}
	if corr == nil {
		return &Sclite{Corr: 100}
	}
	s := Sclite{
		Pra : corr,
		Snt: 1,
	}

	var Corr int
	var Sub int
	var Del int
	var Ins int
	for _, c := range corr {
		Corr += c.Corr
		Sub += c.Sub
		Del += c.Del
		Ins += c.Ins
	}
	Chr := Corr + Sub + Del
	if Chr > 0 {
		s.Sub = float64(Sub)*100 / float64(Chr)
		s.Del = float64(Del)*100 / float64(Chr)
		s.Ins = float64(Ins)*100 / float64(Chr)
		s.Err = float64(Ins+Sub+Del)*100 / float64(Chr)
		s.Chr = Chr
	} else if Ins > 0 {
		s.Ins = 100
		s.Err = 100
	}
	s.Corr = 100 - s.Err
	return &s
}

func GetWordCount(ref string) (int, error) {
	str := StringReplace(ref)
	if len(str) == 0 {
		return 0, nil
	}
	ret, err := ScliteText(str, str)
	if err != nil {
		return 0, err
	}
	return ret.Corr + ret.Sub + ret.Del, nil
}