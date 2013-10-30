package quiet

import (
	"io/ioutil"
	"log"
)

func init() {
	log.SetOutput(ioutil.Discard)
}
