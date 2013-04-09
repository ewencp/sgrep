sgrep
==================

Overview
------------------
sgrep is a simple extension of grep.  Like grep, using sgrep, you can
look for regular expressions in files on your file system.  Unlike
grep, users can specify .sgrep files in their folders.  These provide
custom instructions to sgrep to *not* look in particular files when
trying to match a regular expression.  For instance, a user that did
not want to grep through a subfolder containing many binary files or
external libraries could add a .sgrep file with a rule to ignore it.

.sgrep files
-------------------
When executing sgrep from a directory, PWD_DIR, sgrep checks for any
.sgrep files in any parent of PWD_DIR or in any subdirectory of
PWD_DIR.

sgrep reads the rules of each .sgrep file.  Rules specify which files
to exclude from searching through.  For instance, to avoid sgrep-ing
through a file, "hello.txt", in PWD_DIR, one could write a .sgrep file
in PWD_DIR with the contents:

hello.txt

To avoid searching somefile sub_dir_hello.txt in PWD_DIR/SUB_DIR, one
could either add a .sgrep file in PWD_DIR/SUB_DIR with the contents

hello.txt

or could add a rule to his/her .sgrep file in PWD_DIR that says:

SUB_DIR/hello.txt

With either of these rules, sgrep would still search through any file
named "hello.txt" that was not directly in SUB_DIR, eg.,
PWD_DIR/hello.txt, PWD_DIR/SUB_DIR/SUB_SUB_DIR/hello.txt, or
PWD_DIR/OTHER_SUB_DIR/hello.txt.

sgrep also matches wildcards.  For instance, to prevent matching any 
files that end in "pyc", use the rule

*pyc

sgrep also includes rules derived from parent directories.  As an
example, for a .sgrep file in PARENT_DIR (where PARENT_DIR is the
parent directory of PWD_DIR) containing:

PWD_DIR/hello.txt

running sgrep from PWD_DIR will skip PWD_DIR/hello.txt (but will not
skip PWD_DIR/SUB_DIR/hello.txt).  The rules for wildcard matching from
parent directories are a little involved.  For a rule specified in a
distant parent directory like the following

/a/*/c/*/PWD_DIR/*/hello.txt

in a file system organized as follows:

/a/b/c/PARENT_DIR/PWD_DIR/

sgrep, will take the most specific part of the rule that it matches
and use that as the rule.  Ie, sgrep will apply the rule:

*/hello.txt

Instead of the rule:

*/c/PWD_DIR/*/hello.txt

You can add comments to sgrep files by prefixing the line with '#'.


Arguments
------------------
sgrep takes in all the same arguments as grep, with two minor
differences:

  * By default, sgrep performs a recursive search, and so does not
    need the -R parameter

  * By default, sgrep searches in the directory from which it is
    executed.


Examples
------------------
sgrep "sgrep is great"  --- Searches for any instances of the phrase
"sgrep is great" in subdirectories, subject to .sgrep file rules.
