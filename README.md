# Hyperledger Fabric Demo

This demo runs on docker, the first container is orderer and another 2 containers is network peers to simulate blockchain on VPN.

## Tools

### **Hyperledge Fabric**

Windows Powershell: `scoop install hyperledger-fabric` [Scoop Manifest](https://github.com/mvrpl/windows-apps/blob/main/bucket/hyperledger-fabric.json)
Unix (Mac OS or Linux): `brew install hyperledger-fabric` [Brew Formula](https://github.com/mvrpl/unix-apps/blob/main/Formula/hyperledger-fabric.rb)

## Generate all files of orderer, peer1 and peer2 

```bash
cryptogen generate --config ./crypto-config.yaml --output="organizations"
```

## Create root channel and generate genesis block

```bash
export FABRIC_CFG_PATH="."

configtxgen \
-profile TwoOrgsOrdererGenesis \
-outputBlock ./system-genesis-block/genesis.block \
-channelID system-channel
```

## Run docker containers

```bash
docker compose up -d
```

# Install Smart Contract

## Generate smart-contract package golang

```bash
export CORE_PEER_LOCALMSPID="SampleOrg"
export CORE_PEER_MSPCONFIGPATH="./organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp"
export CORE_PEER_ADDRESS="localhost:7051"

cd ./smart_contract_demo/

go mod tidy

cd -

peer lifecycle chaincode package hardwares.tar.gz --path ./smart_contract_demo/ --lang golang --label hardwares_1
```

## Deploy smart-contract to peer 1

```bash
docker pull hyperledger/fabric-ccenv

peer lifecycle chaincode install hardwares.tar.gz
```

## Create Channel

```bash
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=./organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt

configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel1.tx -channelID channel1
configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ./channel-artifacts/channel1.block -channelID channel1

osnadmin channel join \
--channelID channel1 \
--config-block ./channel-artifacts/channel1.block \
-o localhost:7053 \
--ca-file ./organizations/ordererOrganizations/example.com/orderers/localhost.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
--client-cert ./organizations/ordererOrganizations/example.com/users/Admin@example.com/tls/client.crt \
--client-key ./organizations/ordererOrganizations/example.com/users/Admin@example.com/tls/client.key
```

## Approve the Chaincode Definition

```bash
peer lifecycle chaincode approveformyorg \
--name hardwares \
--version 1.0 \
--package-id hardwares_1:50c9719b46f095c97afc7775e044861a41372b695dd2607121829bf7bdf4219b \
--sequence 1 \
--orderer localhost:7050 \
--tls \
--cafile organizations/ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem
```