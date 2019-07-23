package tools

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var picDomain = "http://picture.eclicks.cn/"

func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

/**
生成返回给客户端的response
*/
func GetResult(code int, msg string, data interface{}) Result {
	resp := Result{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	return resp
}

func JsonEncode(data interface{}) string {
	//fmt.Println(data)
	jsonBytes, err := json.Marshal(data)
	//fmt.Println(jsonBytes)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

//加图片域名
func HandlePicUrl(path string) string {
	if path != "" {
		path = picDomain + path
	}
	return path
}

//中文截串
func SubMbStr(str string, limit int) string {
	if len([]rune(str)) < limit {
		limit = len([]rune(str))
	}
	str = string([]rune(str)[:limit])
	return str
}

//去除htnl标签
func TrimHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")
	return strings.TrimSpace(src)
}

//类似 php array_column
func SliceColumn(structSlice []interface{}, key string) []interface{} {
	rt := reflect.TypeOf(structSlice)
	rv := reflect.ValueOf(structSlice)
	if rt.Kind() == reflect.Slice { //切片类型
		var sliceColumn []interface{}
		elemt := rt.Elem() //获取切片元素类型
		for i := 0; i < rv.Len(); i++ {
			inxv := rv.Index(i)
			if elemt.Kind() == reflect.Struct {
				for i := 0; i < elemt.NumField(); i++ {
					if elemt.Field(i).Name == key {
						strf := inxv.Field(i)
						switch strf.Kind() {
						case reflect.String:
							sliceColumn = append(sliceColumn, strf.String())
						case reflect.Float64:
							sliceColumn = append(sliceColumn, strf.Float())
						case reflect.Int, reflect.Int64:
							sliceColumn = append(sliceColumn, strf.Int())
						default:
							//do nothing
						}
					}
				}
			}
		}
		return sliceColumn
	}
	return nil
}

func ExplodeInt(str string, split string) []int {
	strArray := strings.Split(str, split)
	var result []int
	for _, value := range strArray {
		newValue, error := strconv.Atoi(value)
		if error != nil {
			fmt.Println("字符串转换成整数失败")
		}
		result = append(result, newValue)
	}
	return result
}
func ExplodeStr(str string, split string) []string {
	strArray := strings.Split(str, split)
	return strArray
}
func StrToTime(str string) int64 {
	formatTime, err := time.Parse("2006-01-02 15:04:05", str)
	if err == nil {
		fmt.Println(formatTime)
	}
	return formatTime.Unix()
}

func Int64TOStr(num int64) string {
	return strconv.FormatInt(num, 10)
}

func Int64TOInt(num int64) int {
	return String2Int(Int64TOStr(num))
}

func Implode(list []string, split string) string {
	return strings.Join(list, split)
}

func ValidatePhone(phone string) bool {
	reg := `^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(phone)
}

func DebugType(val interface{}) {
	fmt.Printf("v1 type:%T\n", val)
}

func String2Int(str string) int {
	result, err := strconv.Atoi(str)
	if err != nil {
		println(err.Error())
	}
	return result
}

func Int2String(val int) string {
	result := strconv.Itoa(val)
	return result
}

func Now() int {
	return int(time.Now().Unix())
}

func String2Int64(str string) int64 {
	result, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		println(err.Error())
	}
	return result
}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func Base64_encode(input []byte) string {
	encodeString := base64.StdEncoding.EncodeToString(input)
	return encodeString
}

/*func Base64_decode(encodeStr string) []byte{
	encodeStr = strings.Replace(encodeStr, " ", "", -1)
	decodeBytes, err := base64.StdEncoding.DecodeString(encodeStr)
	fmt.Println(decodeBytes)
	fmt.Println(err)
	if err != nil {
		return nil
	}
	return decodeBytes
}

func Byte2Str(b []byte) string{
	return string(b[:])
}*/

//获取结构体中字段的名称
func GetFieldName(structName interface{}) []string {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return nil
	}
	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		result = append(result, t.Field(i).Name)
	}
	return result
}

//获取结构体中Tag的值，如果没有tag则返回字段值
func GetTagName(structName interface{}, getTagName string) []string {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return nil
	}
	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		tagName := t.Field(i).Name
		if getTagName != "" {
			tagName = t.Field(i).Tag.Get(getTagName)
			/*tags := strings.Split(string(t.Field(i).Tag), "\"")
			if len(tags) > 1 {
				tagName = tags[1]
			}*/
		}
		result = append(result, tagName)
	}
	return result
}

func IsMach(reg string, str string) bool {
	ruleNum, _ := regexp.Compile(reg)
	return ruleNum.MatchString(str)
}

func HandlePhone(phone string) string{
	if phone == "" || len(phone) < 7{
		return ""
	}
	return phone[0:3]+"****"+phone[7:]
}
