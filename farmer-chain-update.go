package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Farmer structure
type Farmer struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Product structure
type Product struct {
	ID       string `json:"id"`
	FarmerID string `json:"farmer_id"`
	Name     string `json:"name"`
	Price    string `json:"price"`
}

// SmartContract defines the chaincode
type SmartContract struct {
	contractapi.Contract
}

const (
	farmerPrefix  = "farmer-"
	productPrefix = "product-"
)

// RegisterFarmer adds a new farmer to the ledger
func (s *SmartContract) RegisterFarmer(ctx contractapi.TransactionContextInterface, id, name, email string) error {
	farmer := Farmer{ID: farmerPrefix + id, Name: name, Email: email}
	farmerJSON, err := json.Marshal(farmer)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(farmer.ID, farmerJSON)
}

// GetFarmer retrieves a farmer by ID
func (s *SmartContract) GetFarmer(ctx contractapi.TransactionContextInterface, id string) (*Farmer, error) {
	fullID := farmerPrefix + id

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

// RegisterProduct adds a new product to the ledger
func (s *SmartContract) RegisterProduct(ctx contractapi.TransactionContextInterface, id, farmerID, name, price string) error {
	product := Product{ID: productPrefix + id, FarmerID: farmerID, Name: name, Price: price}
	productJSON, err := json.Marshal(product)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(product.ID, productJSON)
}

// GetAllFarmers returns all farmers from the ledger
func (s *SmartContract) GetAllFarmers(ctx contractapi.TransactionContextInterface) ([]Farmer, error) {
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

// GetAllProducts returns all products from the ledger
func (s *SmartContract) GetAllProducts(ctx contractapi.TransactionContextInterface) ([]Product, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange(productPrefix, productPrefix+"~")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var products []Product
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var product Product
		err = json.Unmarshal(queryResponse.Value, &product)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
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
