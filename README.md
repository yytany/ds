# ds
data structure design and achieve
数据结构设计和实现 

持续更新一些有意思的数据结构  [队列,跳表]...

## 使用
包导入：
`go get https://github.com/yantao1995/ds`

### 队列类

```
import (
	"github.com/yantao1995/ds/queue"
)
```

#### loopQueue 循环队列

- 限定长度的循环队列

创建循环队列对象示例: `lq, err := queue.NewLoopQueue(10)`


### 链表类

#### SkipList 跳表  

```
import (
	"github.com/yantao1995/ds/skiplist"
)
```

- 实现比较接口，就可以实现升序跳表或者降序跳表 （compare内包含个别实现示例，可以实现复杂的比较逻辑）
```
//比较接口原型
type CompareAble interface {
	//Compare 函数签名  (建议调用侧保证key不为nil,否则需要在函数内实现对nil的处理)
	Compare(a, b interface{}) int // -1 a<b  0 a==b  1 a>b
}
```
- 支持重复key元素与无重复key插入  
- 迭代器层数遍历简单输出跳表结构图

创建跳表对象示例:  `skiplist, err := skiplist.New(&skiplist.CmpInstanceStruct{})` 
