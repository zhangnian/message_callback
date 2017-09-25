package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func PanicOnError(err error, msg string) {
	if err != nil {
		log_msg := fmt.Sprintf("panic: [%s]%s", err.Error(), msg)
		log.Fatalln(log_msg)

		panic(err)
	}
}

func IsExist(fp string) bool {
	_, err := os.Stat(fp)
	return err == nil || os.IsExist(err)
}

func ReadFile(fp string) ([]byte, error) {
	return ioutil.ReadFile(fp)
}
