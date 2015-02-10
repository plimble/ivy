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

type GMBuilder struct {
	args []string
}

func NewGMBuilder() *GMBuilder {
	return &GMBuilder{[]string{"convert"}}
}

func (g *GMBuilder) Strip() *GMBuilder {
	g.args = append(g.args, "-strip")
	return g
}

func (g *GMBuilder) Quality(qty int) *GMBuilder {
	g.args = append(g.args, "-quality", strconv.Itoa(qty))
	return g
}

func (g *GMBuilder) Resize(width, height int, option string) *GMBuilder {
	g.args = append(g.args, "-resize", strconv.Itoa(width)+"x"+strconv.Itoa(height)+option)
	return g
}

func (g *GMBuilder) Crop(width, height int) *GMBuilder {
	g.args = append(g.args, "-crop", strconv.Itoa(width)+"x"+strconv.Itoa(height)+"+0+0")
	return g
}

func (g *GMBuilder) Gravity(pos string) *GMBuilder {
	g.args = append(g.args, "-gravity", pos)
	return g
}

func (g *GMBuilder) Process(in io.Reader, out io.Writer) error {
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

func (g *GMBuilder) killCmd(cmd *exec.Cmd) error {
	if err := cmd.Process.Kill(); err != nil {
		return errs.NewErrors(fmt.Sprintf("Failed to kill command: %v", err))
	}

	return nil
}
