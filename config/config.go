package config

// creating some universal variable names 

var Initial_boot bool = false //stores to check if the initial boot is done, turns true if the files exist and 
var Company_name string = ""
var User_email string
var App_PORT int = 8080 //the port on which the database server is going to run, default to 8080 but can change in the initial_boot function
var DebugMode bool = true 
var Config_file_folder string