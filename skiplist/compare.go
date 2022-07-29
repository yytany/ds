package skiplist

//比较接口原型
type CompareAble interface {
	Compare(a, b interface{}) int // -1 a<b  0 a==b  1 a>b
}

// 实现比较接口实现示例
type CmpInstanceInt int

func (*CmpInstanceInt) Compare(a, b interface{}) int {
	if a.(CmpInstanceInt) < b.(CmpInstanceInt) {
		return -1
	} else if a.(CmpInstanceInt) > b.(CmpInstanceInt) {
		return 1
	}
	return 0
}

// 实现比较接口实现示例
type CmpInstanceStruct struct {
	price      float64 //价格低的排前面  顺位1
	createTime int64   //时间大的排前面(新的在前面)	 顺位2
}

func (*CmpInstanceStruct) Compare(a, b interface{}) int {
	if a.(CmpInstanceStruct).price < b.(CmpInstanceStruct).price {
		return -1
	} else if a.(CmpInstanceStruct).price > b.(CmpInstanceStruct).price {
		return 1
	} else {
		if a.(CmpInstanceStruct).createTime < b.(CmpInstanceStruct).createTime {
			return -1
		} else if a.(CmpInstanceStruct).createTime > b.(CmpInstanceStruct).createTime {
			return 1
		}
	}
	return 0
}
