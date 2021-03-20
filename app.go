package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Kta-M/web-mtg_notifier/mqtt"
	"github.com/Kta-M/web-mtg_notifier/webmtg_status"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "WebMTG Notifier",
		Usage:   "Notify the start and end of WebMTG",
		Version: "v0.1.0",
		Action:  run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "user-name",
				Aliases:  []string{"u"},
				Usage:    "User Name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "endpoint",
				Aliases:  []string{"e"},
				Usage:    "Endpoint of AWS IoT Core",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "root-ca",
				Aliases:  []string{"r"},
				Usage:    "Root CA file path",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "cert",
				Aliases:  []string{"c"},
				Usage:    "Cert file path",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "key",
				Aliases:  []string{"k"},
				Usage:    "Key file path",
				Required: true,
			},
			&cli.IntFlag{
				Name:    "interval",
				Aliases: []string{"i"},
				Usage:   "Check interval (seconds)",
				Value:   60,
			},
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	// フラグの値を取得
	userName := c.String("user-name")
	endpoint := c.String("endpoint")
	rootCAPath := c.String("root-ca")
	certPath := c.String("cert")
	keyPath := c.String("key")
	interval := c.Int("interval")

	// ブローカーに接続
	client, err := mqtt.Connect(userName, endpoint, rootCAPath, certPath, keyPath)
	if err != nil {
		log.Fatal(err)
	}
	defer mqtt.Disonnect(client, 250)

	var status bool = webmtg_status.GetStatus()
	topic := fmt.Sprintf("topic/%v", userName)

	// 初回送信
	sendStatus(client, topic, status)

	// メインループ
	for {
		oldStatus := status
		status = webmtg_status.GetStatus()
		// log.Printf("cur: %v, old: %v\n", status, oldStatus)

		if status != oldStatus {
			sendStatus(client, topic, status)
		}

		time.Sleep(time.Duration(interval) * time.Second)
	}
}

// メッセージ送信
func sendStatus(client mqtt.Client, topic string, status bool) {
	message := fmt.Sprintf(`{"status": "%v"}`, status)
	// log.Printf("publishing %s : %s...\n", topic, message)
	if err := mqtt.Publish(client, topic, 1, false, message); err != nil {
		log.Println(err)
	}
}
