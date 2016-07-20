package lib

import "fmt"

type FdroidRepo struct {
	Applications []Application `xml:"application"`
}

func (f *FdroidRepo) String() (toReturn string) {
	for _, app := range f.Applications {
		toReturn += app.String()
		toReturn += fmt.Sprintf("#######--#######\n")
	}
	return toReturn
}
