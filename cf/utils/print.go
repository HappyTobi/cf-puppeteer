package print

import (
	"bytes"
	"encoding/json"
	"github.com/happytobi/cf-puppeteer/ui"
)

//PrettyJSON prints pretty json
func PrettyJSON(jsonUgly string) error {
	jsonPretty := &bytes.Buffer{}
	err := json.Indent(jsonPretty, []byte(jsonUgly), "", "  ")

	if err != nil {
		ui.Failed("PrettyJSON error %s", err)
		return err
	}

	//ui.Say(jsonPretty.String())

	return nil
}
