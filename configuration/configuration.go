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

func New(path string, tags []string) CueCli {
	return append([]string{"cue", "export", path}, addFlags("-t", tags)...)
}

func (cueCmd CueCli) AsJson() (byte []byte, e error) {
	return append(cueCmd, "--out", "json").Run()
}

func (cueCmd CueCli) AsJsonCmd() (str []string) {
	return append(cueCmd, "--out", "json")
}

func (cueCmd CueCli) AsYaml() (byte []byte, e error) {
	return append(cueCmd, "--out", "yaml").Run()
}

func (cueCmd CueCli) ElementsAsJson(elements []string) (byte []byte, e error) {
	return append(append(cueCmd, addFlags("-e", elements)...), "--out", "json").Run()
}

func (cueCmd CueCli) ElementsAsText(elements []string) (byte []byte, e error) {
	return append(append(cueCmd, addFlags("-e", elements)...), "--out", "text").Run()
}

func (cueCmd CueCli) ElementsToTmpJsonFile(elements []string) (absolutePath string, e error) {
	absPath := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("ddo.parameters.%s.json", ulid.Make().String()),
	)
	_, err := append(
		append(cueCmd, addFlags("-e", elements)...),
		"--out", "json", "--outfile", absPath,
	).Run()

	if err != nil {
		return "", err
	}
	return absPath, nil
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
