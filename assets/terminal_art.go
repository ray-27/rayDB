package assets

import (
	"fmt"
	"io"
	"os"
)

func RayDB_title(choice *int) {

	c := 0
	defaultChoice := 7
	if choice == nil{
		c = defaultChoice
	}else {
		c = &choice
	}
	fileName := fmt.Sprintf("assets/logo_art/%d.txt",c)

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error showing the logo")
	}
	defer file.Close()

	_, err = io.Copy(os.Stdout, file)
	if err != nil {
		fmt.Println("Error reading the logo file")
	}
}
