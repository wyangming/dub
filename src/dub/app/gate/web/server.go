package web

import (
	"dub/common"
	"dub/config"
	"dub/define"
	"dub/frame"
	"dub/utils"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	json "github.com/json-iterator/go"
)

type GateWebServer struct {
	gwsCfg    define.GateWebServerConfig               //web网关服务配置
	logCfg    define.LogConfig                         //日志配置
	regConn   common.IConnector                        //与注册服务器连接
	log       *utils.Logger                            //日志对象
	proxyUrls map[string]*define.ModelRegReqServerType //代理的路径与服务器映射
}

func (g *GateWebServer) Init(cfgPath string) {
	//读取配置
	var err error
	g.gwsCfg, g.logCfg, err = config.GetGateWebServerConfig(cfgPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	//日志配置
	log := utils.NewLogger()
	log.SetConfig(g.logCfg)
	g.log = log

	//连接主服务器
	conn := frame.NewConnector()
	conn.CallBack = g.RegServiceCallBack
	g.regConn = conn
	err = g.regConn.Start(g.gwsCfg.RegAddr)
	if err != nil {
		log.Errorf("server.go Init method g.regCon.Start err %v\n", err)
		os.Exit(2)
	}
	g.Reg()
	log.Infof("reg %s server is %s\n", g.logCfg.DeviceName, g.gwsCfg.Addr)

	g.proxyUrls = make(map[string]*define.ModelRegReqServerType)
	//启动代理服务器
	g.StartProxy()
}

//启动代理服务器的方法
func (g *GateWebServer) StartProxy() {
	for {
		if len(g.proxyUrls) < 1 {
			continue
		}

		//第一种代理方案有问题
		err := http.ListenAndServe(g.gwsCfg.Addr, g.reverseProxy())

		if err != nil {
			g.log.Errorf("server.go Init method start web gate server err %v\n", err)
		}
		g.log.Infoln("web gate server start")
	}
}

//代理方案实现的方法
func (g *GateWebServer) reverseProxy() *httputil.ReverseProxy {
	director := func(req *http.Request) {
		//处理找不到页面报错的问题
		defer func() {
			if err := recover(); err != nil {
				g.log.Errorln("server.go reverseProxy method err recover() do. err is %v", err)
			}
		}()

		//TODO:处理一下找不到代理的问题
		req_url := req.URL.Path
		proxy_url := "/"
		proxy_server := g.proxyUrls[proxy_url]

		req_url_splits := strings.Split(req_url, "/")
		if len(req_url_splits) > 1 && req_url_splits[1] != "" {
			first_url := req_url_splits[1]
			for url_key, url_val := range g.proxyUrls {
				if first_url == url_key[1:] {
					proxy_url = url_key
					proxy_server = url_val
					break
				}
			}
		}

		if proxy_server == nil {
			return
		}

		newUrl, err := url.Parse(fmt.Sprintf("http://%s%s", proxy_server.Addr, strings.Replace(req.URL.Path, proxy_url, "", 1)))
		g.log.Infof("proxy_url is %s and proxy_server is %v\n", proxy_url, proxy_server)
		if err != nil {
			g.log.Errorf("server.go ServeHTTP method http proxy url.Parse err %v\n", err)
			return
		}
		//req里设置被代理服务器的前置url
		//设置web服务被代理的基础路径
		//暂时思路是在代理的服务器上使用request的header里设置一个信息来作为一个被代理的路径
		req.Header.Add(define.Gate_String_Web_Proxy, proxy_url)
		req.URL = newUrl
	}
	return &httputil.ReverseProxy{Director: director}
}

func (g *GateWebServer) Reg() {
	serverInfo := &define.ModelRegReqServerType{
		Addr:       g.gwsCfg.Addr,
		ServerName: define.ServerNameGate_WebServer,
		ServerType: 3,
	}

	for {
		data, err := json.Marshal(serverInfo)
		if err != nil {
			g.log.Errorf("server.go reg method json.Marshal(serverInfo) err %v\n", err)
			continue
		}

		//发送注册命令
		err = g.regConn.Send(define.CmdRegServer_Register, define.CmdSubRegServer_Register_Reg, data)
		if err != nil {
			g.log.Errorf("server.go reg method g.regConn.Send(define.CmdRegServer_Register %d, define.CmdSubRegServer_Register_Reg %d, data) err %v\n", define.CmdRegServer_Register, define.CmdSubRegServer_Register_Reg, err)
			continue
		}
		break
	}
}

//注册服务回调函数
func (g *GateWebServer) RegServiceCallBack(mainId, subId uint16, data []byte) bool {
	g.log.Infof("web gate server recve mainId %d subId %d\n", mainId, subId)
	if mainId == define.CmdRegServer_Register {
		switch subId {
		case define.CmdSubRegServer_Register_Reg_Inform:
			//下层服务上线通知
			res := &define.ModelRegReqServerType{}
			err := json.Unmarshal(data, res)
			if err != nil {
				g.log.Errorf("server.go RegServiceCallBack(%d, %d, data []byte) json.Unmarshal(data,res) err %v\n", mainId, subId, err)
				return true
			}

			//判断是否是逻辑服务
			//web网关接入web应用服务
			if res.ServerType == 2 {
				//配置相应的web微服务
				if len(res.ProxyUrl) < 0 {
					res.ProxyUrl = "/"
				}
				g.proxyUrls[res.ProxyUrl] = res
				g.log.Infof("web service proxy url %s host %s\n", res.ProxyUrl, res.Addr)

				//让注册服务器通知大厅服务器代理的信息
				lobProxyServer := &define.ModelRegReqLobbyProxyServer{}
				for mk, mv := range g.proxyUrls {
					lobProxyServer.ProxyUrls = append(lobProxyServer.ProxyUrls, mk)
					lobProxyServer.ServerNames = append(lobProxyServer.ServerNames, mv.ServerName)
				}
				lobProxyServer.InformServerName = define.ServerNameWeb_ManLobbyServer
				data, err := json.Marshal(lobProxyServer)
				if err != nil {
					g.log.Errorf("server.go RegServiceCallBack method json.Marshal(serverInfo) err %v\n", err)
					return true
				}

				//发送注册命令
				err = g.regConn.Send(define.CmdRegServer_Register, define.CmdSubRegServer_Register_Lobby_Proxy_Server, data)
				if err != nil {
					g.log.Errorf("server.go RegServiceCallBack method g.regConn.Send(define.CmdRegServer_Register %d, define.CmdSubRegServer_Register_Reg %d, data) err %v\n", define.CmdRegServer_Register, define.CmdSubRegServer_Register_Reg, err)
					return true
				}
			}
		case define.CmdSubRegServer_Register_Reg:
			res := &define.ModelRegResServerType{}
			err := json.Unmarshal(data, res)
			if err != nil {
				g.log.Errorf("server.go RegServiceCallBack(%d, %d, data []byte) json.Unmarshal(data,res) err %v\n", mainId, subId, err)
				return true
			}

			if res.Err == 0 {
				g.log.Infoln("service web gate server reg success")
			} else {
				g.log.Infoln("service web gate server reg fail")
			}
		}
	}
	return true
}

var gateWebServer *GateWebServer

func NewGateWebServer() *GateWebServer {
	if gateWebServer == nil {
		gateWebServer = new(GateWebServer)
	}
	return gateWebServer
}
