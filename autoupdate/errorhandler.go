package autoupdate

import "log"

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
