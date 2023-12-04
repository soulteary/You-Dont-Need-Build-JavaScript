package optimizer

import (
	"bytes"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/js"
)

func Exec(rootDir string) {
	src := path.Join(rootDir, "santd@1.1.3/santd.js")
	dest := path.Join(rootDir, "santd@1.1.3/santd.min.js")
	PatchSantd(src, dest)

	const dayjsDir = "core@2023.12.04"

	A(rootDir, dayjsDir)

}

func A(rootDir string, outputDir string) {
	dayjs := []string{
		"dayjs@1.11.10/dayjs.min.js",
		"dayjs@1.11.10/locale/en.min.js",
		"dayjs@1.11.10/locale/zh-cn.min.js",
		"dayjs@1.11.10/plugin/utc.min.js",
		"dayjs@1.11.10/plugin/localeData.min.js",
		"dayjs@1.11.10/plugin/customParseFormat.min.js",
		"dayjs@1.11.10/plugin/weekOfYear.min.js",
		"dayjs@1.11.10/plugin/weekYear.min.js",
		"dayjs@1.11.10/plugin/advancedFormat.min.js",
	}

	buff := [][]byte{}
	for _, file := range dayjs {
		src := path.Join(rootDir, file)
		buf, err := os.ReadFile(src)
		if err != nil {
			panic(err)
		}
		buff = append(buff, buf)
	}

	os.MkdirAll(path.Join(rootDir, outputDir), 0755)
	dayjsFilePath := path.Join(outputDir, "dayjs.min.js")
	os.WriteFile(path.Join(rootDir, dayjsFilePath), bytes.Join(buff, []byte("\n")), 0644)

	cores := []string{
		dayjsFilePath,
		"enquire.js@2.1.6/enquire.min.js",
		"san@3.13.3/san.min.js",
		"san-router@2.0.2/san-router.min.js",
		"san-router@2.0.2/san-router.min.js",
		// "santd@1.1.3/santd.min.js",
		// "esljs@2.2.2/esl.min.js",
	}

	buff = [][]byte{}
	for _, file := range cores {
		src := path.Join(rootDir, file)
		buf, err := os.ReadFile(src)
		if err != nil {
			panic(err)
		}
		buff = append(buff, buf)
	}

	os.MkdirAll(path.Join(rootDir, outputDir), 0755)
	os.WriteFile(path.Join(rootDir, path.Join(outputDir, "core.min.js")), bytes.Join(buff, []byte("\n")), 0644)
	os.Remove(dayjsFilePath)
}

func MinifyScript(input string) []byte {
	m := minify.New()
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)

	result, err := m.String("text/javascript", input)
	if err != nil {
		panic(err)
	}
	return []byte(result)
}

func PatchSantd(src, dest string) {
	code, err := os.ReadFile(src)
	if err != nil {
		panic(err)
	}
	input := string(code)

	var fixSan = regexp.MustCompile(`var\ssan__default\s=.+;`)
	content := fixSan.ReplaceAllString(input, "var san__default = san;")
	var fixDayjs = regexp.MustCompile(`var\sdayjs__default\s=.+;`)
	content = fixDayjs.ReplaceAllString(content, "var dayjs__default = dayjs;")

	// import: dayjs
	content = regexp.MustCompile(`utc\s=\sutc.+?;`).ReplaceAllString(content, "dayjs.extend(window.dayjs_plugin_utc);")
	content = regexp.MustCompile(`localeData\s=\slocaleData.+?;`).ReplaceAllString(content, "dayjs.extend(window.dayjs_plugin_localeData);")
	content = regexp.MustCompile(`customParseFormat\s=\scustomParseFormat.+?;`).ReplaceAllString(content, "dayjs.extend(window.dayjs_plugin_customParseFormat);")
	content = regexp.MustCompile(`weekOfYear\s=\sweekOfYear.+?;`).ReplaceAllString(content, "dayjs.extend(window.dayjs_plugin_weekOfYear);")
	content = regexp.MustCompile(`weekYear\s=\sweekYear.+?;`).ReplaceAllString(content, "dayjs.extend(window.dayjs_plugin_weekYear);")
	content = regexp.MustCompile(`advancedFormat\s=\sadvancedFormat.+?;`).ReplaceAllString(content, "dayjs.extend(window.dayjs_plugin_advancedFormat);")
	// remove: dayjs extend
	content = regexp.MustCompile(`dayjs__default.extend(.+);`).ReplaceAllString(content, "")
	// remove: dayjs locale require
	content = regexp.MustCompile(`require\("dayjs/locale/".+\);`).ReplaceAllString(content, "")

	// fix: enquire = require('enquire.js');
	content = strings.ReplaceAll(content, "require('enquire.js');", "window.enquire;")
	// remove: process.env.NODE_ENV !== 'production' &&
	content = strings.ReplaceAll(content, "process.env.NODE_ENV !== 'production' && ", "")

	os.WriteFile(dest, MinifyScript(content), 0644)
}
