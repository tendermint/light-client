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

  # this sets the base for all client queries in the tests
  export TMHOME=${BASE}/client
  tmcli init --node=tcp://localhost:46657 --genesis=${GENESIS_FILE} > /dev/null 2>&1
  if ! assertTrue "initialized light-client" "$?"; then
    return 1
  fi
}

oneTimeTearDown() {
  printf "\nstopping tendermint..."
  kill -9 $PID_SERVER >/dev/null 2>&1
  sleep 1
}

test01getInsecure() {
  GENESIS=$(tmcli rpc genesis)
  assertTrue "get genesis" "$?"
  MYCHAIN=$(echo ${GENESIS} | jq .genesis.chain_id | tr -d \")
  assertEquals "genesis chain matches" "${CHAIN_ID}" "${MYCHAIN}"

  STATUS=$(tmcli rpc status)
  assertTrue "get status" "$?"
  SHEIGHT=$(echo ${STATUS} | jq .latest_block_height)
  assertTrue "parsed status" "$?"
  assertNotNull "has a height" "${SHEIGHT}"

  VALS=$(tmcli rpc validators)
  assertTrue "get validators" "$?"
  VHEIGHT=$(echo ${VALS} | jq .block_height)
  assertTrue "parsed validators" "$?"
  assertTrue "sensible heights: $SHEIGHT / $VHEIGHT" "test $VHEIGHT -ge $SHEIGHT"
  VCNT=$(echo ${VALS} | jq '.validators | length')
  assertEquals "one validator" "1" "$VCNT"
}

# load and run these tests with shunit2!
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )" #get this files directory
. $DIR/shunit2
