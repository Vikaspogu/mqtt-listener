package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var (
	client mqtt.Client
	path   string
)

func main() {
	godotenv.Load()
	mqttBroker := os.Getenv("MQTT_BROKER")
	topic := os.Getenv("TOPIC")
	path = "/data"

	opts := mqtt.NewClientOptions().AddBroker("tcp://" + mqttBroker + ":1883").SetClientID("frigate_client")
	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	go listen(topic)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health"},
	}))
	router.Use(gin.Recovery())
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "UP",
		})
	})

	router.Run()
}

func listen(topic string) {
	client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		img, _, err := image.Decode(bytes.NewReader(msg.Payload()))
		if err != nil {
			log.Fatalln(err)
		}
		location, err := time.LoadLocation("America/Chicago")
		if err != nil {
			log.Error(err)
		}
		currentTime := time.Now().In(location)
		currentDate := currentTime.Format("2006-01-02")
		currDir := fmt.Sprintf("%s%s%s", path, "/", currentDate)
		timeFormatted := strings.ReplaceAll(currentTime.Format("15:04:05"), ":", "-")
		imageName := fmt.Sprintf("%s%s", timeFormatted, ".jpg")

		if _, err := os.Stat(currDir); os.IsNotExist(err) {
			err := os.Mkdir(currDir, 0755)
			if err != nil {
				log.Errorf("Cannot create folder %s", currentDate)
			}
		}
		fileName := fmt.Sprintf("%s%s%s", currDir, "/", imageName)
		f, err := os.Create(fileName)
		if err != nil {
			log.Error(err)
			return
		}
		defer f.Close()
		opt := jpeg.Options{
			Quality: 90,
		}
		jpeg.Encode(f, img, &opt)
		log.Infof("Uploading...%s", imageName)
	})
}
