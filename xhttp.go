package xhttp

import "time"

type HttpRequest struct {
	Link    string
	Method  string
	Referer string
	Header  map[string]string
	Cookie  map[string]string
	Param   map[string]string
	Timeout time.Duration
	Proxy   string
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

//get请求
func (h *HttpRequest) Get(url string) {

}
