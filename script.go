/*
 * @Author: xiongxinwei 3293172751nss@gmail.com
 * @Date: 2022-07-30 13:51:21
 * @LastEditors: xiongxinwei 3293172751nss@gmail.com
 * @LastEditTime: 2022-07-30 13:58:10
 * @FilePath: \undefinedd:\文档\git\计算机操作系统\script.go
 * @Description: 计算机操作系统的学习
 *
 */
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {

	var a int = 24
	for i := 24; i < 100; i++ {

		a1 := strconv.Itoa(a)
		a2 := strconv.Itoa((a + 1))
		a3 := strconv.Itoa((a - 1))
		filePath := "markdown/" + a1 + ".md"
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
		//在原来的基础上追加666表示访问权限
		if err != nil {
			fmt.Println("文件打开失败", err)
		}
		//及时关闭file句柄
		defer file.Close()

		//写入文件时，使用带缓存的 *Writer
		write := bufio.NewWriter(file)
		write.WriteString("+ [author](https://github.com/3293172751)\n")
		write.WriteString("# 第" + a1 + "节\n")

		//批量加入文件，

		write.WriteString("+ [回到目录](../README.md)\n")
		write.WriteString("+ [回到项目首页](../../README.md)\n")
		write.WriteString("+ [上一节](" + a3 + ".md)\n")
		write.WriteString("> ❤️💕💕真正动手时间操作系统,从应用软件出发\"进入操作系统\",成为掌握计算机关键技术的工程师。Myblog:[http://nsddd.top](http://nsddd.top/)\n")
		write.WriteString("---\n")
		write.WriteString("[TOC]\n")
		for i := 0; i < 5; i++ {
			write.WriteString("\n")
		}
		write.WriteString("## END 链接\n")
		write.WriteString("+ [回到目录](../README.md)\n")
		write.WriteString("+ [上一节](" + a3 + ".md)\n")
		write.WriteString("+ [下一节](" + a2 + ".md)\n")
		write.WriteString("---\n")
		write.WriteString("+ [参与贡献❤️💕💕](https://github.com/3293172751/Block_Chain/blob/master/Git/git-contributor.md)")
		//Flush将缓存的文件真正写入到文件中
		write.Flush()
		a = a + 1
	}

}
