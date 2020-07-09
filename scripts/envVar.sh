#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#

# This is a collection of bash functions used by different scripts

export CORE_PEER_TLS_ENABLED=true
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export PEER0_ORG1_CA=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt

# Set OrdererOrg.Admin globals
setOrdererGlobals() {
  export CORE_PEER_LOCALMSPID="OrdererMSP"
  export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
  export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/ordererOrganizations/example.com/users/Admin@example.com/msp
}

# Set environment variables for the peer org
setGlobals() {
  local USING_ORG=""
  if [ -z "$OVERRIDE_ORG" ]; then
    USING_ORG=$1
  else
    USING_ORG="${OVERRIDE_ORG}"
  fi
  PEER=$2
  echo "Using organization ${USING_ORG} peer ${PEER}"

    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG1_CA
    if [[ $PEER -eq 0 ]]; then 
      export CORE_PEER_ADDRESS=localhost:7051
    elif [[ $PEER -eq 1 ]]; then 
      export CORE_PEER_ADDRESS=localhost:8051
    elif [[ $PEER -eq 2 ]]; then 
      export CORE_PEER_ADDRESS=localhost:9051
    elif [[ $PEER -eq 3 ]]; then 
      export CORE_PEER_ADDRESS=localhost:10051
    elif [[ $PEER -eq 4 ]]; then 
      export CORE_PEER_ADDRESS=localhost:11051
    elif [[ $PEER -eq 5 ]]; then 
      export CORE_PEER_ADDRESS=localhost:7051
    elif [[ $PEER -eq 6 ]]; then 
      export CORE_PEER_ADDRESS=localhost:8051
    elif [[ $PEER -eq 7 ]]; then 
      export CORE_PEER_ADDRESS=localhost:9051
    elif [[ $PEER -eq 8 ]]; then 
      export CORE_PEER_ADDRESS=localhost:10051
    elif [[ $PEER -eq 9 ]]; then 
      export CORE_PEER_ADDRESS=localhost:11051
    elif [[ $PEER -eq 10 ]]; then 
      export CORE_PEER_ADDRESS=localhost:7051
    elif [[ $PEER -eq 11 ]]; then 
      export CORE_PEER_ADDRESS=localhost:8051
    elif [[ $PEER -eq 12 ]]; then 
      export CORE_PEER_ADDRESS=localhost:9051
    elif [[ $PEER -eq 13 ]]; then 
      export CORE_PEER_ADDRESS=localhost:10051
    elif [[ $PEER -eq 14 ]]; then 
      export CORE_PEER_ADDRESS=localhost:11051
    elif [[ $PEER -eq 15 ]]; then 
      export CORE_PEER_ADDRESS=localhost:7051
    elif [[ $PEER -eq 16 ]]; then 
      export CORE_PEER_ADDRESS=localhost:8051
    elif [[ $PEER -eq 17 ]]; then 
      export CORE_PEER_ADDRESS=localhost:9051
    elif [[ $PEER -eq 18 ]]; then 
      export CORE_PEER_ADDRESS=localhost:10051
    elif [[ $PEER -eq 19 ]]; then 
      export CORE_PEER_ADDRESS=localhost:11051
    fi
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp


  if [ "$VERBOSE" == "true" ]; then
    env | grep CORE
  fi
}

# parsePeerConnectionParameters $@
# Helper function that sets the peer connection parameters for a chaincode
# operation
parsePeerConnectionParameters() {

  PEER_CONN_PARMS=""
  PEERS=""
  while [ "$#" -gt 0 ]; do
    setGlobals $1
    PEER="peer0.org$1"
    ## Set peer adresses
    PEERS="$PEERS $PEER"
    PEER_CONN_PARMS="$PEER_CONN_PARMS --peerAddresses $CORE_PEER_ADDRESS"
    ## Set path to TLS certificate
    TLSINFO=$(eval echo "--tlsRootCertFiles \$PEER0_ORG$1_CA")
    PEER_CONN_PARMS="$PEER_CONN_PARMS $TLSINFO"
    # shift by one to get to the next organization
    shift
  done
  # remove leading space for output
  PEERS="$(echo -e "$PEERS" | sed -e 's/^[[:space:]]*//')"
}

verifyResult() {
  if [ $1 -ne 0 ]; then
    echo "!!!!!!!!!!!!!!! "$2" !!!!!!!!!!!!!!!!"
    echo
    exit 1
  fi
}