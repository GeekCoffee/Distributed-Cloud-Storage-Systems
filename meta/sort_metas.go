package meta

import "time"

//自定义fileMeta结构体的排序规则，需要实现sort.Interface中的所有方法，才能达到自定义排序的效果
//==>这里排序规则是使用时间来排序，具体就是使用Unix时间戳来排序


const BaseFormat = "2006-01-02 15:04:05"

type ByTimeSort []FileMeta

//sort.Interface中的方法，用于返回数组中的个数
func (fMetas ByTimeSort)Len() int{
	return len(fMetas)
}


//sort.Interface中的Swap方法，用于交换两个数组中的元素
func (fMetas ByTimeSort)Swap(i,j int){
	fMetas[i], fMetas[j] = fMetas[j], fMetas[i]
}


//核心：具体的自定义排序规则，使用Less(i,j)方法实现
//这里的具体规则是：使用Unix的时间戳来排序，时间戳大的话，说明是最新上传的文件
//固定语意是：当elem[i]<elem[j]是否成立，成立的话返回true，则说明elem[i]是比elem[j]小的，所以就产生了有序的数组
func (fMetas ByTimeSort)Less(i,j int) bool{
	iTime, _ := time.Parse(BaseFormat, fMetas[i].UploadTime)  //返回时间戳，以int64类型来表示
	jTime, _ := time.Parse(BaseFormat, fMetas[j].UploadTime)
	return iTime.UnixNano() > jTime.UnixNano()
}

