package logic

import (
	"crypto/rsa"
	"encoding/json"

	"github.com/gravitl/netmaker/database"
)

// == Public ==

// RetrieveServerTrafficKey - retrieves public key based on node
func RetrieveServerTrafficKey() (rsa.PrivateKey, error) {
	var telRecord, err = fetchTelemetryRecord()
	if err != nil {
		return rsa.PrivateKey{}, err
	}
	var key rsa.PrivateKey
	err = json.Unmarshal([]byte(telRecord.TrafficKey), &key)
	return key, err
}

// StoreNodeTrafficKey - stores the traffic key of node
func StoreNodeTrafficKey(nodeid string, value rsa.PublicKey) error {
	var newKey = trafficKey{
		TrafficKey: value,
	}
	var data, err = json.Marshal(&newKey)
	if err != nil {
		return err
	}
	return database.Insert(nodeid, string(data), database.TRAFFIC_TABLE_NAME)
}

// RetrieveNodeTrafficKey - fetches traffic key of a node
func RetrieveNodeTrafficKey(nodeid string) (rsa.PublicKey, error) {
	var data, err = database.FetchRecord(database.TRAFFIC_TABLE_NAME, nodeid)
	if err != nil {
		return rsa.PublicKey{}, err
	}
	var tKey trafficKey
	if err = json.Unmarshal([]byte(data), &tKey); err != nil {
		return rsa.PublicKey{}, err
	}
	return tKey.TrafficKey, nil
}

// == Private
type trafficKey struct {
	TrafficKey rsa.PublicKey `json:"traffickey" bson:"traffickey"`
}
