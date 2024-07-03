package bootstrap

import "runtime"

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
