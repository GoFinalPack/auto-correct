package tests

import (
	"fmt"
	auto_correct "github/pfinal/auto-correct"
	"testing"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2024/8/22
 * @Desc:
 * @Project: auto-correct
 */

func TestCorrect(T *testing.T) {

	a := auto_correct.AutoCorrect{}
	a.Init()
	text := "golang 使用中文测试"
	fmt.Println(a.Correct(text))

	text = "pfinalclub测试"
	fmt.Println(a.Correct(text))

	text = "json测试" //  JSON 测试
	fmt.Println(a.Correct(text))

	text = "Mysql 测试一下" //  MySQL 测试一下
	fmt.Println(a.Correct(text))
}
