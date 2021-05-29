package extend

import (
	"fmt"
	"io"
	"strings"
	"os/exec"
	"bufio"
    "path/filepath"
    "Two-Card/utils"
	"github.com/astaxie/beego"
)

func ExecCommand(commandName string, params []string) bool {
    cmd := exec.Command(commandName, params...)
    //显示运行的命令
    utils.LogDebug(fmt.Sprintf("执行: %s\n", strings.Join(cmd.Args, " ")))
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        utils.LogError(fmt.Sprintf("error=>%s", err.Error()))
        return false
    }
    cmd.Start() // Start开始执行c包含的命令，但并不会等待该命令完成即返回。Wait方法会返回命令的返回状态码并在命令返回后释放相关的资源。
 
    reader := bufio.NewReader(stdout)
 
    var contentArray []string
    var index int
    //实时循环读取输出流中的一行内容
    for {
        line, err2 := reader.ReadString('\n')
        if err2 != nil || io.EOF == err2 {
            break
        }
        index++
        contentArray = append(contentArray, line)
    }
 
	cmd.Wait()
	// LogDebug(contentArray)
    return true
}

func ExecSclite(reffile, hypfile, output_dir string) bool {
	Home := beego.AppConfig.String("homeroot")
	command := "bin/sclite"
	params := []string{
        "-i",
        "wsj",
        "-r",
        reffile,
        "-h",
        hypfile,
        "-e",
        "utf-8",
        "-o",
        "all",
        "-O",
        output_dir,
        "-c",
		"NOASCII",
	}
	ExecCommand(filepath.Join(Home, command), params)
	return true
}
