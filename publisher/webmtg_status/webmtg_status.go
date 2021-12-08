package webmtg_status

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func GetStatus() bool {
	return getStatusByTools([]string{"Slack", "zoom\\.us", "Microsoft", "Google"}, "grep -E 192\\.168\\.0\\.[0-9]+:[0-9]+$") ||
		getStatusByTools([]string{"Slack"}, "grep -E 10\\.0\\.0\\.[0-9]+:[0-9]+$") ||
		getStatusByTools([]string{"Discord"}, "grep -E '\\*:[0-9]+$'")
}

func getStatusByTools(toolNames []string, optionalGrep string) bool {
	cmdstr := "lsof -iUDP | grep"
	for _, name := range toolNames {
		cmdstr += fmt.Sprintf(" -e %s", name)
	}
	cmdstr += " | "
	cmdstr += optionalGrep
	cmdstr += " | wc -l"
	fmt.Println(cmdstr)

	out, err := exec.Command("sh", "-c", cmdstr).Output()
	if err != nil {
		log.Println(err)
		return false
	}

	num, _ := strconv.Atoi(strings.TrimSpace(string(out)))

	return num > 0
}
