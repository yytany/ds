# ds
data structure design and achieve
数据结构设计和实现

## 队列类

### loopQueue 循环队列
- 限定长度的循环队列

## 链表类

### SkipList 跳表
- 实现比较接口，就可以实现升序跳表或者降序跳表 （compare内包含个别实现示例，可以实现复杂的比较逻辑）
```
//比较接口原型
type CompareAble interface {
	//Compare 函数签名
	Compare(a, b interface{}) int // -1 a<b  0 a==b  1 a>b
}
```
- 支持重复key元素与无重复key插入
- 迭代器层数遍历简单输出跳表结构图