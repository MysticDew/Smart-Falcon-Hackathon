# Smart-Falcon-Hackathon

**INTRODUCTION**

Hyperledger Fabric is a blockchain framework within the Hyperledger project that aims to provide a modular and scalable solution for developing enterprise-level blockchain applications. It is specifically designed for use in permissioned networks, where participants are known entities and have a reason to trust the network. Key features of Hyperledger Fabric:

Permissioned Network
Modular Architecture
Smart Contracts (Chaincode)
Privacy and Confidentiality
Endorsement Policies
Scalability
Pluggable Consensus Mechanisms
It is built by primarily using

Ubuntu 20.04.6: Ubuntu 20.04.6 LTS is a stable and long-term support release of the Ubuntu operating system, providing users with a reliable and secure foundation for their computing needs.
Docker Desktop: Docker Desktop is a powerful platform for developing, shipping, and running applications in containers. It simplifies the deployment and management of containerized applications, allowing developers to focus on writing code without worrying about the intricacies of the underlying infrastructure.
Golang: Go, commonly known as Golang, is a statically typed, compiled programming language designed for simplicity, efficiency, and ease of use. It was created by Google engineers Robert Griesemer, Rob Pike, and Ken Thompson and first released in 2009.

**PREREQUISTES**

The below two lines contain the code to export the downloaded go package to the home as well as bin directories.

 export GOPATH=$HOME/go
 export PATH=$PATH:$GOPATH/bin

 The below code is used to export the bin and config files to the current directory
 export PATH=${PWD}/../bin:$PATH
 export FABRIC_CFG_PATH=$PWD/../config/

To set the environment variables that allow you to operate the peer CLI as Org1:

 export CORE_PEER_TLS_ENABLED=true
 export CORE_PEER_LOCALMSPID="Org1MSP"
 export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
 export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
 export CORE_PEER_ADDRESS=localhost:7051

To set the environment variables that allow you to operate the peer CLI as Org2:

 export CORE_PEER_TLS_ENABLED=true
 export CORE_PEER_LOCALMSPID="Org2MSP"
 export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
 export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
 export CORE_PEER_ADDRESS=localhost:9051


import the following packages:

 "github.com/hyperledger/fabric-contract-api-go/contractapi"
 "github.com/golang/protobuf/ptypes"
 "strconv"
 "log"
 "time"

 **PROCEDURE**

 # step 1: To setup Hyperledger Fabric Test Network
use the commands:-
 ./network.sh down
 ./network.sh up

# step 2: Package and Deploy the chain code into fabric test network 
use the command:-
 ./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go

# step 3: Develop Hyperledger Fabric Chaincode
use the command:-
 peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c '{"function":"InitLedger","Args":[]}'

 # step 4: Test the chain code functionality using Fabric Peer CLI commands
 # 4.1: To create Asset: use the below command
   peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c '{"function":"CreateAsset","Args":["asset4","4590","2002","750000.0","active",
"0.0","","initial deposit"]}'

# 4.2: To update asset: use the below command
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c '{"function":"UpdateAsset","Args":["asset1","1412","1230","99000
0.0","active","0.0","","initial deposit"]}'

# 4.3: To read asset: use the below command
  peer chaincode query -C mychannel -n basic -c '{"Args":["ReadAsset","asset3"]}'

# 4.4: To get Asset Transaction History: use the below command
  peer chaincode query -C mychannel -n basic -c '{"Args":["GetAssetHistory","Asset4"]}'

 
