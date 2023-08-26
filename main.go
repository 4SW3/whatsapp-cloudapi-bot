package main

import (
	"flag"
	"log"
	"os"

	"github.com/4d3v/gowtbot/pkg/config"
	"github.com/4d3v/gowtbot/pkg/db"
	"github.com/4d3v/gowtbot/pkg/driver"
	"github.com/4d3v/gowtbot/pkg/wt"

	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
)

var (
	addr     = flag.String("addr", ":8080", "TCP address to listen to")
	compress = flag.Bool("compress", true, "Whether to enable transparent response compression")
	inDev    = flag.Bool("inDev", true, "Whether mode is in development or production")
	app      = config.AppConfig{}
)

func main() {
	flag.Parse()

	run()

	h := wt.WebhookHandler
	if *compress {
		log.Println("compressing")
		h = fasthttp.CompressHandler(h)
	}

	if err := fasthttp.ListenAndServe(*addr, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %v", err)
	}
}

func run() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app.AppSecret = []byte(os.Getenv("APP_SECRET"))
	app.AccessTkn = []byte(os.Getenv("ACCESS_TOKEN"))
	app.HubToken = []byte(os.Getenv("HUB_TOKEN"))
	app.WTPhoneID = []byte(os.Getenv("WHATSAPP_PHONE_ID"))
	app.InDevelopment = *inDev

	firebaseApp, blockedPhones, err := driver.FirebaseInit()
	if err != nil {
		log.Fatal(err)
	}

	wt.NewHandlers(wt.NewWTRepo(
		&app,
		db.NewFirebaseRepo(firebaseApp),
		make(map[string]bool),
		blockedPhones,
	))
}
