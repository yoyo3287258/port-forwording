package web

import (
	"embed"
	"github.com/gin-gonic/gin"
	"io"
	"io/fs"
	"port-forwording/web/routers"
)

//go:embed dist
var Dist embed.FS

func Run() {
	r := gin.Default()
	index(r)
	routers.InitRouters(r)
	err := r.Run(":8765")
	if err != nil {
		return
	}
}

func index(r *gin.Engine) {
	r.NoRoute(func(c *gin.Context) {
		dist, _ := fs.Sub(Dist, "dist")
		indexFile, err := dist.Open("index.html")
		if err != nil {
			panic(err.Error())
		}
		defer func() {
			_ = indexFile.Close()
		}()
		content, _ := io.ReadAll(indexFile)

		c.Header("Content-Type", "text/html")
		c.Status(200)
		_, _ = c.Writer.WriteString(string(content))
		c.Writer.Flush()
		c.Writer.WriteHeaderNow()
	})
}
