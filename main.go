package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var client mqtt.Client

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mqttBroker := os.Getenv("MQTT_BROKER")
	topic := os.Getenv("TOPIC")

	opts := mqtt.NewClientOptions().AddBroker("tcp://" + mqttBroker + ":1883").SetClientID("frigate_client")
	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	go listen(topic)

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "UP",
		})
	})

	r.Run()
}

func listen(topic string) {
	client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		img, _, err := image.Decode(bytes.NewReader(msg.Payload()))
		if err != nil {
			log.Fatalln(err)
		}
		currentTime := time.Now()
		f, err := os.Create(currentTime.Format("2006.01.02 15:04:05") + ".jpg")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		opt := jpeg.Options{
			Quality: 90,
		}
		jpeg.Encode(f, img, &opt)
		log.Info("done")
	})
}
