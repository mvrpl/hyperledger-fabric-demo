package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Hardware struct {
	VendorID       string  `json:"VendorID"`
	ModelID        string  `json:"ModelID"`
	SerialNumber   string  `json:"SerialNumber"`
	OwnerName      string  `json:"OwnerName"`
	ConditionScore float32 `json:"ConditionScore"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	hardwares := []Hardware{
		{VendorID: "0x10DE", ModelID: "0x2B85", SerialNumber: "24230073RB09A", OwnerName: "Yang Son", ConditionScore: 98.8},
		{VendorID: "Apple", ModelID: "A3084", SerialNumber: "EID_9832749823740923470923400070927049092484", OwnerName: "Jin Soo", ConditionScore: 79.3},
		{VendorID: "Ubiquiti Inc.", ModelID: "UDR7-US", SerialNumber: "245A4C86D9B1", OwnerName: "Kabinov Vassili Rostislavovich", ConditionScore: 99.6},
	}

	for _, hardware := range hardwares {
		hardwareJSON, err := json.Marshal(hardware)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(hardware.SerialNumber, hardwareJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}
	fmt.Println("InitLedger Done!")
	return nil
}

func (s *SmartContract) CreateHardware(ctx contractapi.TransactionContextInterface, hardware Hardware) error {
	id := hardware.SerialNumber
	exists, err := s.HardwareExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the hardware %s already exists", id)
	}

	hardwareJSON, err := json.Marshal(hardware)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, hardwareJSON)
}

func (s *SmartContract) ReadHardware(ctx contractapi.TransactionContextInterface, id string) (*Hardware, error) {
	hardwareJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if hardwareJSON == nil {
		return nil, fmt.Errorf("the hardware %s does not exist", id)
	}

	var hardware Hardware
	err = json.Unmarshal(hardwareJSON, &hardware)
	if err != nil {
		return nil, err
	}

	return &hardware, nil
}

func (s *SmartContract) UpdateHardware(ctx contractapi.TransactionContextInterface, id string, conditionScore float32) error {
	exists, err := s.HardwareExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the hardware %s does not exist", id)
	}

	hardware, err := s.ReadHardware(ctx, id)
	if err != nil {
		return err
	}

	hardware.ConditionScore = conditionScore

	hardwareJSON, err := json.Marshal(hardware)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, hardwareJSON)
}

func (s *SmartContract) DeleteHardware(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.HardwareExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the hardware %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

func (s *SmartContract) HardwareExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	hardwareJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return hardwareJSON != nil, nil
}

func (s *SmartContract) TransferHardware(ctx contractapi.TransactionContextInterface, id string, newOwner string) (string, error) {
	hardware, err := s.ReadHardware(ctx, id)
	if err != nil {
		return "", err
	}

	oldOwner := hardware.OwnerName
	hardware.OwnerName = newOwner

	hardwareJSON, err := json.Marshal(hardware)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(id, hardwareJSON)
	if err != nil {
		return "", err
	}

	return oldOwner, nil
}

func (s *SmartContract) GetAllHardwares(ctx contractapi.TransactionContextInterface) ([]*Hardware, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var hardwares []*Hardware
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var hardware Hardware
		err = json.Unmarshal(queryResponse.Value, &hardware)
		if err != nil {
			return nil, err
		}
		hardwares = append(hardwares, &hardware)
	}

	return hardwares, nil
}

func main() {
	simpleContract := new(SmartContract)

	cc, err := contractapi.NewChaincode(simpleContract)
	if err != nil {
		panic(err.Error())
	}
	if err := cc.Start(); err != nil {
		panic(err.Error())
	}
}
