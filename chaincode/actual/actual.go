package main

/*

#cgo LDFLAGS: -L /usr/lib -liec61850 -Wl,-rpath=/usr/lib
#include "iec61850_model.h"
#include "iec61850_client.h"

#include <stdlib.h>
#include <stdio.h>

#include "hal_thread.h"



float read() {
	char* hostname ="192.168.0.2";
	int tcpPort = 102;
	IedClientError error;
	IedConnection con = IedConnection_create();
	IedConnection_connect(con, &error, hostname, tcpPort);
	if (error == IED_ERROR_OK){
		MmsValue* vol_battery = IedConnection_readObject(con, &error, "RPBBESS/ZBAT1.Vol.mag.f", IEC61850_FC_MX);
		if (error == IED_ERROR_OK){
			float vol_bat = MmsValue_toFloat(vol_battery);
			IedConnection_destroy(con);
			return vol_bat;
		}
	}
	IedConnection_destroy(con);
}


*/
import "C"

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a user
type SmartContract struct {
	contractapi.Contract
}

// User describes basic details of what makes up a user
type User struct {
	Battery     float32 `json:"battery"`
	Generation  float32 `json:"generation"`
	Consumption float32 `json:"consumption"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *User
}

// InitLedger adds a base set of users to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	battery := C.read()

	users := []User{
		User{Battery: float32(battery), Generation: 0, Consumption: 0},
	}

	for i, user := range users {
		userAsBytes, _ := json.Marshal(user)
		err := ctx.GetStub().PutState("USER"+strconv.Itoa(i), userAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

func (s *SmartContract) ReadBattery(ctx contractapi.TransactionContextInterface, userNumber string) error {

	userAsBytes, _ := ctx.GetStub().GetState(userNumber)
	battery := C.read()
	user := new(User)
	_ = json.Unmarshal(userAsBytes, user)
	user.Battery = float32(battery)
	userAsBytes, _ = json.Marshal(user)

	return ctx.GetStub().PutState(userNumber, userAsBytes)
}

func (s *SmartContract) QueryUser(ctx contractapi.TransactionContextInterface, userNumber string) (*User, error) {
	userAsBytes, err := ctx.GetStub().GetState(userNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if userAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", userNumber)
	}

	user := new(User)
	_ = json.Unmarshal(userAsBytes, user)

	return user, nil
}

// QueryAllUsers returns all users found in world state
func (s *SmartContract) QueryAllUsers(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := "USER0"
	endKey := "USER99"

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		user := new(User)
		_ = json.Unmarshal(queryResponse.Value, user)

		queryResult := QueryResult{Key: queryResponse.Key, Record: user}
		results = append(results, queryResult)
	}

	return results, nil
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create personal chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting personal chaincode: %s", err.Error())
	}
}
