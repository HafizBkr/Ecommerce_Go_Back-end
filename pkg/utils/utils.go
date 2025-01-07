package utils

import "strings"

func SplitName(fullName string) []string {
    return strings.Fields(strings.TrimSpace(fullName))
}
