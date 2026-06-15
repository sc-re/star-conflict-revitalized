package main

import (
	"bytes"
	"context"
	"embed"
	"encoding/xml"
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"strings"

	"starconflict/cmd/http-proxy/adm"

	"github.com/elazarl/goproxy"
)

//go:embed image/*
var f embed.FS

var fortune *bool

func deny(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	return req, goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusForbidden, "Forbidden")
}

type Server struct {
	XMLName        xml.Name `xml:"Server"`
	Name           string   `xml:"name,attr"`
	Address        string   `xml:"address,attr"`
	DefaultDev     string   `xml:"default_dev,attr,omitempty"`
	DefaultLive    string   `xml:"default_live,attr,omitempty"`
	DefaultTest    string   `xml:"default_test,attr,omitempty"`
	DefaultPubtest string   `xml:"default_pubtest,attr,omitempty"`
	DefaultDevtest string   `xml:"default_devtest,attr,omitempty"`
	DefaultV2test  string   `xml:"default_v2test,attr,omitempty"`
}

type Servers struct {
	XMLName xml.Name `xml:"Servers"`
	Servers []Server `xml:"Server"`
}

func addressXmlResponse(req *http.Request) *http.Response {
	servers := Servers{
		Servers: []Server{
			{Name: "VM", Address: "192.168.2.69:4801"},
			{Name: "localhost", Address: "localhost"},
			{Name: "live", Address: "node01im-ru.star-conflict.com:3801;node11sv-ru.star-conflict.com:3801", DefaultLive: "1"},
		},
	}
	xmlDoc, err := xml.Marshal(servers)
	xmlStr := string(xmlDoc)
	if err != nil {
		return goproxy.NewResponse(req, "text/xml; charset=utf-8", http.StatusInternalServerError, "")
	}
	return goproxy.NewResponse(req, "text/xml; charset=utf-8", http.StatusOK, xmlStr)
}

func rssResponse(req *http.Request) *http.Response {
	xmlStr, err := adm.GetRssFeed(*fortune)
	if err != nil {
		return goproxy.NewResponse(req, "text/xml; charset=utf-8", http.StatusInternalServerError, "")
	}
	return goproxy.NewResponse(req, "text/xml; charset=utf-8", http.StatusOK, xmlStr)
}

func promoResponse(req *http.Request) *http.Response {
	xmlStr, err := adm.PromoXml()
	if err != nil {
		return goproxy.NewResponse(req, "text/xml; charset=utf-8", http.StatusInternalServerError, "")
	}
	return goproxy.NewResponse(req, "text/xml; charset=utf-8", http.StatusOK, xmlStr)
}

func binaryResponse(r *http.Request, contentType string, status int, body []byte) *http.Response {
	resp := &http.Response{}
	resp.Request = r
	resp.TransferEncoding = r.TransferEncoding
	resp.Header = make(http.Header)
	resp.Header.Add("Content-Type", contentType)
	resp.StatusCode = status
	resp.Status = http.StatusText(status)
	resp.Proto = "HTTP/1.1"
	resp.ProtoMajor = 1
	resp.ProtoMinor = 1
	buf := bytes.NewBuffer(body)
	resp.ContentLength = int64(buf.Len())
	resp.Body = io.NopCloser(buf)
	return resp
}

func catResponse(req *http.Request) *http.Response {
	body, err := f.ReadFile(req.URL.Path[1:])
	if err != nil {
		log.Printf("Error: %v", err)
		return goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusInternalServerError, "")
	}
	return binaryResponse(req, "image/jpeg", http.StatusOK, body)
}

func handleAdmDst(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	var resp *http.Response
	switch {
	case req.URL.Path == "/addresses.php":
		resp = addressXmlResponse(req)
	case strings.HasPrefix(req.URL.Path, "/rss/"):
		resp = rssResponse(req)
	case strings.HasPrefix(req.URL.Path, "/image/"):
		resp = catResponse(req)
	case strings.HasPrefix(req.URL.Path, "/status/"):

	case strings.HasPrefix(req.URL.Path, "/version/"):

	case strings.HasPrefix(req.URL.Path, "/promocontent/"):
		resp = promoResponse(req)
	default:
		resp = goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusForbidden, "Mööp")
	}

	return req, resp
}

func main() {
	addr := flag.String("addr", ":8080", "Listen address")
	fortune = flag.Bool("fortune", false, "Add random fortune quotes to the rss feed")
	flag.Parse()
	proxy := goproxy.NewProxyHttpServer()

	proxy.Tr.DialContext = func(ctx context.Context, network string, addr string) (net.Conn, error) {
		return nil, nil
	}
	proxy.Verbose = true

	dstAdm := goproxy.DstHostIs("adm.star-conflict.com")
	dstApiGaijin := goproxy.DstHostIs("api.gaijinent.com")
	dstPurch := goproxy.DstHostIs("purch.gaijinent.com")
	dstStore := goproxy.DstHostIs("store.gaijin.net")

	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	proxy.OnRequest(dstAdm).DoFunc(handleAdmDst)
	proxy.OnRequest(dstApiGaijin)
	proxy.OnRequest(dstPurch)
	proxy.OnRequest(dstStore)
	proxy.OnRequest(goproxy.Not(dstAdm), goproxy.Not(dstApiGaijin), goproxy.Not(dstPurch), goproxy.Not(dstStore)).HandleConnect(goproxy.AlwaysReject)
	proxy.OnRequest(goproxy.Not(dstAdm), goproxy.Not(dstApiGaijin), goproxy.Not(dstPurch), goproxy.Not(dstStore)).DoFunc(deny)

	http.ListenAndServe(*addr, proxy)
}
