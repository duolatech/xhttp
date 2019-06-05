package xhttp

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpRequest struct {
	Link     string
	Referer  string
	Header   map[string]string
	Cookie   map[string]string
	Timeout  map[string]time.Duration
	Proxy    string
	Response *http.Response
	Error    error
	Reqtime  time.Duration
}

func NewHttp() *HttpRequest {
	return &HttpRequest{}
}

//设置页面referer
func (h *HttpRequest) SetReferer(referer string) {
	h.Referer = referer
}

//设置header
func (h *HttpRequest) SetHeader(header map[string]string) {
	h.Header = header
}

//设置cookie
func (h *HttpRequest) SetCookie(cookie map[string]string) {
	h.Cookie = cookie
}

//设置代理
func (h *HttpRequest) SetProxy(proxy string) {
	h.Proxy = proxy
}

//设置超时时间, 单位s
func (h *HttpRequest) SetTimeout(connectTimeout, readWriteTimeout time.Duration) {
	h.Timeout = map[string]time.Duration{
		"connectTimeout":   time.Second * time.Duration(connectTimeout),
		"readWriteTimeout": time.Second * time.Duration(readWriteTimeout),
	}
}

//get 请求
func (h *HttpRequest) Get(url string) *HttpRequest {

	h.Link = url
	data := map[string]string{}
	resp, err := h.HttpClient("GET", data)

	h.Response = resp
	h.Error = err
	return h
}

//post 请求
func (h *HttpRequest) Post(url string, data map[string]string) *HttpRequest {

	h.Link = url
	resp, err := h.HttpClient("POST", data)

	h.Response = resp
	h.Error = err
	return h
}

//put 请求
func (h *HttpRequest) Put(url string, data map[string]string) *HttpRequest {

	h.Link = url
	resp, err := h.HttpClient("PUT", data)

	h.Response = resp
	h.Error = err
	return h
}

//delete 请求
func (h *HttpRequest) Delete(url string, data map[string]string) *HttpRequest {

	h.Link = url
	resp, err := h.HttpClient("DELETE", data)

	h.Response = resp
	h.Error = err
	return h
}

//获取请求内容
func (h *HttpRequest) GetContent() ([]byte, error) {

	respBody, err := ioutil.ReadAll(h.Response.Body)

	return respBody, err
}

//获取响应header
func (h *HttpRequest) GetHeader() (header http.Header) {
	header = h.Response.Header
	return
}

//获取文档类型
func (h *HttpRequest) GetContentType() (ctype string) {
	ctype = h.Response.Header.Get("Content-Type")
	return
}

//获取响应cookie
func (h *HttpRequest) GetCookies() []*http.Cookie {
	return h.Response.Cookies()
}

//获取响应状态码
func (h *HttpRequest) GetStatudCode() int {
	return h.Response.StatusCode
}

//获取请求消耗时间
func (h *HttpRequest) GetTime() time.Duration {
	return h.Reqtime
}

//设置连接超时
func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout) //设置建立连接超时
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(rwTimeout)) //设置发送接受数据超时
		return conn, nil
	}
}

func (h *HttpRequest) HttpClient(method string, data map[string]string) (*http.Response, error) {

	t1 := time.Now()
	transport := &http.Transport{}
	//默认连接和处理超时时间均为60s
	if len(h.Timeout) == 0 {
		h.SetTimeout(60, 60)
		transport.DialTLS = TimeoutDialer(h.Timeout["connectTimeout"], h.Timeout["readWriteTimeout"])
	}
	//设置代理
	if len(h.Proxy) > 0 {
		proxy, _ := url.Parse(h.Proxy)
		transport.Proxy = http.ProxyURL(proxy)
	}
	//不校验服务器证书
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	c := http.Client{
		Transport: transport,
	}

	//参数设置
	param := url.Values{}
	for k, v := range data {
		param[k] = []string{v}
	}

	req, err := http.NewRequest(method, h.Link, strings.NewReader(param.Encode()))

	if err != nil {
		fmt.Println("req error:" + err.Error())
		return nil, err
	}
	//设置referer
	if len(h.Referer) > 0 {
		req.Header.Set("Referer", h.Referer)
	}

	//设置header
	for headername, headervalue := range h.Header {
		req.Header.Add(headername, headervalue)
	}
	//设置cookie
	cookie := &http.Cookie{}
	for cookiename, cookievalue := range h.Cookie {
		cookie.Name = cookiename
		cookie.Value = cookievalue
		req.AddCookie(cookie)
	}

	resp, err := c.Do(req)
	t2 := time.Now()
	if err != nil {
		fmt.Println("do error,err:" + err.Error())
		return nil, err
	}
	h.Reqtime = t2.Sub(t1)
	return resp, err
}
