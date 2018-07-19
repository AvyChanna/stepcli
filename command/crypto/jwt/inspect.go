package jwt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/smallstep/cli/command/crypto/internal/utils"
	"github.com/smallstep/cli/errs"
	"github.com/urfave/cli"
)

func inspectCommand() cli.Command {
	return cli.Command{
		Name:   "inspect",
		Action: cli.ActionFunc(inspectAction),
		Usage:  `return the decoded JWT without verification`,
		UsageText: `**step crypto jwt inspect**
		**--insecure**`,
		Description: `**step crypto jwt inspect** reads a JWT data structure from STDIN, decodes it,
and outputs the header and payload on STDERR. Since this command does not
verify the JWT you must pass **--insecure** as a misuse prevention mechanism.`,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:   "insecure",
				Hidden: true,
			},
		},
	}
}

func inspectAction(ctx *cli.Context) error {
	if err := errs.NumberOfArguments(ctx, 0); err != nil {
		return err
	}

	if !ctx.Bool("insecure") {
		return errs.InsecureCommand(ctx)
	}

	token, err := utils.ReadString(os.Stdin)
	if err != nil {
		return err
	}

	return printToken(token)
}

func printToken(token string) error {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return errors.New("error decoding token: JWT must have three parts")
	}

	header, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return errors.Wrapf(err, "error decoding token")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return errors.Wrapf(err, "error decoding token")
	}

	m := make(map[string]json.RawMessage)
	m["header"] = header
	m["payload"] = payload
	m["signature"] = []byte(`"` + parts[2] + `"`)

	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return errors.Wrapf(err, "error marshaling token data")
	}
	fmt.Fprintln(os.Stderr, string(b))

	return nil
}