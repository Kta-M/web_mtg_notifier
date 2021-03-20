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

	status := false
	topic := fmt.Sprintf("topic/%v", userName)

	// メインループ
	for {
		oldStatus := status
		status = webmtg_status.GetStatus()
		// log.Printf("cur: %v, old: %v\n", status, oldStatus)

		if status != oldStatus {
			// log.Printf("publishing %s...\n", topic)
			message := fmt.Sprintf(`{"status": "%v"}`, status)
			if err := mqtt.Publish(client, topic, 1, false, message); err != nil {
				log.Println(err)
			}
		}

		time.Sleep(time.Duration(interval) * time.Second)
	}
}
