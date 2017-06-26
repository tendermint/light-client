#!/bin/bash

oneTimeSetUp() {
  BASE=~/.lc_init_test
  rm -rf "$BASE"
  mkdir -p "$BASE"

  SERVER="${BASE}/server"
  SERVER_LOG="${BASE}/tendermint.log"

  tendermint init --home="$SERVER" >> "$SERVER_LOG"
  if ! assertTrue $?; then return 1; fi

  GENESIS_FILE=${SERVER}/genesis.json
  CHAIN_ID=$(cat ${GENESIS_FILE} | jq .chain_id | tr -d \")

  printf "starting tendermint...\n"
  tendermint node --proxy_app=dummy --home="$SERVER" >> "$SERVER_LOG" 2>&1 &
  sleep 5
  PID_SERVER=$!
  disown
  if ! ps $PID_SERVER >/dev/null; then
    echo "**STARTUP FAILED**"
    cat $SERVER_LOG
    return 1
  fi
}

oneTimeTearDown() {
  printf "\nstopping tendermint..."
  kill -9 $PID_SERVER >/dev/null 2>&1
  sleep 1
}

test01goodInit() {
  export TMHOME=${BASE}/client-01
  assertFalse "ls ${TMHOME} 2>/dev/null >&2"

  echo y | tmcli init --node=tcp://localhost:46657 --chain-id="${CHAIN_ID}" > /dev/null
  assertTrue "initialized light-client" $?
  checkDir $TMHOME 3
}

test02badInit() {
  export TMHOME=${BASE}/client-02
  assertFalse "ls ${TMHOME} 2>/dev/null >&2"

  # no node where we go
  echo y | tmcli init --node=tcp://localhost:9999 --chain-id="${CHAIN_ID}" > /dev/null 2>&1
  assertFalse "invalid init" $?
  # dir there, but empty...
  checkDir $TMHOME 0

  # try with invalid chain id
  echo y | tmcli init --node=tcp://localhost:46657 --chain-id="bad-chain-id" > /dev/null 2>&1
  assertFalse "invalid init" $?
  checkDir $TMHOME 0

  # reject the response
  echo n | tmcli init --node=tcp://localhost:46657 --chain-id="${CHAIN_ID}" > /dev/null 2>&1
  assertFalse "invalid init" $?
  checkDir $TMHOME 0
}

test03noDoubleInit() {
  export TMHOME=${BASE}/client-03
  assertFalse "ls ${TMHOME} 2>/dev/null >&2"

  # init properly
  echo y | tmcli init --node=tcp://localhost:46657 --chain-id="${CHAIN_ID}" > /dev/null 2>&1
  assertTrue "initialized light-client" $?
  checkDir $TMHOME 3

  # try again, and we get an error
  echo y | tmcli init --node=tcp://localhost:46657 --chain-id="${CHAIN_ID}" > /dev/null 2>&1
  assertFalse "warning on re-init" $?
  checkDir $TMHOME 3

  # unless we --force-reset
  echo y | tmcli init --force-reset --node=tcp://localhost:46657 --chain-id="${CHAIN_ID}" > /dev/null 2>&1
  assertTrue "re-initialized light-client" $?
  checkDir $TMHOME 3
}

test04acceptGenesisFile() {
  export TMHOME=${BASE}/client-04
  assertFalse "ls ${TMHOME} 2>/dev/null >&2"

  # init properly
  tmcli init --node=tcp://localhost:46657 --genesis=${GENESIS_FILE} > /dev/null 2>&1
  assertTrue "initialized light-client" $?
  checkDir $TMHOME 3
}

# XXX Ex: checkDir $DIR $FILES
# Makes sure directory exists and has the given number of files
checkDir() {
  assertTrue "ls ${1} 2>/dev/null >&2"
  assertEquals "no files created" "$2" $(ls $1 | wc -l)
}

# load and run these tests with shunit2!
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )" #get this files directory
. $DIR/shunit2
