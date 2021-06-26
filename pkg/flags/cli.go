package flags

import (
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/spf13/cobra"
)

var (
	noColorFlag   = "no-color"
	namespaceFlag = "namespace"
)

type CliOpts struct {
	NoColoring    bool
	AllNameSpaces bool
	Namespace     string
	AskOpts       survey.AskOpt
}

func NewCliOptions(cmd *cobra.Command) (*CliOpts, error) {
	var err error
	c := &CliOpts{
		AskOpts: func(opt *survey.AskOptions) error {
			opt.Stdio = terminal.Stdio{
				In:  os.Stdin,
				Out: os.Stdout,
				Err: os.Stderr,
			}
			return nil
		},
	}
	c.NoColoring, err = cmd.Flags().GetBool(noColorFlag)
	if err != nil {
		return nil, err
	}
	c.Namespace, err = cmd.Flags().GetString(namespaceFlag)
	if err != nil {
		return nil, err
	}
	return c, err
}

func (c *CliOpts) Ask(qs []*survey.Question, ans interface{}) error {
	return survey.Ask(qs, ans, c.AskOpts)
}

func AddPacCliOptions(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolP(noColorFlag, "C", false, "disable coloring")
}
