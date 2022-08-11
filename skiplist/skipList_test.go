package skiplist

import (
	"fmt"
	"testing"
)

func Test_Operate(t *testing.T) {

	var cmp *CmpInstanceInt
	//sl, _ := New(cmp, WithAllowTheSameKey(true), WithLevelCacheSize(10, 2, 3, 4, 7, 4, 5))
	sl, _ := New(cmp, WithAllowTheSameKey(false))
	sl.addNode(CmpInstanceInt(1), 1)
	sl.addNode(CmpInstanceInt(3), 3)
	sl.addNode(CmpInstanceInt(2), 2)
	sl.updateBatchByKey(CmpInstanceInt(2), 22)
	//fmt.Println(sl.addNode(CmpInstanceInt(1), 11))
	//fmt.Println(sl.updateByKey(CmpInstanceInt(1), 111))
	//sl.updateBatchByKey(CmpInstanceInt(1), -9)
	//sl.delByKey(CmpInstanceInt(2))
	fmt.Println(sl.delByKey(CmpInstanceInt(2)))
	fmt.Println(sl.DeleteByRank(3))
	sl.addNode(CmpInstanceInt(4), 4)
	sl.addNode(CmpInstanceInt(5), 5)
	sl.addNode(CmpInstanceInt(11), 11)
	sl.addNode(CmpInstanceInt(6), 6)
	sl.delByKey(CmpInstanceInt(6))
	fmt.Println("------------------")
	fmt.Println(sl.searchByRank(1).data)
	fmt.Println(sl.searchByRank(4).data)
	fmt.Println(sl.searchByRank(5).data)
	fmt.Println(sl.searchByRank(6))
	fmt.Println("------------------")
	list := sl.searchByRankRange(1, 3)
	for k := range list {
		fmt.Println(list[k].data)
	}
	fmt.Println("------------------")
	fmt.Println(sl.searchNodeAndRankByKey(CmpInstanceInt(4)))
	fmt.Println("------------------")
	it := NewIterator(sl)
	it.PrintGraph()
}

func Benchmark_AddNode(b *testing.B) {
	var cmp *CmpInstanceInt
	sl, _ := New(cmp)
	for i := 0; i < b.N; i++ {
		sl.addNode(CmpInstanceInt(i), i)
	}
}
func Benchmark_delNodeWithSameKey(b *testing.B) {
	var cmp *CmpInstanceInt
	sl, _ := New(cmp)
	for i := 0; i < b.N; i++ {
		sl.addNode(CmpInstanceInt(i), i)
		sl.delByKey(CmpInstanceInt(i))
	}
}
func Benchmark_delNodeWithOutSameKey(b *testing.B) {
	var cmp *CmpInstanceInt
	sl, _ := New(cmp, WithAllowTheSameKey(false))
	for i := 0; i < b.N; i++ {
		sl.addNode(CmpInstanceInt(i), i)
		sl.delByKey(CmpInstanceInt(i))
	}
}
func Benchmark_SearchRandKey(b *testing.B) {
	var cmp *CmpInstanceInt
	sl, _ := New(cmp, WithAllowTheSameKey(false))
	for i := 0; i < b.N; i++ {
		sl.addNode(CmpInstanceInt(i), i)
		sl.searchRandOneByKey(CmpInstanceInt(i))
	}
}

func Benchmark_SearchRank(b *testing.B) {
	var cmp *CmpInstanceInt
	sl, _ := New(cmp, WithAllowTheSameKey(false))
	for i := 0; i < b.N; i++ {
		sl.addNode(CmpInstanceInt(i), i)
		sl.searchByRank(i)
	}
}

func Benchmark_SearchRankRange(b *testing.B) {
	var cmp *CmpInstanceInt
	sl, _ := New(cmp, WithAllowTheSameKey(false))
	for i := 0; i < b.N; i++ {
		sl.addNode(CmpInstanceInt(i), i)
		if i > 10 {
			sl.searchByRankRange(i-i%10, i)
		}
	}
}
