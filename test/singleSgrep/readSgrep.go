package main

// import "fmt"
import "ReadSgrep"
import "SgrepRules"


func main() {

	sgrep_rule_list := ReadSgrep.ReadSgrepFile(".")
	SgrepRules.PrintRuleList(sgrep_rule_list)
	// sgrep_rule_list.PrintRuleList()

}