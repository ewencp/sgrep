package main

// import "fmt"
import "ReadSgrep"
// import "SgrepRules"


func main() {
	rule_tree := ReadSgrep.GetRuleTree(".")
	ReadSgrep.PrettyPrint(rule_tree)
}