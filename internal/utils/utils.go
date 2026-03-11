package utils

import (
	"net"
	"net/url"

	"github.com/yylego/must"
)

// ExtractPort extracts port number from URL endpoint
// ExtractPort 从 URL 端点提取端口号
func ExtractPort(endpoint *url.URL) string {
	must.Full(endpoint)
	_, port, _ := net.SplitHostPort(must.Nice(endpoint.Host))
	return must.Nice(port)
}
