// --
// Copyright (c) 2014 Thaniyarasu Kannusamy <thaniyarasu@gmail.com>.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
// ++

package common

import (
	"flag"
	"fmt"
	"github.com/arasuresearch/arasu/app"
	"net"
	"os"
	"os/exec"
	"time"
)

type SubCmd struct {
	Title string
	Short string
	Long  string
	Flag  flag.FlagSet
	Args  []string
	Help  bool
	App   *app.App
}

func StartAssetServer(dir, mode, port string) error {
	//fmt.Print("starting static server ,please wait")
	args := []string{"serve", "--mode", mode, "--port", port}
	//args := []string{"serve", "--port", port}

	fmt.Println(args)
	cmd := exec.Command("pub", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	err := cmd.Start()
	if err != nil {
		return err
	}
	dialUrl := ":" + port
	for {
		time.Sleep(300 * time.Millisecond)
		if _, err := net.Dial("tcp", dialUrl); err == nil {
			break
		}
	}
	return nil
}
