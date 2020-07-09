/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package main

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
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Address    string `json:"address"`
	Generation string `json:"generation"`
	Coins      int    `json:"coins"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *User
}

// InitLedger adds a base set of users to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	users := []User{
		User{Name: "Miguel", Surname: "Gayo", Address: "C/ inventada", Generation: "PV", Coins: 1000},
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

// CreateUser adds a new user to the world state with given details
func (s *SmartContract) CreateUser(ctx contractapi.TransactionContextInterface, userNumber string, name string, surname string, address string, generation string, coins int) error {
	user := User{
		Name:       name,
		Surname:    surname,
		Address:    address,
		Generation: generation,
		Coins:      coins,
	}

	userAsBytes, _ := json.Marshal(user)

	return ctx.GetStub().PutState(userNumber, userAsBytes)
}

// CreateUser adds a new user to the world state with given details
func (s *SmartContract) AddCoins(ctx contractapi.TransactionContextInterface, userNumber string, coins int) error {

	userAsBytes, _ := ctx.GetStub().GetState(userNumber)

	user := new(User)
	_ = json.Unmarshal(userAsBytes, user)
	user.Coins = user.Coins*1 + coins

	userAsBytes, _ = json.Marshal(user)

	return ctx.GetStub().PutState(userNumber, userAsBytes)
}

// Query returns the user stored in the world state with given id
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
