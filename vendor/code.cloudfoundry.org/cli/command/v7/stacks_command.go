package v7

import (
	"sort"

	"code.cloudfoundry.org/cli/actor/sharedaction"
	"code.cloudfoundry.org/cli/actor/v7action"
	"code.cloudfoundry.org/cli/command"
	"code.cloudfoundry.org/cli/command/v7/shared"
	"code.cloudfoundry.org/cli/util/sorting"
	"code.cloudfoundry.org/cli/util/ui"
)

//go:generate counterfeiter . StacksActor

type StacksActor interface {
	GetStacks() ([]v7action.Stack, v7action.Warnings, error)
}

type StacksCommand struct {
	usage           interface{} `usage:"CF_NAME stacks"`
	relatedCommands interface{} `related_commands:"app, push"`

	UI          command.UI
	Config      command.Config
	SharedActor command.SharedActor
	Actor       StacksActor
}

func (cmd *StacksCommand) Setup(config command.Config, ui command.UI) error {
	cmd.UI = ui
	cmd.Config = config
	cmd.SharedActor = sharedaction.NewActor(config)

	ccClient, _, err := shared.NewClients(config, ui, true, "")
	if err != nil {
		return err
	}
	cmd.Actor = v7action.NewActor(ccClient, config, nil, nil)

	return nil
}

func (cmd StacksCommand) Execute(args []string) error {
	err := cmd.SharedActor.CheckTarget(false, false)
	if err != nil {
		return err
	}

	user, err := cmd.Config.CurrentUser()
	if err != nil {
		return err
	}

	cmd.UI.DisplayTextWithFlavor("Getting stacks as {{.Username}}...", map[string]interface{}{
		"Username": user.Name,
	})
	cmd.UI.DisplayNewline()

	stacks, warnings, err := cmd.Actor.GetStacks()
	cmd.UI.DisplayWarnings(warnings)
	if err != nil {
		return err
	}

	sort.Slice(stacks, func(i, j int) bool { return sorting.LessIgnoreCase(stacks[i].Name, stacks[j].Name) })

	cmd.displayTable(stacks)

	return nil
}

func (cmd StacksCommand) displayTable(stacks []v7action.Stack) {
	if len(stacks) > 0 {
		var keyValueTable = [][]string{
			{"name", "description"},
		}
		for _, stack := range stacks {
			keyValueTable = append(keyValueTable, []string{stack.Name, stack.Description})
		}

		cmd.UI.DisplayTableWithHeader("", keyValueTable, ui.DefaultTableSpacePadding)
	}
}
