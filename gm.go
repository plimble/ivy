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

//GMProcessor is GraphicsMagick processor
type GMProcessor struct{}

//NewGMProcessor create GraphicsMagick processor
func NewGMProcessor() *GMProcessor {
	return &GMProcessor{}
}

//Process process image
func (gm *GMProcessor) Process(params *params, filePath string, file *bytes.Buffer) (*bytes.Buffer, error) {
	if params.isDefault {
		return file, nil
	}

	gmb := newGMBuilder()
	gmb.Strip()

	if params.enableResize {
		switch {
		case params.width > 1 && params.height > 1:
			gmb.Resize(params.width, params.height, "!")
		case params.width < 1 && params.height > 1:
			gmb.Resize(1, params.height, "^")
		case params.width > 1 && params.height < 1:
			gmb.Resize(params.width, 1, "^")
		}
	}

	if params.enableGravity {
		gmb.Gravity(gm.getGravity(params.gravity))
	}

	if params.enableCrop {
		gmb.Crop(params.cropWidth, params.cropHeight)
	}

	if params.quality != -1 {
		gmb.Quality(params.quality)
	}

	out := &bytes.Buffer{}
	err := gmb.Process(file, out)

	return out, err
}

func (gm *GMProcessor) getGravity(cropPos string) string {
	switch cropPos {
	case "nw":
		return "NorthWest"
	case "n":
		return "North"
	case "ne":
		return "NorthEast"
	case "w":
		return "West"
	case "c":
		return "Center"
	case "e":
		return "East"
	case "sw":
		return "SouthWest"
	case "s":
		return "South"
	case "se":
		return "SouthEast"
	}

	return "NorthWest"
}

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
