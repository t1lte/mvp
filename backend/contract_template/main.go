package main

import (
    "log"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
    "chaincode"
)

func main() {
    chaincodeSmartContract, err := contractapi.NewChaincode(&chaincode.SmartContract{})
    if err != nil {
        log.Panicf("Error creating chaincode: %v", err)
    }
    if err := chaincodeSmartContract.Start(); err != nil {
        log.Panicf("Error starting chaincode: %v", err)
    }
}