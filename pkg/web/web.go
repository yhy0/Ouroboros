package web

import (
    "embed"
    "encoding/json"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "github.com/yhy0/InfiniteSnake/conf"
    "github.com/yhy0/InfiniteSnake/pkg/process"
    "github.com/yhy0/logging"
    "html/template"
    "net/http"
    "time"
)

/**
   @author yhy
   @since 2023/10/15
   @desc //TODO
**/

//go:embed templates
var templates embed.FS

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true // 允许跨域请求
    },
}

// handleWebSocket 使用 websocket 推送数据，同步页面更改
func handleWebSocket(c *gin.Context) {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        logging.Logger.Errorln("Failed to upgrade request to WebSocket:", err)
        return
    }
    defer conn.Close()
    
    // 数据更改
    for {
        messageType, p, err := conn.ReadMessage()
        if err != nil {
            return
        }
        if messageType == websocket.TextMessage {
            var message map[string]interface{}
            if err := json.Unmarshal(p, &message); err == nil {
                processId := message["pid"].(string)
                checkboxValue := message["checkboxValue"].(bool)
                if checkboxValue {
                    process.MonitorPid[processId] = process.Processes[processId]
                    delete(process.Processes, processId)
                } else {
                    // 关闭监控
                    process.Processes[processId] = process.MonitorPid[processId]
                    delete(process.MonitorPid, processId)
                }
            }
        }
        time.Sleep(1 * time.Second)
    }
}

func Init() {
    logging.Logger.Infoln("Start SCopilot web service at :" + conf.WebPort)
    gin.SetMode("release")
    router := gin.Default()
    
    // 设置模板资源
    router.SetHTMLTemplate(template.Must(template.New("").ParseFS(templates, "templates/*")))
    
    router.GET("/ws", handleWebSocket)
    
    // basic 认证
    authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
        conf.WebUser: conf.WebPass,
    }))
    
    authorized.GET("/", func(c *gin.Context) {
        c.Redirect(302, "/index")
    })
    
    authorized.GET("/index", func(c *gin.Context) {
        process.GetUserProcesses()
        c.HTML(http.StatusOK, "index.html", gin.H{
            "webPort":        conf.WebPort,
            "MonitorProcess": process.MonitorPid,
            "process":        process.Processes,
            "year":           time.Now().Year(),
        })
    })
    
    router.Run(":" + conf.WebPort)
}
