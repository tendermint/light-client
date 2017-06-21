#!/bin/bash

oneTimeSetUp() {
  PASS=qwertyuiop
  export TM_HOME=$HOME/.lc_test
  tmcli reset_all
  assertTrue $?
}

# oneTimeTearDown() {
#   echo "tear down"
# }

newKey(){
  assertNotNull "keyname required" "$1"
  KEYPASS=${2:-qwertyuiop}
  echo $KEYPASS | tmcli keys new $1 >/dev/null 2>&1
  assertTrue "created $1" $?
}

testMakeKeys() {
  USER=demouser
  assertFalse "already user $USER" "tmcli keys get $USER"
  newKey $USER
  assertTrue "no user $USER" "tmcli keys get $USER"
}

# load and run these tests with shunit2!
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )" #get this files directory
. $DIR/shunit2
