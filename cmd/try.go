package main

type per struct{
	name string
	age int32
}

func (p *per) print_name() error {
	print("name")
	println(p.name)
	print("Age")
	println(p.age)

	return nil
}