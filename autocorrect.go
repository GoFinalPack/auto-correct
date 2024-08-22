package auto_correct

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"os"
	"regexp"
	"strings"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2024/8/22
 * @Desc:
 * @Project: auto-correct
 */

type AutoCorrect struct {
	DictsMap map[string]string
}

//go:embed dicts.txt
var defaultDicts embed.FS

func (a *AutoCorrect) Init() {
	// 加载docts
	a.DictsMap = make(map[string]string)
	envVar := "DICTPATH"
	dictsPath := "dicts.txt"
	if dictPath, ok := os.LookupEnv(envVar); ok {
		dictsPath = dictPath
		err := a.loadDicts(dictsPath)
		if err != nil {
			fmt.Printf("加载字典文件失败: %v\n", err)
		}
	} else {
		// 使用 embed.FS 加载 dicts.txt
		err := a.loadDictsFromEmbedFS("dicts.txt")
		if err != nil {
			fmt.Printf("加载嵌入的字典文件失败: %v\n", err)
		}
	}

}

func (a *AutoCorrect) loadDicts(filePath string) error {
	// 打开字典文件
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 逐行读取文件内容
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// 假设每一行的格式是 key:value
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue // 忽略格式不正确的行
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		// 将键值对存入 DictsMap
		a.DictsMap[key] = value
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (a *AutoCorrect) loadDictsFromEmbedFS(fileName string) error {
	// 从嵌入的文件系统中读取文件内容
	data, err := defaultDicts.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("读取嵌入文件失败: %v", err)
	}

	// 使用 bufio.Scanner 逐行处理文件内容
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		// 假设每一行的格式是 key:value
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue // 忽略格式不正确的行
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		// 将键值对存入 DictsMap
		a.DictsMap[key] = value
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
func (a *AutoCorrect) autoSpace(content string) string {
	// 汉字范围匹配，使用 Unicode 范围表示汉字
	re1 := regexp.MustCompile(`([\x{4e00}-\x{9fff}])([a-zA-Z0-9+$@#\[\(/‘“])`)
	content = re1.ReplaceAllString(content, `$1 $2`)

	// 英文数字和其他字符后面跟随汉字的情况
	re2 := regexp.MustCompile(`([a-zA-Z0-9+$’”\]\)@#!\/]|[\d[年月日]]{2,})([\x{4e00}-\x{9fff}])`)
	content = re2.ReplaceAllString(content, `$1 $2`)

	// 英文数字后面跟随某些符号的情况
	re3 := regexp.MustCompile(`([a-zA-Z0-9]+)([\[\(‘“])`)
	content = re3.ReplaceAllString(content, `$1 $2`)

	// 符号后面跟随英文数字的情况
	re4 := regexp.MustCompile(`([\)\]’”])([a-zA-Z0-9]+)`)
	content = re4.ReplaceAllString(content, `$1 $2`)

	return content
}

func (a *AutoCorrect) autoCorrect(content string) string {
	for from, to := range a.DictsMap {
		escapedFrom := regexp.QuoteMeta(from)
		pattern := fmt.Sprintf(`(?:^|[^a-z.])(%s)(?:[^a-z.]|$)`, escapedFrom)
		re := regexp.MustCompile(pattern)

		// 替换匹配的内容
		result := re.ReplaceAllStringFunc(content, func(match string) string {
			// 查找匹配的捕获组
			// 需要处理分组匹配
			// match[1] 是实际匹配的字符串
			return strings.Replace(match, from, to, 1)
		})

		content = result
	}
	return content
}

func (a *AutoCorrect) Correct(text string) string {
	// 将输入的文本正则匹配
	text = a.autoSpace(text)
	text = a.autoCorrect(text)
	return text
}
