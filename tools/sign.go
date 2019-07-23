package tools

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/syyongx/php2go"
	"go-pc_home/pay"
	"io"
	"net/url"
	"sort"
	"strings"
)

type NotifyResponse struct {
	AppID string
	AppKey string
	Result string
	Sign string
}

type QueryResponse struct {
	SerialNumber string
	Platform string
}

type SyncShopLog struct {
	SyncData string
}

type ImgCaptcha struct {}

type ImgCaptchaConsume struct {
	Captcha string  `signName:"captcha"`
	CaptchaTicket string  `signName:"captcha_ticket"`
	ApiTicket string  `signName:"api_ticket"`
}

type SmsCaptcha struct {
	Phone string  `signName:"phone"`
	App string  `signName:"app"`
	AppID int  `signName:"app_id"`
}

type SmsCaptchaConsume struct {
	Captcha string  `signName:"captcha"`
	SmsCaptcha
}

type Kvpair struct {
	k, v string
}

type Kvpairs []Kvpair

func (t Kvpairs) Less(i, j int) bool {
	return t[i].k < t[j].k
}

func (t Kvpairs) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t Kvpairs) Len() int {
	return len(t)
}

func (t Kvpairs) Sort() {
	sort.Sort(t)
}

func (t Kvpairs) RemoveEmpty() (t2 Kvpairs) {
	for _, kv := range t {
		if kv.v != "" {
			t2 = append(t2, kv)
		}
	}
	return
}

func (t Kvpairs) Join() string {
	var strs []string
	for _, kv := range t {
		strs = append(strs, kv.k+"="+kv.v)
	}
	return strings.Join(strs, "&")
}

func Md5Sign(str, key string) string {
	h := md5.New()
	io.WriteString(h, str)
	if key != ""{
		io.WriteString(h, key)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (r NotifyResponse) Md5Sign(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (resp NotifyResponse)CheckSign(c pay.AliPayConfig) bool{
	if resp.Sign == "" {
		return false
	}
	return resp.Sign == resp.GenSign(c)
}

func (r NotifyResponse)GenSign(c pay.AliPayConfig) (sign string){
	p := Kvpairs{}
	if r.Result == "" {
		return sign
	}
	p = append(p, Kvpair{"appid", c.AppID})
	p = append(p, Kvpair{"appkey", c.AppKey})
	p = append(p, Kvpair{"result", r.Result})
	p = p.RemoveEmpty()
	p.Sort()
	sign = r.Md5Sign(p.Join())
	return sign
}

func (r NotifyResponse) DecodeResult() (decodeJsonRes *simplejson.Json, err error){
	return simplejson.NewJson(r.Base64urldecode(r.Result)) //反序列化
}

func (r NotifyResponse) Base64urlEncode(str string) (result string){
	return php2go.Rtrim(php2go.Strtr(Base64_encode([]byte(str)), "+/", "-_"), "=")
}

func (r NotifyResponse) Base64urldecode(str string) (result []byte){
	result,err := base64.RawURLEncoding.DecodeString(str)
	if err != nil{
		return nil
	}
	return result
}

func (q QueryResponse)GenSign() (sign string){
	if q.SerialNumber == "" || q.Platform == ""{
		return sign
	}
	addStr := "DPl96dWEaszJos3d"
	p := Kvpairs{}
	p = append(p, Kvpair{"serialNumber", q.SerialNumber})
	p = append(p, Kvpair{"platform", q.Platform})
	p.Sort()
	sign = Md5Sign(p.Join() + addStr, "")
	return sign
}

func (s SyncShopLog) GenSign() (sign string){
	if s.SyncData == ""{
		return sign
	}
	addStr := "DPl96dWEaszJos3d"
	p := Kvpairs{}
	p = append(p, Kvpair{"sync_data", s.SyncData})
	p.Sort()
	sign = Md5Sign(p.Join() + addStr, "")
	return sign
}


func (i ImgCaptcha) GenSign() (sign string){
	addStr := "DPl96dWEaszJos3d"
	sign = Md5Sign(addStr, "")
	return sign
}

func (c ImgCaptchaConsume) GenSign() (sign string){
	if c.Captcha == "" || c.CaptchaTicket == "" || c.ApiTicket == ""{
		return sign
	}
	addStr := "DPl96dWEaszJos3d"
	p := Kvpairs{}
	p = append(p, Kvpair{"captcha", c.Captcha})
	p = append(p, Kvpair{"captcha_ticket", c.CaptchaTicket})
	p = append(p, Kvpair{"api_ticket", c.ApiTicket})
	p.Sort()
	sign = Md5Sign(p.Join() + addStr, "")
	return sign
}

func (c SmsCaptcha) GenSign() (sign string){
	if c.Phone == "" || c.App == "" || c.AppID == 0{
		return sign
	}
	addStr := "DPl96dWEaszJos3d"
	p := Kvpairs{}
	p = append(p, Kvpair{"phone", c.Phone})
	p = append(p, Kvpair{"app", c.App})
	p = append(p, Kvpair{"app_id", Int2String(c.AppID)})
	p.Sort()
	sign = Md5Sign(p.Join() + addStr, "")
	return sign
}

func (c SmsCaptchaConsume) GenSign() (sign string){
	if c.Captcha == "" || c.Phone == "" || c.App == "" || c.AppID == 0{
		return sign
	}
	addStr := "DPl96dWEaszJos3d"
	p := Kvpairs{}
	p = append(p, Kvpair{"captcha", c.Captcha})
	p = append(p, Kvpair{"phone", c.Phone})
	p = append(p, Kvpair{"app", c.App})
	p = append(p, Kvpair{"app_id", Int2String(c.AppID)})
	p.Sort()
	sign = Md5Sign(p.Join() + addStr, "")
	return sign
}

func VerifySign(c pay.AliPayConfig, u url.Values) (err error) {
	p := Kvpairs{}
	sign := ""
	for k := range u {
		v := u.Get(k)
		switch k {
		case "sign":
			sign = v
			continue
		case "sign_type":
			continue
		}
		p = append(p, Kvpair{k, v})
	}
	if sign == "" {
		err = fmt.Errorf("sign not found")
		return
	}
	p = p.RemoveEmpty()
	p.Sort()
	//fmt.Println(u)
	if Md5Sign(p.Join(), c.Key) != sign {
		err = fmt.Errorf("sign invalid")
		return
	}
	return
}
