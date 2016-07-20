package main

import (
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"github.com/k0pernicus/fdcrawler/lib"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	APK_REPO  = "./repo"
	REPO_URL  = "http://f-droid.org/repo/"
	INDEX_URL = REPO_URL + "index.xml"
)

// APK_REPO/ (package name)
// -> APK/
// -> -> /APK_VERSION/
// -> -> -> info.json
// -> -> -> APK_VERSION.apk

var APPS map[string]*lib.Application

func Download(applicationName string, applicationAPK string, versionCode string) error {
	currentPathDir := APK_REPO + "/" + applicationName + "/" + versionCode + "/"
	currentApkPath := currentPathDir + applicationAPK
	fmt.Println("Path: " + currentPathDir + " / APK: " + currentApkPath)
	if _, err := os.Stat(currentPathDir); err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(currentPathDir, 0777); err != nil {
				return errors.New("Impossible to create the directory " + currentPathDir + ", due to error " + err.Error())
			}
		}
	}
	out, erro := os.Create(currentApkPath)
	if erro != nil {
		return errors.New("Canno't create directory in " + currentPathDir + ", due to error " + erro.Error())
	}
	defer out.Close()
	var resp *http.Response
	var err error
	for {
		resp, err = http.Get(REPO_URL + "/" + applicationAPK)
		if err == nil {
			break
		}
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func main() {

	APPS = make(map[string]*lib.Application)

	// timeSleep := flag.Int("sleep", 2, "the time sleep between 2 crawls (in hours)")

	flag.Parse()

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
					fmt.Println(p.String())
					app.PackagesList[p.Versioncode] = true
					// TODO: GoRoutine
					go Download(app.Id, p.Apkname, strconv.Itoa(p.Versioncode))
				}
				APPS[app.Id] = &app
				fmt.Printf("[%d] ADDING NEW APP (%s)\n", len(APPS), app.Name)
				// app.String()
			} else {
				if APPS[app.Id].Lastupdated != app.Lastupdated {
					for _, _package := range app.Packages {
						if _, ok := APPS[app.Id].PackagesList[_package.Versioncode]; !ok {
							if err := APPS[app.Id].AddPackage(&_package); err == nil {
								// go func() {
								// 	Download(app.Id, _package.Apkname, strconv.Itoa(_package.Versioncode))
								// }()
								fmt.Printf("ADDING NEW VERSION OF %s: %d\n", app.Id, _package.Versioncode)
							} else {
								fmt.Println(err)
							}
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
