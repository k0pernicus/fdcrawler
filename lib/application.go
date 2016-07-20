package lib

import (
	"errors"
	"fmt"
)

type Application struct {
	Id            string    `xml:"id"`
	Added         string    `xml:"added"`
	Lastupdated   string    `xml:"lastupdated"`
	Name          string    `xml:"name"`
	Summary       string    `xml:"summary"`
	Icon          string    `xml:"icon"`
	License       string    `xml:"license"`
	Categories    []string  `xml:"categories"`
	Category      string    `xml:"category"`
	Web           string    `xml:"web"`
	Source        string    `xml:"source"`
	Tracker       string    `xml:"tracker"`
	Marketversion string    `xml:"marketversion"`
	Marketvercode int       `xml:"marketvercode"`
	Requirements  []string  `xml:"requirements"`
	Packages      []Package `xml:"package"`
	PackagesList  map[int]bool
}

func (a *Application) AddPackage(p *Package) error {
	if _, ok := a.PackagesList[p.Versioncode]; ok {
		return errors.New("Package already exists")
	}
	a.Packages = append(a.Packages, *p)
	a.PackagesList[p.Versioncode] = true
	if err := p.Download(a.Id); err != nil {
		fmt.Printf("Canno't download the current application - due to error %s\n", err)
		return err
	}
	return nil
}

func (a *Application) Compare(oa *Application) bool {
	if a.Name == oa.Name && a.Id == oa.Id && a.Lastupdated == oa.Lastupdated {
		return true
	}
	return true
}

func (a *Application) String() (toReturn string) {

	toReturn += fmt.Sprintf("* Application %s\n", a.Name)
	toReturn += fmt.Sprintf("-> Package: %s\n", a.Id)
	toReturn += fmt.Sprintf("-> Description: %s\n", a.Summary)
	toReturn += fmt.Sprintf("-> Category: %s\n", a.Category)
	toReturn += fmt.Sprintf("-> Source code available in %s\n", a.Source)
	for _, p := range a.Packages {
		toReturn += p.String()
	}
	toReturn += fmt.Sprintf("-> Market version: %s\n", a.Marketversion)
	toReturn += fmt.Sprintf("-> License %s\n", a.License)

	return toReturn

}
