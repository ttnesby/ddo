package configuration

import (
	"bytes"
	"ddo/alogger"
	"ddo/path"
	"fmt"
	"github.com/oklog/ulid/v2"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

var l = alogger.New()

type CueCli []string

func addFlags(flag string, slice []string) []string {

	var adder func(result []string, slice []string) []string
	adder = func(result []string, slice []string) []string {
		switch len(slice) {
		case 0:
			return result
		case 1:
			return append(result, flag, slice[0])
		default:
			return adder(append(result, flag, slice[0]), slice[1:])
		}
	}

	return adder(nil, slice)
}

func New(path string, tags []string) (cmd CueCli) {
	cmd = append([]string{"cue", "export", path}, addFlags("-t", tags)...)
	l.Debugf("cueCmd: %v", cmd)
	return cmd
}

func (cueCmd CueCli) AsJson() (cmd CueCli) {
	cmd = append(cueCmd, "--out", "json")
	l.Debugf("cueCmd: %v", cmd)
	return cmd
}

func (cueCmd CueCli) AsYaml() (cmd CueCli) {
	cmd = append(cueCmd, "--out", "yaml")
	l.Debugf("cueCmd: %v", cmd)
	return cmd
}

func (cueCmd CueCli) ElementsAsJson(elements []string) (cmd CueCli) {
	cmd = append(append(cueCmd, addFlags("-e", elements)...), "--out", "json")
	l.Debugf("cueCmd: %v", cmd)
	return cmd
}

func (cueCmd CueCli) ElementsAsText(elements []string) (cmd CueCli) {
	cmd = append(append(cueCmd, addFlags("-e", elements)...), "--out", "text")
	l.Debugf("cueCmd: %v", cmd)
	return cmd
}

func (cueCmd CueCli) ElementsToTmpJsonFile(elements []string) (cmd CueCli, absolutePath string) {
	absolutePath = filepath.Join(
		os.TempDir(),
		fmt.Sprintf("ddo.parameters.%s.json", ulid.Make().String()),
	)
	cmd = append(
		append(cueCmd, addFlags("-e", elements)...),
		"--out", "json", "--outfile", absolutePath,
	)

	l.Debugf("cueCmd: %v", cmd)
	return cmd, absolutePath
}

func (cueCmd CueCli) Run() (byte []byte, e error) {

	l.Debugf("cueCmd: %v", cueCmd)

	cmd := exec.Command(cueCmd[0], cueCmd[1:]...)
	cmd.Dir = path.RepoRoot()

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(&stdoutBuf) //os.Stdout
	cmd.Stderr = io.MultiWriter(&stderrBuf) //os.Stderr

	err := cmd.Run()
	if err != nil {
		return nil, l.Error(fmt.Errorf("%v failed: %s", cueCmd, err))
	}

	out, errStr := stdoutBuf.Bytes(), string(stderrBuf.Bytes())
	if errStr != "" {
		return nil, l.Error(fmt.Errorf("%v returned: %s", cueCmd, errStr))
	}

	return bytes.TrimRight(out, "\r\n"), nil
}
