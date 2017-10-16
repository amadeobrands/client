package main

import (
	"fmt"
	"os"

	"github.com/loomnetwork/client/auth"
	"github.com/loomnetwork/client/client"
	"github.com/loomnetwork/client/config"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	app   = kingpin.New("loom", "Loom network deployment tool.")
	debug = app.Flag("debug", "Enable debug mode.").Bool()

	login   = app.Command("login", "Logins you into loom network. It also retrieves api key.")
	network = login.Arg("network", "Either Linkedin/Github. Leave blank for a prompt").String()

	upload  = app.Command("upload", "Upload an application package.")
	slug    = upload.Arg("appName", "The shortname for your application. Will create the domain <appname>.loomapps.io .").Required().String()
	zipfile = upload.Arg("zipfile", "File to upload.").Required().String()

	setapikey = app.Command("setapikey", "Set api key.")
	apikey    = setapikey.Arg("key", "Key.").Required().String()
)

var DEFAULT_HOST = "https://platform.loomx.io"

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case login.FullCommand():
		auth.Login(*network)
	case upload.FullCommand():
		if *zipfile == "" {
			fmt.Println("Please supply a zipfile")
			return
		}
		c := config.ReadConfig()
		if c.Apikey == "" || len(c.Apikey) < 3 {
			fmt.Println("Missing api key or api key is invalid. Please set it with the setapikey command.")
			return
		}
		hostName := c.HostName
		if hostName == "" {
			hostName = DEFAULT_HOST
		}
		client.UploadApp(hostName, c.Apikey, *zipfile, *slug)
	case setapikey.FullCommand():
		if *apikey == "" || len(*apikey) < 3 {
			fmt.Println("Missing api key or api key is invalid")
			return
		}
		config.WriteConfig(*apikey)
		fmt.Println("Api successfully set!")
	}
}
