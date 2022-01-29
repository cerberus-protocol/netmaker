package mq

import (
	"github.com/gravitl/netmaker/logger"
	"github.com/gravitl/netmaker/logic"
	"github.com/gravitl/netmaker/netclient/ncutils"
)

func decryptMsg(nodeid string, msg []byte) ([]byte, error) {
	logger.Log(0, "found message for decryption: %s \n", string(msg))
	trafficKey, trafficErr := logic.RetrieveServerTrafficKey()
	if trafficErr != nil {
		return nil, trafficErr
	}
	return ncutils.DecryptWithPrivateKey(msg, &trafficKey), nil
}

func encryptMsg(nodeid string, msg []byte) ([]byte, error) {
	var node, err = logic.GetNodeByID(nodeid)
	if err != nil {
		return nil, err
	}
	var key, fetchErr = logic.RetrieveNodeTrafficKey(node.ID) // use nodes traffic key to encrypt msg
	if fetchErr != nil {
		return nil, fetchErr
	}
	encrypted, encryptErr := ncutils.EncryptWithPublicKey(msg, &key)
	if encryptErr != nil {
		return nil, encryptErr
	}
	return encrypted, nil
}

func publish(nodeid string, dest string, msg []byte) error {
	client := SetupMQTT()
	defer client.Disconnect(250)
	encrypted, encryptErr := encryptMsg(nodeid, msg)
	if encryptErr != nil {
		return encryptErr
	}
	if token := client.Publish(dest, 0, false, encrypted); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
