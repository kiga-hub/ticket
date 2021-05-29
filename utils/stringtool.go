package utils

import "fmt"
import "crypto/md5"
import "math/rand"
import "math"
import "time"
import "os"
import "strconv"
import "strings"
import "path/filepath"
import "path"
import "github.com/satori/go.uuid"
import "io/ioutil"

//将字符串加密成 md5
func String2md5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has) //将[]byte转成16进制
}

//RandomString 在数字、大写字母、小写字母范围内生成num位的随机字符串
func RandomString(length int) string {
	// 48 ~ 57 数字
	// 65 ~ 90 A ~ Z
	// 97 ~ 122 a ~ z
	// 一共62个字符，在0~61进行随机，小于10时，在数字范围随机，
	// 小于36在大写范围内随机，其他在小写范围随机
	rand.Seed(time.Now().UnixNano())
	result := make([]string, 0, length)
	for i := 0; i < length; i++ {
		t := rand.Intn(62)
		if t < 10 {
			result = append(result, strconv.Itoa(rand.Intn(10)))
		} else if t < 36 {
			result = append(result, string(rand.Intn(26)+65))
		} else {
			result = append(result, string(rand.Intn(26)+97))
		}
	}
	return strings.Join(result, "")
}

// 随机文件名
func RandomFilename(fullname string, prefix string, length int) string {
	paths, fileName := filepath.Split(fullname)
	suffix := path.Ext(fileName)
	currentTime := time.Now().Format("20060102150405")
	var newName string
	if prefix != "" {
		newName = fmt.Sprintf("%s_%s_%s%s", prefix, currentTime, RandomString(length), suffix)
	} else {
		newName = fmt.Sprintf("%s_%s%s", currentTime, RandomString(length), suffix)
	}
	return filepath.Join(paths, newName)
}

// 获取同名wav文件
func WaveFilename(fullname string) string {
	paths, fileName := filepath.Split(fullname)
	suffix := path.Ext(fileName)
	name := fileName[:len(fileName)-len(suffix)]
	newFileName := fmt.Sprintf("%s.wav", name)
	return filepath.Join(paths, newFileName)
}

// 字符串转时间
func StrToTime(strTime string) (time.Time, error) {
	local, _ := time.LoadLocation("Local")
	return time.ParseInLocation("2006-01-02 15:04:05", strTime, local)
}

// 日期字符转时间
func StrToDate(strTime string) (time.Time, error) {
	local, _ := time.LoadLocation("Local")
	return time.ParseInLocation("2006-01-02", strTime, local)
}

// 时间格式化
func TimeToStr(Time time.Time) string {
	if Time.Unix() > 0 {
		return Time.Format("2006-01-02 15:04:05")
	}
	return ""
}

// 时间转日期
func DateToStr(Time time.Time) string {
	if Time.Unix() > 0 {
		return Time.Format("2006-01-02")
	}
	return ""
}

// 获取uuid
func Uuid() string {
	u2 := uuid.NewV4()
	return fmt.Sprintf("%s", u2)
}

// 保留小数点后3位
func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", value), 64)
	return value
}

func SecondFormat(second float64) string {
	msstr := strings.Split(fmt.Sprintf("%.03f", second), ".")[1]
	s := int(math.Floor(second))
	sstr := fmt.Sprintf("%.02d", s%60)
	mstr := fmt.Sprintf("%.02d", (s/60)%60)
	hstr := fmt.Sprintf("%.02d", s/3600)
	return fmt.Sprintf("%s:%s:%s.%s", hstr, mstr, sstr, msstr)
}

func CreateDir(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

func WriteFile(name, content string) error {
	data := []byte(content)
	return ioutil.WriteFile(name, data, 0644)
}
