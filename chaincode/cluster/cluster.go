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
	Id       string `json:"id"`
	Suitable bool   `json:"suitable"`
}

type Offer struct {
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
	Owner    string `json:"owner"`
}

type Match struct {
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
	Seller   string `json:"seller"`
	Buyer    string `json:"buyer"`
}

// QueryResult structure used for handling result of query
type QueryResultUsers struct {
	Key    string `json:"Key"`
	Record *User
}
type QueryResultOffers struct {
	Key    string `json:"Key"`
	Record *Offer
}
type QueryResultMatches struct {
	Key    string `json:"Key"`
	Record *Match
}

// InitLedger adds a base set of users to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	users := []User{
		User{Id: "Miguel", Suitable: true},
		User{Id: "Arturo", Suitable: false},
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
func (s *SmartContract) CreateUser(ctx contractapi.TransactionContextInterface, userNumber string, id string, suitable bool) error {
	user := User{
		Id:       id,
		Suitable: suitable,
	}

	userAsBytes, _ := json.Marshal(user)

	return ctx.GetStub().PutState(userNumber, userAsBytes)
}

// CreateOffer adds a new offer to the world state with given details
func (s *SmartContract) CreateOffer(ctx contractapi.TransactionContextInterface, offerNumber string, quantity int, price int, owner string) error {
	offer := Offer{
		Quantity: quantity,
		Price:    price,
		Owner:    owner,
	}

	offerAsBytes, _ := json.Marshal(offer)

	return ctx.GetStub().PutState(offerNumber, offerAsBytes)
}

func (s *SmartContract) MatchOffers(ctx contractapi.TransactionContextInterface) error {
	startKey := "OFFER0"
	endKey := "OFFER99"
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return err
	}
	defer resultsIterator.Close()

	buyers := []Offer{}
	sellers := []Offer{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return err
		}
		offer := new(Offer)
		err = json.Unmarshal(queryResponse.Value, offer)
		if err != nil {
			return err
		}
		if offer.Quantity > 0 {
			buyers = append(buyers, *offer)
		} else {
			sellers = append(sellers, *offer)
		}
	}

	i, j, num := 0, 0, 3

	for i < len(sellers) && j < len(buyers) {
		match := new(Match)
		match.Buyer = buyers[j].Owner
		match.Seller = sellers[i].Owner

		if sellers[i].Quantity*-1 > buyers[j].Quantity {
			match.Quantity = buyers[j].Quantity
		} else {
			match.Quantity = sellers[i].Quantity * -1
		}

		match.Price = j
		matchAsBytes, err := json.Marshal(match)
		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
		err = ctx.GetStub().PutState("MATCH"+strconv.Itoa(num), matchAsBytes)
		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}

		sellers[i].Quantity = sellers[i].Quantity + match.Quantity
		buyers[j].Quantity = buyers[j].Quantity - match.Quantity

		if sellers[i].Quantity == 0 {
			i++
		}
		if buyers[j].Quantity == 0 {
			j++
		}

		num++
	}

	return nil
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
func (s *SmartContract) QueryAllUsers(ctx contractapi.TransactionContextInterface) ([]QueryResultUsers, error) {
	startKey := "USER0"
	endKey := "USER99"

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResultUsers{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		user := new(User)
		_ = json.Unmarshal(queryResponse.Value, user)

		QueryResultUsers := QueryResultUsers{Key: queryResponse.Key, Record: user}
		results = append(results, QueryResultUsers)
	}

	return results, nil
}
func (s *SmartContract) QueryAllOffers(ctx contractapi.TransactionContextInterface) ([]QueryResultOffers, error) {
	startKey := "OFFER0"
	endKey := "OFFER99"

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResultOffers{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		offer := new(Offer)
		_ = json.Unmarshal(queryResponse.Value, offer)

		QueryResultOffers := QueryResultOffers{Key: queryResponse.Key, Record: offer}
		results = append(results, QueryResultOffers)
	}

	return results, nil
}
func (s *SmartContract) QueryAllMatches(ctx contractapi.TransactionContextInterface) ([]QueryResultMatches, error) {
	startKey := "MATCH0"
	endKey := "MATCH99"

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResultMatches{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		match := new(Match)
		_ = json.Unmarshal(queryResponse.Value, match)

		QueryResultMatches := QueryResultMatches{Key: queryResponse.Key, Record: match}
		results = append(results, QueryResultMatches)
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
