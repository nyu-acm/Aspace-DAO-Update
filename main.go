package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/nyudlts/go-aspace/lib"
	"os"
)

var aspace = lib.Client

var t = newTrue()

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
	//get the repository id from the command line
	repPntr := flag.Int("repository", 0, "the repository")
	flag.Parse()

	logFile, err := os.Create(fmt.Sprintf("logfile-%d.txt", *repPntr))
	checkPanic(err)
	defer logFile.Close()
	logWriter := bufio.NewWriter(logFile)

	daos, err := aspace.GetDigitalObjectsByRepositoryId(*repPntr)
	checkPanic(err)

	for _, daoId := range daos {

		dao, err := aspace.GetDigitalObject(*repPntr, daoId)
		checkPanic(err)

		//check if there are file versions
		if len(dao.FileVersions) > 0 {
			var update = false
			//check fileversion subrecords for ones where publish value is nil
			for _, fv := range dao.FileVersions {
				if fv.Publish == nil {
					//If the parent DAO's publish field is set to True, update
					if dao.Publish == true {
						update = true
						fv.Publish = newTrue()
					}
				}
			}

			if update == true {
				fmt.Println("Updating", dao.URI)
				updateResponse, err := aspace.PostDigitalObject(*repPntr, daoId, dao)
				if err != nil {
					logWriter.WriteString(dao.URI + "\t" + err.Error())
					logWriter.Flush()
				}
				fmt.Println(updateResponse)
			}
		}

	}

	logWriter.Flush()

}
