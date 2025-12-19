package models

import (
	"strings"

	"github.com/ManoloEsS/burrow_prototype/internal/config"
)

type Request struct {
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	ContentType string            `json:"content-type,omitempty"`
	Body        string            `json:"body,omitempty"`
	Params      map[string]string `json:"params,omitempty"`
	Auth        map[string]string `json:"auth,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
}

func (req *Request) ParseMethod(method string) {
	correctMethod := strings.ToUpper(method)
	req.Method = correctMethod
}

func (req *Request) ParseUrl(cfg *config.Config, url string) {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		req.URL = url
		return
	}
	if url == cfg.DefaultPort || url == "" {
		req.URL = "http://localhost" + cfg.DefaultPort
		return
	}

	req.URL = "https://" + url
}

func (req *Request) ParseHeaders(headersStr string) {
	if req.Headers == nil {
		req.Headers = make(map[string]string)
	}
	headers := strings.Fields(headersStr)
	for _, h := range headers {
		parsedHeader := strings.Split(h, ":")
		req.Headers[parsedHeader[0]] = parsedHeader[1]
	}
}

func (req *Request) ParseBody(body string) {
	req.Body = body
}

func (req *Request) ParseAuth(auth string) {
	if req.Auth == nil {
		req.Auth = make(map[string]string)
	}
	req.Auth[auth] = "placeholder"
}

func (req *Request) ParseParams(paramsStr string) {
	if req.Params == nil {
		req.Params = make(map[string]string)
	}
	params := strings.Fields(paramsStr)
	for _, p := range params {
		parsedParams := strings.Split(p, ":")
		req.Params[parsedParams[0]] = parsedParams[1]
	}
}
