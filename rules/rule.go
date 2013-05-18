package rules

import "fmt"
import "strings"
import "path"
import "path/filepath"
import "os"

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
func PrintRuleList(rule_list []Rule) {
	for _, rule := range rule_list {
		rule.print_rule_data()
	}
}

func ReprRuleList(rule_list []Rule) string {
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
	// rule_path is path from directory executing sgrep in.  eg.,
	// if have .sgrep folder in parent directory, then rule_path
	// will be ".."
	rule_path            string
	original_rule_text   string
	rule_from_parent_dir bool
}

/**
FIXME!!!  Will only work on select operating systems (eg., no windows).
*/
func is_wildcard(rule_part string) bool {
	rule_part = strings.TrimSpace(rule_part)
	// FIXME: only useful for limited file system types
	return ((rule_part == "./*") || (rule_part == "*"))
}

func begins_with_wildcard(file_string string) bool {
	if len(file_string) == 0 {
		return false
	}

	if is_wildcard(string(file_string[0])) {
		return true
	}

	if len(file_string) < 3 {
		return false
	}

	return is_wildcard(file_string[0:3])
}

func (er ExcludeRule) Grep_arg_root_rule() string {
	return "--exclude=" + strings.TrimSpace(er.original_rule_text)
}

/**
Compare rule path to current working directory.  For each .. in
rule_path, produce the name of the folder that would have been
required to get from er.rule_path to the pwd.

Example, if pwd is
/a/b/c/d/
and rule_path is ../../, then should return "c/d"

Use this information when composing the grep arg rule for rules that
originated from parent directories.  Se note in Grep_arg_rule
function.

Note: should only be called on ExcludeRules generated from parent
directories.
*/
func (er ExcludeRule) sgrep_dir_to_pwd_dir_difference() string {
	if !er.rule_from_parent_dir {
		panic("Cannot construct path differences from subdirectories")
	}

	abs_rule_path, err := filepath.Abs(er.rule_path)
	if err != nil {
		panic("Error with rule path when calculating differences")
	}
	pwd, err := os.Getwd()
	if err != nil {
		panic("Error getting pwd when calculating differences")
	}
	abs_pwd_path, err := filepath.Abs(pwd)
	if err != nil {
		panic("Error with pwd path when calculating differences")
	}

	// list_rule_path := filepath.SplitList(abs_rule_path)
	// list_pwd_path := filepath.SplitList(abs_pwd_path)
	list_rule_path := strings.Split(
		filepath.Clean(filepath.ToSlash(abs_rule_path)),
		"/")
	list_pwd_path := strings.Split(
		filepath.Clean(filepath.ToSlash(abs_pwd_path)),
		"/")

	// list_rule_path should be a sub-list of list_pwd_path
	// (because we're guaranteed to be in a parent directory).
	// Therefore, can just get the difference by appending all
	// elements that are in list_pwd_path that are not in
	// list_rule_path.

	difference_list := list_pwd_path[len(list_rule_path):]
	difference_str := ""

	for index, element := range difference_list {
		if index == 0 {
			difference_str = element
		} else {
			difference_str = path.Join(difference_str, element)
		}
	}

	return difference_str + string(filepath.Separator)
}

func (er ExcludeRule) Grep_arg_rule() string {
	if !er.rule_from_parent_dir {
		// the rule was not gathered from a parent directory,
		// but rather the directory from which sgrep was run
		// or a subdirectory of that directory.
		return path.Join(
			er.rule_path, strings.TrimSpace(er.original_rule_text))
	}

	// algorithm: Start with rule statement.  Deal with two cases:
	//   1) Rule statement does not begin with a *:

	//      Find the difference between the cwd and the directory
	//      that the .sgrep file is located.
	//
	//      Take the rule in the file.  For instance, the rule
	//      a/b/*py.  First try to match the entire rule to the
	//      difference.  If that works, then just grab the last
	//      element from the difference.  If it does not, then
	//      take the subdirectory: a/b/ and compare it to the
	//      difference.  If that works, then the rule can be *py.
	//      If it doesn't, then try to match a/.  If that works,
	//      then the rule is b/*py.  If it does not, then emit no
	//      rule, because the rule was designed for a different
	//      folder.

	// distance from cwd to .sgrep file folder
	difference := er.sgrep_dir_to_pwd_dir_difference()

	// rule_to_use will contain the text that we should use to
	// filter grep.  (Ie, text we in as arg to --exclude=)
	rule_to_use := ""
	remaining_rule := er.original_rule_text
	for {
		dir, file_string := filepath.Split(remaining_rule)
		rule_to_use = path.Join(file_string, rule_to_use)

		did_match, err := filepath.Match(dir, difference)
		if err != nil {
			panic("Aborted on match when producing grep arg rule.")
		}

		if did_match {
			// means that we should use this rule because
			// it matched the directory we're using.
			last_element := filepath.Dir(dir)
			if is_wildcard(last_element) {
				// if the reason we matched was
				// because the last rule considered
				// was just a wildcard, preserve the
				// wildcard.
				rule_to_use = path.Join(last_element, rule_to_use)
			}
			break
		}

		if begins_with_wildcard(file_string) {
			break
		}

		if dir == remaining_rule {
			return ""
		}
		remaining_rule = dir
	}

	return rule_to_use
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
func ParseRule(single_line, dir_abs_path string, line_no uint32, rule_from_parent_dir bool) Rule {

	// ignore comments
	comment_index := strings.Index(single_line, "#")
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
	ex_rule.rule_from_parent_dir = rule_from_parent_dir
	return ex_rule
}
