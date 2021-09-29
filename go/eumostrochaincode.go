package main

import (
  "encoding/json"
  "fmt"
  "log"
  "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
  contractapi.Contract
}

type LotProduct struct {
  name       string  `json:"name"`
  quantity   int     `json:"quantity"`
  unitValue  int     `json:"unitValue"`
  unit       string  `json:"unit"`
}

type Lot struct {
  ID                string        `json:"ID"`
  nfId              string        `json:"nfId"`
  lotProducts       []LotProduct  `json:"lotProducts"`
  Owner             string        `json:"owner"`
  OwnerId           string        `json:"ownerId"`
  lotType           string        `json:"lotType"`
  createdAt         string        `json:"createdAt"`
  total             int           `json:"total"`
  formatedAddress   string        `json:"formatedAddress"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
  assets := []Lot{}

  for _, asset := range assets {
    assetJSON, err := json.Marshal(asset)
    if err != nil {
      return err
    }

    err = ctx.GetStub().PutState(asset.ID, assetJSON)
    if err != nil {
      return fmt.Errorf("failed to put to world state. %v", err)
    }
  }

  return nil
}

func (s *SmartContract) CreateAsset(
  ctx contractapi.TransactionContextInterface,
  ID string,
  nfId string,
  lotProducts []LotProduct,
  Owner string,
  OwnerId int,
  lotType string,
  createdAt string,
  total int,
  formatedAddress string,
) error {
  exists, err := s.AssetExists(ctx, id)
  if err != nil {
    return err
  }
  if exists {
    return fmt.Errorf("the asset %s already exists", id)
  }
  asset := Lot{
    ID:               ID,
    nfId:             nfId,
    lotProducts:      lotProducts,
    Owner:            Owner,
    OwnerId:          OwnerId,
    lotType:          lotType,
    createdAt:        createdAt,
    total:            total,
    formatedAddress:  formatedAddress,
  }
  assetJSON, err := json.Marshal(asset)
  if err != nil {
    return err
  }

  return ctx.GetStub().PutState(id, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Lot, error) {
  assetJSON, err := ctx.GetStub().GetState(id)
  if err != nil {
    return nil, fmt.Errorf("failed to read from world state: %v", err)
  }
  if assetJSON == nil {
    return nil, fmt.Errorf("the asset %s does not exist", id)
  }

  var asset Lot
  err = json.Unmarshal(assetJSON, &asset)
  if err != nil {
    return nil, err
  }

  return &asset, nil
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
  assetJSON, err := ctx.GetStub().GetState(id)
  if err != nil {
    return false, fmt.Errorf("failed to read from world state: %v", err)
  }

  return assetJSON != nil, nil
}

// TransferAsset updates the owner field of asset with given id in world state.
func (s *SmartContract) TransferAsset(
  ctx contractapi.TransactionContextInterface,
  id string,
  newOwner string,
  newOwnerId int
) error {
  asset, err := s.ReadAsset(ctx, id)
  if err != nil {
    return err
  }

  asset.Owner = newOwner
  asset.OwnerId = newOwnerId
  assetJSON, err := json.Marshal(asset)
  if err != nil {
    return err
  }

  return ctx.GetStub().PutState(id, assetJSON)
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Lot, error) {
  // range query with empty string for startKey and endKey does an
  // open-ended query of all assets in the chaincode namespace.
  resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
  if err != nil {
    return nil, err
  }
  defer resultsIterator.Close()

  var assets []*Lot
  for resultsIterator.HasNext() {
    queryResponse, err := resultsIterator.Next()
    if err != nil {
      return nil, err
    }

    var asset Lot
    err = json.Unmarshal(queryResponse.Value, &asset)
    if err != nil {
      return nil, err
    }
    assets = append(assets, &asset)
  }

  return assets, nil
}

func main() {
  assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
  if err != nil {
    log.Panicf("Error creating asset-transfer-basic chaincode: %v", err)
  }

  if err := assetChaincode.Start(); err != nil {
    log.Panicf("Error starting asset-transfer-basic chaincode: %v", err)
  }
}
