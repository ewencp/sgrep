package main

import "ReadSgrep"

func main() {
	rule_tree := ReadSgrep.GetRuleTree(".")
	ReadSgrep.PrettyPrint(rule_tree)
}