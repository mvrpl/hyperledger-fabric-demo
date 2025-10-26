# Hyperledger Fabric Demo

This demo runs on docker, the first container is orderer and another 2 containers is network peers to simulate blockchain on VPN.

## Tools

### **Hyperledge Fabric**

Windows Powershell: `scoop install hyperledger-fabric`

Unix (Mac OS or Linux): `brew install hyperledger-fabric`

### **Hyperledger Fabric CA**

Windows Powershell: `scoop install hyperledger-fabric-ca`

Unix (Mac OS or Linux): `brew install hyperledger-fabric-ca`

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
