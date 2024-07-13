package main

import "github.com/ray-27/rayDB.git/assets"

func main() {
	for i:=0; i<16; i++{

		println(i)
		assets.RayDB_title(&i)
		println()
	}
}
