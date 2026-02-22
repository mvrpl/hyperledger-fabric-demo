#!/bin/bash

set -e

export FABRIC_CFG_PATH="."

cd ./smart_contract_demo

go build -o smart_contract .

cd ..

peer lifecycle chaincode package hardwares.tar.gz --path ./smart_contract_demo/ --lang golang --label hardwares_1

export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_MSPCONFIGPATH=./organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

peer lifecycle chaincode install hardwares.tar.gz

peer lifecycle chaincode queryinstalled

peer lifecycle chaincode approveformyorg \
  --channelID system-channel \
  --name hardwares \
  --version 1.0 \
  --sequence 1 \
  --init-required \
  --package-id <PACKAGE_ID> \
  --orderer localhost:7050 \
  --tls --cafile $ORDERER_CA