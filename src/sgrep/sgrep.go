package main

import "ReadSgrep"
import "os"
import "os/exec"
import "io"
import "log"


// FIXME: probably more generic ways to do this (eg., for windows)
var GREP_BIN_PATH = "grep"
var FIND_BIN_PATH = "find"


/**
 Produces grep args from all sgrep files in subdirectories as well
 as the master .sgrep file in home.
*/
func find_args_from_sgrep_files() [] string{
	
	rule_tree := ReadSgrep.GetRuleTree(".")
	var grep_arg_array [] string

	for _, grep_arg := range ReadSgrep.ProduceGrepArgs(rule_tree) {
		if grep_arg != "" {
			grep_arg_array = append(grep_arg_array,grep_arg)
		}
	}

	return grep_arg_array
}


func main() {
	// read sgrep files
	exclude_rule_list := find_args_from_sgrep_files()
	var arg_array [] string

	if len(exclude_rule_list) == 0 {
		// FIXME: wildcard is not os agnostic
		arg_array = append(arg_array,".","-path","*")
	}else {
		arg_array = append(arg_array,".")
		for index, element := range exclude_rule_list {
			arg_array = append(arg_array,"-path")
			arg_array = append(arg_array,"./" + element)
			if index != len(exclude_rule_list) -1 {
				arg_array = append(arg_array, "-o")
			}
		}
		
		arg_array = append(arg_array,"-prune","-o")
	}
	// tell to execute grep 
	arg_array = append(arg_array,"-exec", GREP_BIN_PATH)
	// ... searching with args passed into the command line
	arg_array = append(arg_array,os.Args[1:]...)
	// pass -H to grep (telling it to display filename and
	// matching line), {} for taking arguments from find, and ';'
	// to end exec statement.
	arg_array = append(arg_array,"-H","{}",";")
	
	cmd := exec.Command(FIND_BIN_PATH,arg_array...)		
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