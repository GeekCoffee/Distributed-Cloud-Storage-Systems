package main

import (
	"fmt"
	"log"
	"os"
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


	//ts := fmt.Sprintf("%x", time.Now().Unix())
	//fmt.Println(ts[:8])
	//token := util.MD5([]byte("admin"+ts+"%!*(")) + ts[:8]
	//fmt.Println(token)
	//
	//tokenCreateTime, err := strconv.ParseUint(token[32:],16, 64)
	//if err != nil{
	//	log.Fatal(err)
	//}
	//
	//fmt.Println(tokenCreateTime)


	//nowTime := time.Now().Unix()  //1min = 60秒 , 60min = 3600s , 2h = 120min = 7200
	//if nowTime - tokenCreateTime >= 7200 { //nowTime时间点 >= tokenTime时间点 + 2h
	//	fmt.Println("token_time has been failed.")
	//}else{
	//	fmt.Println("token_time is valid.")
	//}

	//sumToken := 3600 * 24 * 30 * 2  //2个月60天
	//fmt.Println(sumToken)

	err := os.MkdirAll("./data/" + "uploadId",0666)
	if err != nil{
		log.Fatal(err)
	}
	_, err = os.Create("./data/" + "uploadId" + "/hello")
	if err != nil{
		log.Fatal(err)
	}


}