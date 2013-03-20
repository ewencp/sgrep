package ReadSgrep

import "SgrepRules"
import "os"
import "path/filepath"
import "io"
import "bufio"
import "io/ioutil"
import "fmt"
import "strings"

type RuleTree struct {

	// the rules contained in a single folder
	cur_dir_rules [] SgrepRules.Rule

	// not the absolute path.  just the name of the folder
	// relative to its parent folder.
	dir_name string

	// additional folders to scan through with their rule
	// subtrees.
	sub_directories [] *RuleTree
}

/**
 * From the directory specified by current_dir, generate a RuleTree
 * for the current directory and all subdirectories based on the .sgrep
 * files they contain.
 */
func GetRuleTree(current_dir string) *RuleTree {
	root := new(RuleTree)
	root.dir_name = current_dir
	walk_down(current_dir,root)

	return root
}


func ProduceGrepArgs(rt *RuleTree) [] string {
	var grep_args [] string
	for _, sgrep_rule := range rt.cur_dir_rules {
		grep_args = append(grep_args,sgrep_rule.Grep_arg_rule())
	}

	for _, subdir_tree := range rt.sub_directories {
		new_args := ProduceGrepArgs(subdir_tree)
		for _, additional_arg := range new_args {
			grep_args = append(grep_args,additional_arg)
		}

	}

	return grep_args
}



/**
Mostly for debugging, prints all sgrep rules found from the root
of the rule tree provided.
*/
func  PrettyPrint (rt *RuleTree) {
	fmt.Println(rt.pretty_print_helper())
}

func (rt *RuleTree) pretty_print_helper ( ) string {

	str := "***" + rt.dir_name + "***\n"
	rule_str := SgrepRules.ReprRuleList(rt.cur_dir_rules)
	str += indent_str(rule_str, 1)

	for _, subtree := range rt.sub_directories {
		str += indent_str(
			subtree.pretty_print_helper(),
			1)
	}

	return str
}

/**
 Helper function: indents a string with the number of tabs
 specified by how_much.  */
func indent_str(original string, how_much uint32)  string {
	indented := ""
	split := strings.Split(original,"\n")
	for _, line_text := range split {
		indented += "\n\t" + line_text
	}
	return indented
}



/**

 Look for .sgrep files in current_dir and read all of its
 subdirectories looking for .sgrep files as well.

 FIXME: probably a way to cull running through certain directories early.

 FIXME: Does not handle cycles at all.
 */
func walk_down(
	current_dir string, root *RuleTree){

	// determine my current rules
	root.cur_dir_rules = ReadSgrepFile(current_dir)

	file_dir_list, _ := ioutil.ReadDir(current_dir)

	var file_dir_node os.FileInfo
	for _, file_dir_node = range file_dir_list {

		if file_dir_node.IsDir() {
			// recurse, looking for rules in this folder
			new_leaf := new(RuleTree)
			new_leaf_dir_path := filepath.Join(current_dir,file_dir_node.Name())
			new_leaf.dir_name = file_dir_node.Name()
			root.sub_directories = append(root.sub_directories,new_leaf)
			walk_down(new_leaf_dir_path,new_leaf)
		}
	}
}



/**
 @param {String} dir_abs_path --- The file path relative to the current
 directory.  

 @returns{List of SgrepRules} --- All rules that were read from the
 .sgrep file located in dir_abs_path.
*/
func ReadSgrepFile(
	dir_abs_path string) []  SgrepRules.Rule {

	var rules [] SgrepRules.Rule
	
	fi,err := os.Open(filepath.Join(dir_abs_path,".sgrep"))
	if err != nil {
		// no .sgrep file in this folder
		return rules
	}
	// at end of function close fi
	defer fi.Close()

	file_reader := bufio.NewReader(fi)
	single_line := ""
	var line_no uint32;
	line_no = 0;
	for {
		line_no += 1
		single_line, err =  file_reader.ReadString('\n')
		if single_line != "" {
			new_rule := SgrepRules.ParseRule(
				single_line,dir_abs_path,line_no)

			if new_rule != nil {
				// FIXME: I wonder if there's any
				// overhead from using slices in this
				// way.  eg, copies all elements over
				// again.
				rules = append(rules,new_rule)
			}
		}
		
		if (err == io.EOF){
			break
		}
	}

	return rules
}