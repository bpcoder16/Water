package bootstrap

import (
	"github.com/bpcoder16/Water/env"
	"os"
	"path"
	"path/filepath"
	"syscall"
)

// HookStd 同时劫持 Stderr 和 Stdout
func hookStd(logPath string) {
	hookStderr(logPath)
	hookStdout(logPath)
}

// hookStderr 劫持 Stderr
func hookStderr(logPath string) {
	var filename string
	if len(logPath) > 0 {
		filename = logPath + "/std/stderr.log"
	} else {
		filename = path.Join(env.RootPath(), "log/std/stderr.log")
	}
	dirname := filepath.Dir(filename)
	if err := os.MkdirAll(dirname, 0755); err != nil {
		panic("mkdirAll " + dirname + " failed")
	}
	fh, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic("open stderr.log failed: " + err.Error())
	}
	err = syscall.Dup2(int(fh.Fd()), 2)
	if err != nil {
		panic("stderr.log syscall.Dup2 failed: " + err.Error())
	}
}

// hookStdout 劫持 Stdout
func hookStdout(logPath string) {
	var filename string
	if len(logPath) > 0 {
		filename = logPath + "/std/stdout.log"
	} else {
		filename = path.Join(env.RootPath(), "log/std/stdout.log")
	}
	dirname := filepath.Dir(filename)
	if err := os.MkdirAll(dirname, 0755); err != nil {
		panic("mkdirAll " + dirname + " failed")
	}
	fh, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic("open stdout.log failed: " + err.Error())
	}
	err = syscall.Dup2(int(fh.Fd()), 1)
	if err != nil {
		panic("stdout.log syscall.Dup2 failed: " + err.Error())
	}
}
