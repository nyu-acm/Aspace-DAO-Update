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

	//create a logfile and writer object
	logFile, err := os.Create(fmt.Sprintf("logfile-%d.txt", *repPntr))
	checkPanic(err)
	defer logFile.Close()
	logWriter := bufio.NewWriter(logFile)

	//get a list of doas in the repository from the Aspace API
	daos, err := aspace.GetDigitalObjectsByRepositoryId(*repPntr)
	checkPanic(err)

	//iterate through the daos
	for _, daoId := range daos {

		//request a dao from the aspace API
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

			//check to see if the dao should be updated
			if update == true {
				fmt.Println("Updating", dao.URI)
				// post the updated dao to aspace
				updateResponse, err := aspace.PostDigitalObject(*repPntr, daoId, dao)
				//log any errors
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
