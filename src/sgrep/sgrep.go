package main

import "ReadSgrep"
import "os"
import "os/exec"
import "io"
import "os/user"
import "log"

// FIXME: probably more generic ways to do this (eg., for windows)
var GREP_BIN_PATH = "grep"

/**
 If array_a = ['a','b'] and array_b = ['c','d'], returns
 ['a','b','c','d']
*/
func merge_string_arrays (array_a,array_b []string) []string {
	var merged_array[] string
	for _, str := range array_a {
		merged_array = append(merged_array,str)
	}
	for _, str := range array_b {
		merged_array = append(merged_array,str)
	}
	return merged_array
}

/**
 Produces grep args from all sgrep files in subdirectories as well
 as the master .sgrep file in home.
*/
func grep_args_from_sgrep_files() [] string{
	// used to get the home directory of the user
	usr, err := user.Current()
	if err != nil {
		log.Fatal( err )
	}
	
	// sgrep_root_rule_list is a list of sgrep rules		
	sgrep_root_rule_list := ReadSgrep.ReadSgrepFile(
		usr.HomeDir)
	rule_tree := ReadSgrep.GetRuleTree(".")
	
	var grep_arg_array [] string
	for _, sgrep_rule := range sgrep_root_rule_list {
		grep_arg_array = append(
			grep_arg_array, sgrep_rule.Grep_arg_root_rule())
	}

	for _, grep_arg := range ReadSgrep.ProduceGrepArgs(rule_tree) {
		grep_arg_array = append(grep_arg_array,grep_arg)
	}

	return grep_arg_array
}


func main() {
	// read sgrep files
	grep_arg_array := grep_args_from_sgrep_files()
	// make the query recursive
	grep_arg_array = append(grep_arg_array,"-R")
	grep_arg_array = merge_string_arrays(os.Args[1:],grep_arg_array)

	// search in the current directory
	grep_arg_array = append(grep_arg_array, ".")

	cmd := exec.Command(GREP_BIN_PATH,grep_arg_array...)	
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	go io.Copy(os.Stdout, stdout) 
	go io.Copy(os.Stderr, stderr)
	cmd.Wait()
}