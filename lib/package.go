package lib

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

const (
	APK_REPO = "./repo"
	REPO_URL = "http://f-droid.org/repo/"
)

type Package struct {
	Version          string   `xml:"version"`
	Versioncode      int      `xml:"versioncode"`
	Apkname          string   `xml:"apkname"`
	Srcname          string   `xml:"srcname"`
	Hash             string   `xml:"hash"`
	Sig              string   `xml:"sig"`
	Size             int      `xml:"size"`
	Sdkver           int      `xml:"sdkver"`
	TargetSdkVersion int      `xml:"targetSdkVersion"`
	Added            string   `xml:"added"`
	Permissions      []string `xml:"permissions"`
}

func (p *Package) Download(applicationName string) error {
	currentPathDir := APK_REPO + "/" + applicationName + "/" + strconv.Itoa(p.Versioncode) + "/"
	currentApkFile := currentPathDir + p.Apkname
	if _, err := os.Stat(currentPathDir); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(currentPathDir, 0777); err != nil {
				return errors.New("Impossible to create the directory " + currentPathDir + ", due to error " + err.Error())
			}
			out, err := os.Create(currentApkFile)
			if err != nil {
				return errors.New("Canno't create directory in " + currentPathDir)
			}
			defer out.Close()
			resp, err := http.Get(REPO_URL + "/" + p.Apkname)
			if err != nil {
				return errors.New("Canno't get the APK " + p.Apkname + " from " + REPO_URL)
			}
			defer resp.Body.Close()
			_, err = io.Copy(out, resp.Body)
			if err != nil {
				return err
			}
			return nil
		} else {
			return errors.New("Repository " + currentPathDir + " already exists!")
		}
	} else {
		return errors.New("Error: stat is ok...")
	}
}

func (p *Package) String() (toReturn string) {

	toReturn += fmt.Sprintf("-> Version %s, added %s\n", p.Version, p.Added)
	toReturn += fmt.Sprintf("\t| Apk name: %s\n", p.Apkname)
	toReturn += fmt.Sprintf("\t| Size: %d bytes\n", p.Size)
	toReturn += fmt.Sprintf("\t| Sdk version: %d\n", p.TargetSdkVersion)

	return toReturn

}
