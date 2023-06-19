package main

import "container/list"

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

/**
 * 代码中的类名、方法名、参数名已经指定，请勿修改，直接返回方法规定的值即可
 *
 *
 * @param root TreeNode类
 * @return int整型一维数组
 */
func preorderTraversal(root *TreeNode) []int {
	var ret []int
	var l list.List
	ret = append(ret, root.Val)
	if root.Right != nil {
		l.PushFront(root.Right)
	}
	if root.Left != nil {
		l.PushFront(root.Left)
	}

	for l.Len() > 0 {
		ele := l.Front()
		root := ele.Value.(*TreeNode)
		l.Remove(ele)
		ret = append(ret, root.Val)
		if root.Right != nil {
			l.PushFront(root.Right)
		}
		if root.Left != nil {
			l.PushFront(root.Left)
		}
	}

	return ret
}
