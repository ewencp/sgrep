package SgrepRules

import "fmt"

/**
 Any rule that is contained in a .sgrep file should match the Rule interface.

 Currently, we only support one type of rule, ExcludeRule, which tells
 us to skip a file when sgrep-ing.
*/


type RuleType uint32
const EXCLUDE = 0

type RulePriority int32

type Rule interface {

	// Currently, just have exclude rules
	type_of_rule() RuleType
	// For two conflicting rules, choose to follow the one with
	// the higher priority.
	priority () RulePriority

	
	// For debugging:
	
	// should return the original text of the rule.  
	text() string
	// should return the relative location of the sgrep file that
	// the rule was created from
	rule_location() string
	// just pretty prints the rule's internal data
	print_rule_data()
	
}

/*
 * Pretty prints a list of rule objects.
 */
func PrintRuleList (rule_list [] Rule)  {
	for _, rule := range rule_list {
		rule.print_rule_data()
	}
}


/**
 * Specifies a rule for skipping files/folders when sgrep-ing.
 */
type ExcludeRule struct {
	rule_path string
	priority_ RulePriority
	original_rule_text string
}

func (er ExcludeRule) priority() RulePriority {
	return er.priority_
}

func (er ExcludeRule) text() string {
	return er.original_rule_text
}

func (er ExcludeRule) rule_location() string {
	return er.rule_path
}

func (er ExcludeRule) type_of_rule() RuleType {
	return EXCLUDE
}

func (er ExcludeRule) print_rule_data() {

	fmt.Println("ExcludeRule")
	fmt.Println("\t" + er.original_rule_text)
	fmt.Println("\t" + er.rule_path)
}

/*
This function originally returned an interface pointer.
But kept throwing an error.  Plus, found this doc:
https://groups.google.com/forum/?fromgroups=#!topic/golang-nuts/DwFGXLYgatY
*/
func ParseRule(single_line,dir_abs_path string, priority RulePriority)  Rule{
	if single_line == "" {
		return nil
	}
	
	ex_rule := ExcludeRule{}
	ex_rule.rule_path = dir_abs_path
	ex_rule.priority_ = priority
	ex_rule.original_rule_text = single_line

	return ex_rule
}