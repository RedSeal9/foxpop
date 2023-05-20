package main

import (
	"context"
	"embed"
	"fmt"
	foxpop "foxpop/prefs"
	"io"
	"net"
	"os"
	"path"

	"github.com/tawesoft/golib/v2/dialog"
	"gopkg.in/ini.v1"
)

const PropertyServer = "foxprops.loc"

//go:embed help.txt
var content embed.FS

func main() {
	server := PropertyServer

	rawArgs := os.Args[1:]
	if len(rawArgs) != 0 {
		if rawArgs[0] == "h" { // help
			hf, _ := content.Open("help.txt")
			hc, _ := io.ReadAll(hf)
			fmt.Printf(string(hc))
			os.Exit(0)
		}
		if rawArgs[0] == "c" { // create new propsfile
			ingestFile := rawArgs[1]
			if _, err := os.Stat(ingestFile); err != nil {
				dialog.Error("Propsfile does not exist")
				os.Exit(1)
			}

			exportFile := rawArgs[2]
			err := foxpop.ParseDataFile(ingestFile, exportFile)
			if err != nil {
				fmt.Printf("err.Error(): %v\n", err.Error())
			}
			os.Exit(0)
		}
		if rawArgs[0] == "l" {
			server = rawArgs[1]
		}
	}

	profPath := findProfilePath()
	if profPath == "" {
		dialog.Error("Firefox install was not found.")
		os.Exit(1)
	}
	jsFileLocation := path.Join(profPath, "user.js")
	fmt.Printf("jsFileLocation: %v\n", jsFileLocation)

	if !canFindPrefServer(server) {
		dialog.Error("Server was not found.")
		os.Exit(1)
	}

	d, err := foxpop.ReturnProperties(server)
	if err != nil {
		dialog.Error("Error encountered while fetching properties!")
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	ujs := foxpop.StringifyUserJS(d)
	if res, _ := dialog.Ask("Are you okay with %v change(s)?", len(d.Entries)); res {
		dialog.Info("Making changes to 'user.js'")
		err := os.WriteFile(jsFileLocation, []byte(ujs), 0666)
		if err != nil {
			dialog.Error("Changes failed")
		}
	}
}

func findProfilePath() string {
	ad := os.Getenv("AppData")
	ad = path.Join(ad, "Mozilla", "Firefox")
	pil := path.Join(ad, "profiles.ini")
	if _, err := os.Stat(ad); err == nil {
		pfs, _ := ini.Load(pil)
		pl := pfs.Section("Profile0").Key("Path").String()
		return path.Join(ad, pl)
	} else {
		return ""
	}
}

func canFindPrefServer(server string) bool {
	_, err := net.DefaultResolver.LookupIP(context.Background(), "ip4", server)
	return err == nil
}
