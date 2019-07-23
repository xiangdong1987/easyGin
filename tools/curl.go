package tools

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"go-pc_home/config"
	"go-pc_home/pay"
	"io/ioutil"
	"net/http"
	"net/url"
)

var APP_chelun = "chelun"
var APPID_pc_home = 82

func HttpPay(orderID int, price int, notifyUrl string) (res *simplejson.Json, err error) {
	params := url.Values{}
	//Url, err := url.Parse(config.GetConfig("pay_host").(string) + "/api.php")
	Url, err := url.Parse(config.GetConfig("pay_host").(string) + "/api.php")
	if err != nil {
		return nil,err
	}
	appID := config.GetConfig("appid").(string)
	platform := "wap"

	params.Set("c", "order")
	params.Set("v", "startpay")
	params.Set("appid", appID)
	params.Set("platform", platform)
	params.Set("order_id", Int2String(orderID))
	params.Set("price", Int2String(price))
	params.Set("notify_url", notifyUrl)
	params.Set("user_id", "0")

	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	//fmt.Println(urlPath)
	resp, err :=   http.Get(urlPath)
	if err != nil {
		return nil,err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil,err
	}

	//fmt.Println(string(body))
	res, err = simplejson.NewJson(body) //反序列化
	if err != nil {
		return nil,err
	}
	return res,nil
}

func HttpQuery(serialNumberStr string, c pay.AliPayConfig) (res *simplejson.Json, err error) {
	params := url.Values{}
	Url := &url.URL{}
	if c.Type=="pc_qr" {
		Url, err = url.Parse(config.GetConfig("chelun_app_http").(string) + "/web/pc_qr_query_order")
	}else{
		Url, err = url.Parse(config.GetConfig("chelun_app_http").(string) + "/web/pc_query_order")
	}

	//Url, err := url.Parse(config.GetConfig("chelun_app_http").(string) + "/api.php")
	if err != nil {
		return nil,err
	}

	params.Set("serialNumber", serialNumberStr)
	params.Set("platform", c.Platform)
	params.Set("sign", QueryResponse{SerialNumber:serialNumberStr,Platform:c.Platform}.GenSign())
	//fmt.Println("HttpQuery:", params)
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	//fmt.Println(urlPath)
	resp, err :=   http.Get(urlPath)
	if err != nil {
		return nil,err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil,err
	}

	//fmt.Println(string(body))
	res, err = simplejson.NewJson(body) //反序列化
	if err != nil {
		return nil,err
	}
	resData, ok := res.CheckGet("data")
	if !ok{
		return nil,err
	}
	return resData,nil
}

func SyncCwShopLog(name string, phone string, nums int, shopId int, addr string, province string, city string, district string, remark string ) bool{
	syncData := make(map[string]interface{})
	syncData["ctime"] = MakeTimestamp()
	syncData["addr"] = addr
	syncData["gold"] = 0
	syncData["nums"] = nums
	syncData["phone"] = phone
	syncData["shop_id"] = shopId
	syncData["name"] = name
	syncData["uid"] = 0
	syncData["province"] = province
	syncData["city"] = city
	syncData["district"] = district
	syncData["remark"] = remark
	//fmt.Println("syncdata:" , syncData)
	params := url.Values{}
	syncDataJson := JsonEncode(syncData)
	params.Set("sync_data", syncDataJson)
	params.Set("sign", SyncShopLog{SyncData:syncDataJson}.GenSign())
	//fmt.Println("SyncCwShopLogparams:" ,params)
	Url, err := url.Parse(config.GetConfig("chelun_app_http").(string) + "/web/pc_shoplog_insert")
	if err != nil {
		return false
	}
	_,err = HttpGet_SimpleJson(Url, params)
	if err !=nil {
		return false
	}
	return true
}

func HttpGet_SimpleJson(Url *url.URL,params url.Values) (res *simplejson.Json, err error){
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	//fmt.Println(urlPath)
	resp, err :=   http.Get(urlPath)
	if err != nil {
		return nil,err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil,err
	}

	//fmt.Println(string(body))
	res, err = simplejson.NewJson(body) //反序列化
	if err != nil {
		return nil,err
	}
	return res,nil
}

func HttpGetImgCaptcha() (res *simplejson.Json, err error) {
	params := url.Values{}
	params.Set("sign", ImgCaptcha{}.GenSign())

	Url, err := url.Parse(config.GetConfig("chelun_app_http").(string) + "/web/pc_img_captcha_get")
	if err != nil {
		return nil,err
	}
	res,err = HttpGet_SimpleJson(Url, params)
	if err !=nil {
		return nil,err
	}
	fmt.Println(res)
	resData, ok := res.CheckGet("data")
	if !ok{
		return nil,err
	}
	return resData,nil
}

func HttpConsumeImgCaptcha(captcha string, ticket string, apiTicket string ) (res *simplejson.Json, err error) {
	params := url.Values{}
	params.Set("captcha", captcha)
	params.Set("captcha_ticket", ticket)
	params.Set("api_ticket", apiTicket)
	params.Set("sign", ImgCaptchaConsume{Captcha:captcha, CaptchaTicket:ticket, ApiTicket:apiTicket}.GenSign())

	Url, err := url.Parse(config.GetConfig("chelun_app_http").(string) + "/web/pc_img_captcha_check")
	if err != nil {
		return nil,err
	}
	res,err = HttpGet_SimpleJson(Url, params)
	if err !=nil {
		return nil,err
	}
	fmt.Println(res)
	resData, ok := res.CheckGet("data")
	if !ok{
		return nil,err
	}
	return resData,nil
}

func HttpGetSmsCaptcha(phone string, app string, appID int) (res *simplejson.Json, err error) {
	params := url.Values{}
	params.Set("phone", phone)
	params.Set("app", app)
	params.Set("app_id", Int2String(appID))
	params.Set("sign", SmsCaptcha{Phone:phone,App:app,AppID:appID}.GenSign())
	Url, err := url.Parse(config.GetConfig("chelun_app_http").(string) + "/web/pc_sms_captcha_get")
	if err != nil {
		return nil,err
	}
	res,err = HttpGet_SimpleJson(Url, params)
	fmt.Println(res)
	if err !=nil {
		return nil,err
	}
	resData, ok := res.CheckGet("data")
	if !ok{
		return nil,err
	}
	return resData,nil
}

func HttpConsumeSmsCaptcha(captcha string, phone string, app string, appID int ) (res *simplejson.Json, err error) {
	params := url.Values{}
	params.Set("captcha", captcha)
	params.Set("phone", phone)
	params.Set("app", app)
	params.Set("app_id", Int2String(appID))
	params.Set("sign", SmsCaptchaConsume{Captcha:captcha,SmsCaptcha:SmsCaptcha{Phone:phone,App:app,AppID:appID}}.GenSign())

	Url, err := url.Parse(config.GetConfig("chelun_app_http").(string) + "/web/pc_sms_captcha_check")
	if err != nil {
		return nil,err
	}
	res,err = HttpGet_SimpleJson(Url, params)
	fmt.Println(res)
	if err !=nil {
		return nil,err
	}
	resData, ok := res.CheckGet("data")
	if !ok{
		return nil,err
	}
	return resData,nil
}