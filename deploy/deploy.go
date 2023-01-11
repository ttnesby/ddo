package deploy

import (
	dl "ddo/deploy/level"
	do "ddo/deploy/operation"
	rr "ddo/reporoot"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func name(context string) string {
	sha1 := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(context))
	wo := strings.Split(sha1.String(), "-")
	return context + "-" + wo[0] + wo[4]
}

func verify(id, bicep, json string) (validId uuid.UUID, bicepPath string, jsonPath string, e error) {

	anError := func(e error) (uuid.UUID, string, string, error) {
		return uuid.Nil, "", "", e
	}

	i, err := uuid.Parse(id)
	if err != nil {
		return anError(err)
	}

	b := filepath.Join(rr.Get(), bicep)
	j := filepath.Join(rr.Get(), json)

	if _, err := os.Stat(b); err != nil {
		return anError(err)
	}
	if _, err := os.Stat(j); err != nil {
		return anError(err)
	}

	return i, b, j, nil
}

func New(
	level dl.Level,
	op do.Operation,
	context,
	id,
	rgOrLocation,
	templateFile,
	parameterFile string) (string, error) {

	if !level.Valid() || !op.Valid() {
		return "", fmt.Errorf("invalid level %s or operation %s", level, op)
	}

	i, b, j, err := verify(id, templateFile, parameterFile)
	if err != nil {
		return "", err
	}

	prefix := strings.Join(
		[]string{
			"az deployment",
			string(level),
			string(op),
			"--name",
			name(context),
		},
		" ")

	infix := func() string {
		switch level {
		case dl.ResourceGroup:
			return strings.Join(
				[]string{
					"--subscription",
					i.String(),
					"--resource-group",
					rgOrLocation,
				},
				" ")
		case dl.Subscription:
			return strings.Join(
				[]string{
					"--subscription",
					i.String(),
					"--location",
					rgOrLocation,
				},
				" ")
		case dl.ManagementGroup:
			return strings.Join(
				[]string{
					"--management-group-id",
					i.String(),
					"--location",
					rgOrLocation,
				},
				" ")
		}
		return "" // should never happen
	}()

	postfix := strings.Join(
		[]string{
			"--template-file",
			b,
			"--out",
			"yaml",
			"--parameters",
			"@" + j,
		},
		" ")

	return prefix + " " + infix + " " + postfix, nil
}
