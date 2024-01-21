package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
        "strconv"
        "log"
        "time"
        "github.com/golang/protobuf/ptypes"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// HistoryQueryResult structure used for returning result of history query
type HistoryQueryResult struct {
	Record    *Asset    `json:"record"`
	TxId      string    `json:"txId"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"isDelete"`
}

// Asset describes basic details of what makes up a simple asset
// Insert struct field in alphabetic order => to achieve determinism across languages
// golang keeps the order when marshal to json but doesn't order automatically
type Asset struct {
	DEALERID         string  `json:"DEALERID"`
	MSISDN           int     `json:"MSISDN"`
	MPIN             int     `json:"MPIN"`
	BALANCE          float64 `json:"BALANCE"`
	STATUS           string  `json:"STATUS"`
        TRANSAMOUNT      float64 `json:"TRANSAMOUNT"`
        TRANSTYPE        string  `json:"TRANSTYPE"`
        REMARKS          string  `json:"REMARKS"`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{
		{DEALERID: "asset1", MSISDN: 9876, MPIN: 5121, BALANCE: 10000.0, STATUS: "ACTIVE", TRANSAMOUNT: 0.0, TRANSTYPE:"", REMARKS: "Initital Deposit"},
		{DEALERID: "asset2", MSISDN: 5432, MPIN: 1254, BALANCE: 25000.0, STATUS: "ACTIVE", TRANSAMOUNT: 0.0, TRANSTYPE:"", REMARKS: "Initital Deposit"},
		{DEALERID: "asset3", MSISDN: 1098, MPIN: 1100, BALANCE: 50000.0, STATUS: "ACTIVE", TRANSAMOUNT: 0.0, TRANSTYPE:"", 
REMARKS: "Initital Deposit"},

	}

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.DEALERID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, dealerid string, msisdn int, mpin int, balance float64, status string, transamount float64, transtype string, remarks string) error {
	exists, err := s.AssetExists(ctx, dealerid)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", dealerid)
	}

	asset := Asset{
		DEALERID:            dealerid,
		MSISDN:                msisdn,
		MPIN:                    mpin,
		BALANCE:              balance,
		STATUS:                status,
                TRANSAMOUNT:      transamount,
                TRANSTYPE:          transtype,
                REMARKS:              remarks,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(dealerid, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, dealerid string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(dealerid)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", dealerid)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, dealerid string, msisdn int, mpin int, balance float64, status string, transamount float64, transtype string, remarks string) error {
	exists, err := s.AssetExists(ctx, dealerid)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", dealerid)
	}

	// overwriting original asset with new asset
	asset := Asset{
		DEALERID:            dealerid,
		MSISDN:                msisdn,
		MPIN:                    mpin,
		BALANCE:              balance,
		STATUS:                status,
                TRANSAMOUNT:      transamount,
                TRANSTYPE:          transtype,
                REMARKS:              remarks,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(dealerid, assetJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, dealerid string) error {
	exists, err := s.AssetExists(ctx, dealerid)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", dealerid)
	}

	return ctx.GetStub().DelState(dealerid)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, dealerid string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(dealerid)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// TransferAsset updates the owner field of asset with given id in world state, and returns the old owner.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, dealerid string, newOwner int) (string, error) {
	asset, err := s.ReadAsset(ctx, dealerid)
	if err != nil {
		return "", err
	}

	oldOwner := asset.MPIN
	asset.MPIN = newOwner

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(dealerid, assetJSON)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(oldOwner), nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

// GetAssetHistory returns the chain of custody for an asset since issuance.
func (s *SmartContract) GetAssetHistory(ctx contractapi.TransactionContextInterface, dealerid string) ([]HistoryQueryResult, error) {
	log.Printf("GetAssetHistory: ID %v", dealerid)

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(dealerid)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &asset)
			if err != nil {
				return nil, err
			}
		} else {
			asset = Asset{
				DEALERID: dealerid,
			}
		}

		timestamp, err := ptypes.Timestamp(response.Timestamp)
		if err != nil {
			return nil, err
		}

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: timestamp,
			Record:    &asset,
			IsDelete:  response.IsDelete,
		}
		records = append(records, record)
	}

	return records, nil
}
