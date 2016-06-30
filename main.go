package main

import (
	"github.com/davecheney/i2c"
	"github.com/gin-gonic/gin"
	"github.com/quinte17/bme280"
	"log"
	"time"
	"html/template"
)

type MyData struct {
	bme280.Envdata
	MmHg      float64
	Timestamp time.Time
}

var data MyData

func main() {

	dev, err := i2c.New(0x76, 0)
	if err != nil {
		log.Print(err)
	}

	bme, err := bme280.NewI2CDriver(dev)
	if err != nil {
		log.Print(err)
	}

	ticker := time.NewTicker(time.Second * 3)

	go func() {
		for t := range ticker.C {
			env, err := bme.Readenv()
			if err != nil {
				log.Print("[ERROR] %v", err)
				continue
			}
			data.Envdata = env
			data.MmHg = env.Press / 1.333224
			data.Timestamp = t
			log.Print(data)
		}
	}()

	r := gin.Default()
	t := template.New("index.html")
	t.Parse(string(MustAsset("index.html")))
	r.SetHTMLTemplate(t)

	r.GET("/bme280", func(c *gin.Context) {
		c.Header("Refresh", "5")
		c.HTML(200, "index.html", data)
	})

	r.GET("/bme280/json", func(c *gin.Context) {
		c.JSON(200, data)
	})

	log.Print("Listen at :8280")
	r.Run(":8280")
}
