package lib

import (
	"fmt"
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

func (p *Package) String() (toReturn string) {

	toReturn += fmt.Sprintf("-> Version %s, added %s\n", p.Version, p.Added)
	toReturn += fmt.Sprintf("\t| Apk name: %s\n", p.Apkname)
	toReturn += fmt.Sprintf("\t| Size: %d bytes\n", p.Size)
	toReturn += fmt.Sprintf("\t| Sdk version: %d\n", p.TargetSdkVersion)

	return toReturn

}
