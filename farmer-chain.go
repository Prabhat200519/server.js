package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Farmer structure
type Farmer struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

// Consumer structure
type Consumer struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

// Transaction structure
type Transaction struct {
	ID         string `json:"id"`
	FarmerID   string `json:"farmer_id"`
	ConsumerID string `json:"consumer_id"`
	Amount     string `json:"amount"`
	Timestamp  string `json:"timestamp"`
}

// SmartContract defines the chaincode
type SmartContract struct {
	contractapi.Contract
}

const (
	farmerPrefix      = "farmer-"
	consumerPrefix    = "consumer-"
	transactionPrefix = "transaction-"
)

// RegisterFarmer adds a new farmer to the ledger
func (s *SmartContract) RegisterFarmer(ctx contractapi.TransactionContextInterface, id, name, location string) error {
	farmer := Farmer{ID: farmerPrefix + id, Name: name, Location: location}
	farmerJSON, err := json.Marshal(farmer)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(farmer.ID, farmerJSON)
}

//Manual edit by myself
func (s *SmartContract) GetFarmer(ctx contractapi.TransactionContextInterface, id string) (*Farmer, error) {
    fullID := "farmer-" + id  // Ensure consistent prefix

    farmerJSON, err := ctx.GetStub().GetState(fullID)
    if err != nil {
        return nil, fmt.Errorf("failed to read from world state: %v", err)
    }
    if farmerJSON == nil {
        return nil, fmt.Errorf("farmer %s does not exist", fullID)
    }

    var farmer Farmer
    err = json.Unmarshal(farmerJSON, &farmer)
    if err != nil {
        return nil, err
    }

    return &farmer, nil
}


// RegisterConsumer adds a new consumer to the ledger
func (s *SmartContract) RegisterConsumer(ctx contractapi.TransactionContextInterface, id, name, location string) error {
	consumer := Consumer{ID: consumerPrefix + id, Name: name, Location: location}
	consumerJSON, err := json.Marshal(consumer)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(consumer.ID, consumerJSON)
}

// RecordTransaction stores a new transaction on the ledger
func (s *SmartContract) RecordTransaction(ctx contractapi.TransactionContextInterface, id, farmerID, consumerID, amount, timestamp string) error {
	transaction := Transaction{ID: transactionPrefix + id, FarmerID: farmerID, ConsumerID: consumerID, Amount: amount, Timestamp: timestamp}
	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(transaction.ID, transactionJSON)
}

// QueryAllFarmers returns all farmers from the ledger
func (s *SmartContract) QueryAllFarmers(ctx contractapi.TransactionContextInterface) ([]Farmer, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange(farmerPrefix, farmerPrefix+"~")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var farmers []Farmer
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var farmer Farmer
		err = json.Unmarshal(queryResponse.Value, &farmer)
		if err != nil {
			return nil, err
		}
		farmers = append(farmers, farmer)
	}
	return farmers, nil
}

// QueryAllConsumers returns all consumers from the ledger
func (s *SmartContract) QueryAllConsumers(ctx contractapi.TransactionContextInterface) ([]Consumer, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange(consumerPrefix, consumerPrefix+"~")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var consumers []Consumer
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var consumer Consumer
		err = json.Unmarshal(queryResponse.Value, &consumer)
		if err != nil {
			return nil, err
		}
		consumers = append(consumers, consumer)
	}
	return consumers, nil
}

// QueryAllTransactions returns all transactions from the ledger
func (s *SmartContract) QueryAllTransactions(ctx contractapi.TransactionContextInterface) ([]Transaction, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange(transactionPrefix, transactionPrefix+"~")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var transactions []Transaction
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var transaction Transaction
		err = json.Unmarshal(queryResponse.Value, &transaction)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating chaincode: %s", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
