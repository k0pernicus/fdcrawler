package main

import (
	"encoding/xml"
	"fmt"
	"github.com/k0pernicus/fdcrawler/lib"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	REPO_URL  = "http://f-droid.org/repo/"
	INDEX_URL = REPO_URL + "index.xml"
)

// APK_REPO/ (package name)
// -> APK/
// -> -> /APK_VERSION/
// -> -> -> info.json
// -> -> -> APK_VERSION.apk

var APPS map[string]*lib.Application

func main() {

	APPS = make(map[string]*lib.Application)

	for {

		fmt.Println("Getting data...")

		http_get, err := http.Get(INDEX_URL)

		if err != nil {
			fmt.Println("Error fetching the webpage %s: %s", INDEX_URL, err)
			return
		}

		data, _ := ioutil.ReadAll(http_get.Body)

		http_get.Body.Close()

		var apps lib.FdroidRepo

		fmt.Println("Parsing data...")

		err = xml.Unmarshal([]byte(data), &apps)

		if err != nil {
			fmt.Println("error: %v", err)
			return
		}

		for _, app := range apps.Applications {
			if _, ok := APPS[app.Id]; !ok {
				app.PackagesList = make(map[int]bool)
				for _, p := range app.Packages {
					app.PackagesList[p.Versioncode] = true
					if err := p.Download(app.Id); err != nil {
						fmt.Println(err)
					}
				}
				APPS[app.Id] = &app
				fmt.Printf("[%d] ADDING NEW APP (%s)\n", len(APPS), app.Name)
				// app.String()
			} else {
				if APPS[app.Id].Lastupdated != app.Lastupdated {
					for _, _package := range app.Packages {
						if _, ok := APPS[app.Id].PackagesList[_package.Versioncode]; !ok {
							APPS[app.Id].AddPackage(&_package)
							fmt.Printf("ADDING NEW VERSION OF %s: %d\n", app.Id, _package.Versioncode)
						}
					}
					APPS[app.Id].Lastupdated = app.Lastupdated
				}
			}
		}

		fmt.Println("Sleeping 2 minutes...")

		time.Sleep(2 * time.Minute)

	}
}
