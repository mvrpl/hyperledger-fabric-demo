package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"github.com/mvrpl/hyperledger-fabric-smart-contract-demo/mocks"
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

func (h *Hardware) Set(s string) error {
	err := json.Unmarshal([]byte(s), h)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}
	return nil
}

func (h *Hardware) String() string {
	return fmt.Sprintf("VendorID: %s, ModelID: %s, SerialNumber: %s, OwnerName: %s, ConditionScore: %.1f", h.VendorID, h.ModelID, h.SerialNumber, h.OwnerName, h.ConditionScore)
}

func main() {
	sc := SmartContract{}

	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	hardware := Hardware{}

	funcName := flag.String("function", "InitLedger", "A function name of smart contract to execute")
	newOwner := flag.String("newOwner", "", "The new owner of the hardware to transfer")
	hardwareID := flag.String("hardwareID", "", "The ID of the hardware")
	conditionScore := flag.Float64("conditionScore", 0.0, "The condition score of the hardware")
	flag.Var(&hardware, "hardware", "A JSON string of the hardware object (e.g., '{\"VendorID\":\"ManufactureName\",\"ModelID\":\"ModelName\",\"SerialNumber\":\"SerialNumber\",\"OwnerName\":\"OwnerName\",\"ConditionScore\":98.8}')")

	flag.Parse()

	switch *funcName {
	case "InitLedger":
		err := sc.InitLedger(transactionContext)
		if err != nil {
			fmt.Println("Error initializing ledger:", err)
			os.Exit(166)
		}
	case "GetAllHardwares":
		hardwares, err := sc.GetAllHardwares(transactionContext)
		if err != nil {
			fmt.Println("Error getting all hardwares:", err)
			os.Exit(167)
		} else {
			fmt.Println("All Hardwares:", hardwares)
		}
	case "GetHardware":
		hardware, err := sc.ReadHardware(transactionContext, *hardwareID)
		if err != nil {
			fmt.Println("Error getting hardware:", err)
			os.Exit(167)
		} else {
			fmt.Println("Hardware:", hardware)
		}
	case "TransferHardware":
		oldOwner, err := sc.TransferHardware(transactionContext, *hardwareID, *newOwner)
		if err != nil {
			fmt.Println("Error transferring hardware:", err)
			os.Exit(168)
		} else {
			fmt.Println(oldOwner, " => ", *newOwner)
		}
	case "DeleteHardware":
		err := sc.DeleteHardware(transactionContext, *hardwareID)
		if err != nil {
			fmt.Println("Error deleting hardware:", err)
			os.Exit(169)
		} else {
			fmt.Println("Hardware deleted:", *hardwareID)
		}
	case "UpdateHardware":
		err := sc.UpdateHardware(transactionContext, *hardwareID, float32(*conditionScore))
		if err != nil {
			fmt.Println("Error updating hardware:", err)
			os.Exit(170)
		} else {
			fmt.Println("Hardware updated:", *hardwareID)
		}
	case "CreateHardware":
		err := sc.CreateHardware(transactionContext, hardware)
		if err != nil {
			fmt.Println("Error creating hardware:", err)
			os.Exit(171)
		} else {
			fmt.Println("Hardware created:", hardware.SerialNumber)
		}
	default:
		fmt.Println("Invalid function name. Use -function flag to specify a valid function (e.g., InitLedger, GetAllHardwares, GetHardware, TransferHardware, DeleteHardware).")
		os.Exit(199)
	}

}
