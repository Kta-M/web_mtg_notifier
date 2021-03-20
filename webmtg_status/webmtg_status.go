package webmtg_status

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func GetStatus() bool {
	return getStatusByTools([]string{"Slack", "zoom\\.us", "Microsoft", "Google"}, true) ||
		getStatusByTools([]string{"Discord"}, false)
}

func getStatusByTools(toolNames []string, useFilter bool) bool {
	cmdstr := "lsof -iUDP"
	for _, name := range toolNames {
		cmdstr += fmt.Sprintf(" | grep -e %s", name)
	}
	if useFilter {
		cmdstr += " | grep -E 192\\.168\\.0\\.5:\\[0-9]+$"
	}
	cmdstr += " | wc -l"

	out, err := exec.Command("sh", "-c", cmdstr).Output()
	if err != nil {
		log.Println(err)
		return false
	}
	num, _ := strconv.Atoi(strings.TrimSpace(string(out)))

	return num > 0
}
