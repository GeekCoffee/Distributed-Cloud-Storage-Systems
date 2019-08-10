package main

import (
	"fmt"
)

//用于测试标准库的包和函数的

func ff(f func(int,int),a ,b int)  {
	f(a,b)
}

func sum(a,b int) {
	fmt.Println(a+b)
}

func main() {
	//ff(sum,1,1)
	//file, err := os.Open("./tmp/web展示.zip")
	//if err != nil{
	//	fmt.Printf("error: %s", err.Error())
	//}
	//
	//fileHash := util.FileSha1(file)
	//fmt.Println(fileHash)

	//arr := []int{0,1,2,3,4,5,6,7,8,9}
	//fmt.Println(arr[0:5])  //从下标0开始，返回包括0在内的5个元素

	//fmt.Println(util.Sha1([]byte("hello"))) //40个字符
	//fmt.Println(util.MD5([]byte("hello")))  //32个字符
	//fmt.Println(util.Sha256([]byte("hello")))  //64个字符
	//
	//count := 0
	//for i,c := range "5d41402abc4b2a76b9719d911017c592"{
	//	i = i
	//	c = c
	//	count++
	//}
	//fmt.Println(count)

}