package main

import (
	"github.com/ray-27/rayDB.git/services"
)

type Db_data struct {
	name string
	num  int
}

func main() {

	services.Boot()

}
