package txai

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/axgle/mahonia"
)

// http 请求

// HttpPost 发起POST请求
func (ai *TxAi) HttpPost(url string, u url.Values) (respBody []byte, err error) {
	contentType := "application/x-www-form-urlencoded"
	client := http.Client{Timeout: ai.timeout}
	request, err := http.NewRequest("POST", url, strings.NewReader(u.Encode()))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", contentType)
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	respBody, err = ioutil.ReadAll(resp.Body)

	return respBody, err
}

// RequestAPI 发起请求，并解析响应结果
func (ai *TxAi) RequestAPI(uri string, u url.Values, response BaseResponseInterface) error {
	apiURL := BaseURL + uri
	if ai.debug == true {
		log.Println("Request-URI:", apiURL)
		log.Println("Request-Body:", u)
	}
	// 发起请求
	body, err := ai.HttpPost(apiURL, u)
	if err != nil {
		return err
	}
	if ai.debug == true {
		log.Println("Response-Body", string(body))
	}
	// log.Println(u)
	// log.Println(string(body))
	// 判断是否是gbk参数的接口
	switch uri {
	case "/nlp/nlp_wordseg", "/nlp/nlp_wordpos", "/nlp/nlp_wordner", "/nlp/nlp_wordsyn":
		enc := mahonia.NewDecoder("gbk")
		body = []byte(enc.ConvertString(string(body)))
		break
	}
	// 解析json数据
	err = json.Unmarshal(body, response)
	if err != nil {
		return err
	}
	ret := response.GetRet()
	if ret != 0 {
		return GetError(ret)
	}
	return nil
}
