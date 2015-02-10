package ivy

import (
	"bytes"
	"fmt"
	"github.com/plimble/errs"
	"io"
	"os/exec"
	"strconv"
	"time"
)

type gmBuilder struct {
	args []string
}

func newGMBuilder() *gmBuilder {
	return &gmBuilder{[]string{"convert"}}
}

func (g *gmBuilder) Strip() *gmBuilder {
	g.args = append(g.args, "-strip")
	return g
}

func (g *gmBuilder) Quality(qty int) *gmBuilder {
	g.args = append(g.args, "-quality", strconv.Itoa(qty))
	return g
}

func (g *gmBuilder) Resize(width, height int, option string) *gmBuilder {
	g.args = append(g.args, "-resize", strconv.Itoa(width)+"x"+strconv.Itoa(height)+option)
	return g
}

func (g *gmBuilder) Crop(width, height int) *gmBuilder {
	g.args = append(g.args, "-crop", strconv.Itoa(width)+"x"+strconv.Itoa(height)+"+0+0")
	return g
}

func (g *gmBuilder) Gravity(pos string) *gmBuilder {
	g.args = append(g.args, "-gravity", pos)
	return g
}

func (g *gmBuilder) Process(in io.Reader, out io.Writer) error {
	g.args = append(g.args, "-", "-")
	cmd := exec.Command("gm", g.args...)

	stderr := bytes.NewBuffer(nil)
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = stderr

	cmd.Start()
	cmdDone := make(chan error, 1)
	go func() {
		cmdDone <- cmd.Wait()
	}()

	select {
	case <-time.After(time.Duration(500000) * time.Millisecond):
		err := g.killCmd(cmd)
		if err != nil {
			return err
		}
		<-cmdDone
		return errs.NewErrors("Command timed out")
	case err := <-cmdDone:
		if err != nil {
			return errs.NewErrors(stderr.String())
		}
	}

	return nil
}

func (g *gmBuilder) killCmd(cmd *exec.Cmd) error {
	if err := cmd.Process.Kill(); err != nil {
		return errs.NewErrors(fmt.Sprintf("Failed to kill command: %v", err))
	}

	return nil
}
