package ReadSgrep

import "SgrepRules"
import "os"
import "path/filepath"
import "io"
import "bufio"

/**
 @param {String} dir_abs_path --- The file path relative to the current
 directory.  

 @param {SgrepRules.RulePriority} --- The current directory has a
 priority of 0.  Directories closer to the root have lower priorities.
 Directories further from the root have higher priorities.

 @returns{List of SgrepRules} --- All rules that were read from the
 .sgrep file located in dir_abs_path.
*/
func ReadSgrepFile(
	dir_abs_path string, priority SgrepRules.RulePriority) []  SgrepRules.Rule {

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
				single_line,dir_abs_path,priority,line_no)

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