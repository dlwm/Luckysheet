package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/text/encoding/charmap"
	"io"
	"io/ioutil"
	"log"
	"math/rand/v2"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // 允许所有来源，实际应用中应该更严格
		},
	}
	clients   = make(map[*websocket.Conn]struct{})
	mx        = sync.RWMutex{}
	broadcast = make(chan struct {
		d []byte
		c *websocket.Conn
	}, 1024)
	content = DefContent
)

func init() {
	f, _ := os.Open("./x.json")
	buf := bytes.Buffer{}
	io.Copy(&buf, f)
	content = buf.String()
}

func handleUpdate(ctx *gin.Context) {
	uid := strconv.Itoa(rand.Int())
	name := "user_" + uid
	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	mx.Lock()

	clients[ws] = struct{}{}

	mx.Unlock()

	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			delete(clients, ws)
			break
		}
		switch msgType {
		case websocket.TextMessage:
			data, err := ungzip(msg)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(data))
			var rsp struct {
				Type     int    `json:"type"`
				Id       string `json:"id,omitempty"`
				UserName string `json:"username,omitempty"`
				Data     string `json:"data"`
			}

			var msgSend struct {
				T string `json:"t"`
			}
			_ = sonic.Unmarshal(data, &msgSend)
			switch msgSend.T {
			case "v", "rv", "rv_end", "cg", "all", "fc", "drc", "arc", "f", "fsc", "fsr", "sha", "shc", "shd", "shr", "shre", "sh", "c", "na":
				rsp.Type = 2
			case "mv":
				rsp.Type = 3
				rsp.Id = uid
				rsp.UserName = name
			case "": //离线情况下把更新指令打包批量下发给客户端
				rsp.Type = 4
			default:
				rsp.Type = 1
			}

			f, ok := updMap[msgSend.T]
			if ok {
				f(data, &content)
			}
			//fmt.Println(content)

			rsp.Data = string(data)
			byts, _ := sonic.Marshal(rsp)

			broadcast <- struct {
				d []byte
				c *websocket.Conn
			}{d: byts, c: ws}
		}
	}
}
func handleLoad(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, content) // todo lock
}

func main() {
	engine := gin.Default()
	grp := engine.Group("/", corsMiddleware)
	grp.GET("/collaborate", handleUpdate)
	grp.POST("/collaborate", handleLoad)

	go func() {
		for {
			select {
			case msg := <-broadcast:
				mx.RLock()
				for conn := range clients {
					if msg.c == conn {
						continue
					}
					if err := conn.WriteMessage(websocket.TextMessage, msg.d); err != nil {
						fmt.Println(err)
					}
				}
				mx.RUnlock()
			}
		}
	}()

	engine.Run(":2234")
}

func corsMiddleware(ctx *gin.Context) {
	// 设置跨域响应头
	ctx.Header("Access-Control-Allow-Origin", "*") // 允许所有域访问（生产环境请限制域名）
	ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	ctx.Header("Access-Control-Allow-Headers", "*")

	// 预检请求处理
	if ctx.Request.Method == http.MethodOptions {
		ctx.JSON(http.StatusOK, nil)
		return
	}
}

func ungzip(gzipmsg []byte) (reqmsg []byte, err error) {
	if len(gzipmsg) == 0 {
		return
	}
	if string(gzipmsg) == "rub" {
		reqmsg = gzipmsg
		return
	}
	e := charmap.ISO8859_1.NewEncoder()
	encodeMsg, err := e.Bytes(gzipmsg)
	if err != nil {
		return
	}
	b := bytes.NewReader(encodeMsg)
	r, err := gzip.NewReader(b)
	if err != nil {
		return
	}
	defer r.Close()
	reqmsg, err = ioutil.ReadAll(r)
	if err != nil {
		return
	}
	reqstr, err := url.QueryUnescape(string(reqmsg))
	if err != nil {
		return
	}
	reqmsg = []byte(reqstr)
	return
}
