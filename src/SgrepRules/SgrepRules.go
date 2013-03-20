package SgrepRules

import "fmt"
import "strings"
import "path"

/**
 Any rule that is contained in a .sgrep file should match the Rule interface.

 Currently, we only support one type of rule, ExcludeRule, which tells
 us to skip a file when sgrep-ing.
*/


type RuleType uint32
const EXCLUDE = 0

type Rule interface {

	// Currently, just have exclude rules
	type_of_rule() RuleType

	// returns 
	Grep_arg_root_rule() string
	Grep_arg_rule() string
	
	// For debugging:
	
	// should return the original text of the rule.  
	text() string
	// should return the relative location of the sgrep file that
	// the rule was created from
	rule_location() string
	// just pretty prints the rule's internal data
	print_rule_data()

	repr_rule_data() string
	
}

/*
 * Pretty prints a list of rule objects.
 */
func PrintRuleList (rule_list [] Rule)  {
	for _, rule := range rule_list {
		rule.print_rule_data()
	}
}

func ReprRuleList (rule_list [] Rule) string {
	str := ""
	for _, rule := range rule_list {
		str += rule.repr_rule_data() + "\n"
	}
	return str
}



/**
 * Specifies a rule for skipping files/folders when sgrep-ing.
 */
type ExcludeRule struct {
	rule_path string
	original_rule_text string
}

func (er ExcludeRule) Grep_arg_root_rule() string {
	return "--exclude=" + strings.TrimSpace(er.original_rule_text)
}

func (er ExcludeRule) Grep_arg_rule() string {
	return "--exclude=" + path.Join(er.rule_path, strings.TrimSpace(er.original_rule_text))
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
	fmt.Println(er.repr_rule_data())
}

func (er ExcludeRule) repr_rule_data() string {
	str := "ExcludeRule \n"
	str += "\t" + er.original_rule_text + "\n"
	str += "\t" + er.rule_path + "\n"
	return str
}


/*
This function originally returned an interface pointer.
But kept throwing an error.  Plus, found this doc:
https://groups.google.com/forum/?fromgroups=#!topic/golang-nuts/DwFGXLYgatY
*/
func ParseRule(single_line,dir_abs_path string, line_no uint32)  Rule{

	// ignore comments
	comment_index := strings.Index(single_line,"#")
	if comment_index != -1 {
		single_line = single_line[0:comment_index]
	}
	single_line = strings.TrimSpace(single_line)
	
	// ignore blank lines
	if single_line == "" {
		return nil
	}

	ex_rule := ExcludeRule{}
	ex_rule.rule_path = dir_abs_path
	ex_rule.original_rule_text = single_line

	return ex_rule
}