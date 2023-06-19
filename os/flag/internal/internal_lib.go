package internal

import "fmt"

/*
访问权限：模块级私有
internal代码包中声明的公开程序实体仅能被该代码包的直接父包及其子包中的代码引用。
当然，引用前需要先导入这个internal包。对于其他代码包，导入该internal包都是非法的，无法通过编译。

该示例中，internal.Hello() 只能被 flag 目录下的文件访问
*/
func Hello() {
	fmt.Println("Hello")
}
