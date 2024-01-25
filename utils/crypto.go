package utils

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

// String2md5 -
func String2md5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has)
}

// RandomString -
func RandomString(length int) string {
	// 48 ~ 57 numbers
	// 65 ~ 90 A ~ Z
	// 97 ~ 122 a ~ z
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

// Uuid -
func Uuid() string {
	u2 := uuid.NewV4()
	return fmt.Sprintf("%s", u2)
}
