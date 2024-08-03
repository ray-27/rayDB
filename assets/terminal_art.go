package assets

import (
	"embed"
	"fmt"
	"io"
	"os"
)

//go:embed logo_art/*.txt
var logoArtFS embed.FS

func RayDB_logo(choice int) {
	/*
		choice: int
			takes the choice of the .txt file to be used as the logo
			the logo are saved in the ./assets/logo_art
	*/

	fileName := fmt.Sprintf("logo_art/%d.txt", choice)

	file, err := logoArtFS.Open(fileName)
	if err != nil {
		fmt.Println("Error showing the logo:", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(os.Stdout, file)
	if err != nil {
		fmt.Println("Error reading the logo file:", err)
	}
}
