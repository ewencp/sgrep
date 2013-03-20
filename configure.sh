
#!/usr/bin/bash

root_proj=`pwd`
export GOPATH=$GOPATH:$root_proj

cd src/SgrepRules
go install

cd $root_proj
cd src/ReadSgrep
go install

cd $root_proj
cd src/sgrep
go install

