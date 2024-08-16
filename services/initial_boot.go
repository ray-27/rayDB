package services

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ray-27/rayDB.git/config"
)

func create_config_state() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter company name: ")
		input, err := reader.ReadString('\n')
		if err == nil {
			config.Company_name = input
			break
		}else{
			fmt.Println("Error reading the company name, re-enter :", err)
		}
	}

	for {
		fmt.Print("Enter email: ")
		input, err := reader.ReadString('\n')
		if err == nil {
			config.User_email = input
			break
		}else{
			fmt.Println("Error reading the email, re-enter", err)
		}
	}

	print("company : ",config.Company_name)
	print("email : ", config.User_email)
	
}

func Boot() {
	//define which path if the config_state file path
	// config_dir := "config/config_state.json"

	// //check if the file exists and if not then turn the `Initial_boot = true`
	// if _, err := os.Stat(config_dir); os.IsNotExist(err) {
	// 	//file does not exist, create a config_state.json file

	// 	err := os.MkdirAll(config_dir, os.ModePerm)
	// 	if err != nil {
	// 		fmt.Println("Error creating directory:", err)
	// 		return
	// 	}
	// 	fmt.Println("Directory created successfully")
	// } else if err != nil {
	// 	// An error other than "directory does not exist" occurred
	// 	fmt.Println("Error checking directory:", err)
	// } else {
	// 	// Directory exists
	// 	fmt.Println("Directory already exists")
	// }

	create_config_state()

}
