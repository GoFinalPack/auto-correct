## auto-correct

自动给中英文之间加入合理的空格并纠正专用名词大小写

### 安装

```shell
go get -u https://github.com/GoFinalPack/auto-correct

```

### 使用

```go
	a := auto_correct.AutoCorrect{}
	a.Init()
	text := "golang 使用中文测试"   // Golang 使用中文测试
	fmt.Println(a.Correct(text))

	text = "pfinalclub测试"     //  PFinal Club 测试
	fmt.Println(a.Correct(text))

```
