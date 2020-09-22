package main

import (
	"fmt"
	"github.com/nyudlts/go-aspace/lib"
)

var aspace = lib.Client

func checkPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func newTrue() *bool {
	b := true
	return &b
}

func main() {

	repositoryId := 2
	daos, err := aspace.GetDigitalObjectsByRepositoryId(repositoryId)
	checkPanic(err)

	for daoId := range daos {
		fmt.Println(repositoryId, daoId)

		dao, err := aspace.GetDigitalObject(repositoryId, daoId)
		checkPanic(err)

		if len(dao.FileVersions) > 0 {
			for _, fv := range dao.FileVersions {
				if fv.Publish == nil {
					if dao.Publish == true {
						fv.Publish = newTrue()
					}
				}
			}

		}

		//check that fileversions are updated

		// update dao in aspace

	}

}
