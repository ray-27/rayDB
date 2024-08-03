package services

type Arr_int []int

func (arr Arr_int) Print(){
	l := len(arr)
	for i:=0;i < l;i++{
		print(arr[i]," ")
	}
}