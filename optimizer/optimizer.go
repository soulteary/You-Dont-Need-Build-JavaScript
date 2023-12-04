package optimizer

import (
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
