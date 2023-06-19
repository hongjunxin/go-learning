#!/bin/bash
go build -buildmode=c-archive -o obj/main.a

# 生成 main.a 静态库的同时，也会生成一个 main.h 头文件。但是 main.h 头文件中
# 只有 main 包中导出的 goPrintln 函数的声明，并没有 number 子包导出函数的声明。
# 其实 number_add_mod 函数在生成的 C 静态库中是存在的，我们可以直接使用。