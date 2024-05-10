package bootstrap

import (
	"context"
	"github.com/bpcoder16/Water/conf"
	"github.com/bpcoder16/Water/utils"
	"log"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	// 重定向 stderr 和 stdout 到指定文件
	// stderr 输出到 log/std/stderr.log
	// stdout 输出到 log/std/stdout.log
	// 未 recover 的 panic 以及一些其他的 crash 信息都会输出到 stderr 里去,
	// 所以应对 stderr 监控
	// 对于线上应用，若不将 stderr 和 stdout 重定向，运行容器会将一般会将其重定向，
	hookStd()
	log.Println("stderr and stdout will redirect to log/std/")
}

func MustInit(ctx context.Context, conf *conf.AppConfig) {
	initLoggers(ctx, conf)
}

// RegisterDict 优势引入简单，劣势 key 类型大小写不敏感
func RegisterDict(dictFile string, configPtr interface{}) {
	err := utils.ParseJSONFile("./conf/dicts/"+dictFile, configPtr)
	if err != nil {
		panic("RegisterDict[" + dictFile + "], Err:" + err.Error())
	}
}
