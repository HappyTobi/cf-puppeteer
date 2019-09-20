package env

import (
	"encoding/csv"
	"github.com/happytobi/cf-puppeteer/ui"
	"strings"
)

func Convert(envVars []string) (conVars map[string]string) {
	replacer := strings.NewReplacer("'", "\"", "`", "\"")
	envs := make(map[string]string)
	for _, v := range envVars {
		fv := replacer.Replace(v)
		r := csv.NewReader(strings.NewReader(fv))
		r.Comma = '='
		r.TrimLeadingSpace = true

		fields, err := r.Read()
		if err != nil {
			ui.Warn("could not convert environment variable %s - error:%s", v, err)
		}

		if len(fields) > 1 {
			envKey := strings.TrimSpace(fields[0])
			envVal := strings.TrimSpace(fields[1])
			if len(fields) > 2 {
				//only possible is value contains "=" sign so we have to join them together again
				envVal = strings.Join(fields[1:], "=")
			}

			envs[envKey] = envVal
		}
	}

	return envs
}
