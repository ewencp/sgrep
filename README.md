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

Usage
------------------
man grep

Examples
------------------
sgrep "sgrep is great" 
