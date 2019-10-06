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
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// SmartContract struct
type SmartContract struct {
}

//AuditLog struct
type AuditLog struct {
	LogID            string `json:"logID"`
	Timestamp        string `json:"timestamp"`
	UserIdentity     string `json:"userIdentity"`
	UserRole         string `json:"userRole"`
	MSP              string `json:"msp"`
	Channel          string `json:"channel"`
	ChaincodeVersion string `json:"chaincodeVersion"`
	ChaincodeName    string `json:"chaincodeName"`
	FunctionCall     string `json:"functionCall"`
	RecordID         string `json:"recordID"`
	RecordHash       string `json:"recordHash"`
	Status           string `json:"status"`
}

//Init method
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

//Invoke method
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	function, args := APIstub.GetFunctionAndParameters()

	switch function {
	case "initLedger":
		return s.initLedger(APIstub)

	case "addAuditLogEntry":
		return s.addAuditLogEntry(APIstub, args)

	case "addAuditLogEntryFromChaincode":
		return s.addAuditLogEntryFromChaincode(APIstub, args)

	case "getAuditLogEntry":
		if hasAttribute(APIstub, "auditor") {
			return s.getAuditLogEntry(APIstub, args)
		}
		return shim.Error("Unauthorised access dectected")

	case "getAuditLogByRange":
		if hasAttribute(APIstub, "auditor") {
			return s.getAuditLogByRange(APIstub, args)
		}
		return shim.Error("Unauthorised access dectected")

	case "getHistoryForAuditLogKey":
		if hasAttribute(APIstub, "auditor") {
			return s.getHistoryForAuditLogKey(APIstub, args)
		}
		return shim.Error("Unauthorised access dectected")

	default:
		return shim.Error("Invalid Smart Contract function name.")
	}

}

//Inisilize the ledger
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	txntmsp, err := APIstub.GetTxTimestamp()
	if err != nil {
		return shim.Error("Error converting timestamp")
	}

	time := time.Unix(txntmsp.Seconds, int64(txntmsp.Nanos)).String()
	auditLog := AuditLog{LogID: "0", Timestamp: time, UserIdentity: "testUser", UserRole: "doctor", MSP: "hospital1MSP", Channel: "hospital1channel", ChaincodeVersion: "v1.0", ChaincodeName: "patient", FunctionCall: "updatePatient", RecordID: "Patient1", RecordHash: "24b941962ded00992a17d4f73b3dc2740e005af50e871479576a867af06937c5"}
	auditLogAsBytes, _ := json.Marshal(auditLog)
	APIstub.PutState("0", auditLogAsBytes)
	return shim.Success(nil)
}

// Add audtit log entry from chaincode
func (s *SmartContract) addAuditLogEntryFromChaincode(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 10 {
		return shim.Error("Incorrect number of arguments.")
	}

	txntmsp, err := APIstub.GetTxTimestamp()
	if err != nil {
		return shim.Error("Error converting timestamp")
	}

	time := time.Unix(txntmsp.Seconds, int64(txntmsp.Nanos)).String()
	auditLog := AuditLog{}
	auditLog.LogID = args[0] + time
	auditLog.Timestamp = time
	auditLog.UserIdentity = args[1]
	auditLog.UserRole = args[2]
	auditLog.MSP = args[3]
	auditLog.Channel = args[4]
	auditLog.ChaincodeVersion = args[5]
	auditLog.ChaincodeName = args[6]
	auditLog.FunctionCall = args[7]
	auditLog.RecordID = args[8]
	auditLog.RecordHash = args[9]

	auditLogAsBytes, _ := json.Marshal(auditLog)
	APIstub.PutState(args[0], auditLogAsBytes)

	return shim.Success(nil)
}

//Add audit log entry
func (s *SmartContract) addAuditLogEntry(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 11 {
		return shim.Error("Incorrect number of arguments.")
	}

	txntmsp, err := APIstub.GetTxTimestamp()
	if err != nil {
		return shim.Error("Error converting timestamp")
	}
	time := time.Unix(txntmsp.Seconds, int64(txntmsp.Nanos)).String()

	auditLog := AuditLog{}
	auditLog.LogID = args[0]
	auditLog.Timestamp = time
	auditLog.UserIdentity = args[1]
	auditLog.UserRole = args[2]
	auditLog.MSP = args[3]
	auditLog.Channel = args[4]
	auditLog.ChaincodeVersion = args[5]
	auditLog.ChaincodeName = args[6]
	auditLog.FunctionCall = args[7]
	auditLog.RecordID = args[8]
	auditLog.RecordHash = args[9]
	auditLog.Status = args[10]

	auditLogAsBytes, _ := json.Marshal(auditLog)
	APIstub.PutState(args[0], auditLogAsBytes)

	return shim.Success(nil)
}

//Get audit log entry
func (s *SmartContract) getAuditLogEntry(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.")
	}
	auditLogAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(auditLogAsBytes)
}

// Get Audit Log By Range
func (s *SmartContract) getAuditLogByRange(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments.")
	}

	resultsIterator, err := APIstub.GetStateByRange(args[0], args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	buffer := buildStateTimeLine(resultsIterator)

	fmt.Printf("- getAuditLogByRange queryResult:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())

}

//Get history for an audit log key
func (s *SmartContract) getHistoryForAuditLogKey(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments.")
	}

	resultsIterator, err := APIstub.GetHistoryForKey(args[0])

	if err != nil {
		return shim.Error(err.Error())
	}

	defer resultsIterator.Close()

	buffer := buildTimeLine(resultsIterator)

	fmt.Printf("- getHistoryForAuditKey returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())

}

//Start the chaincode
func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

//Helper methods
func hasAttribute(APIstub shim.ChaincodeStubInterface, attr string) bool {
	err := cid.AssertAttributeValue(APIstub, attr, "true")
	if err != nil {
		return false
	}
	return true
}

func buildTimeLine(resultsIterator shim.HistoryQueryIteratorInterface) bytes.Buffer {

	var buffer bytes.Buffer

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return buffer
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")
		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString(", \"value\":")
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"isDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return buffer
}

func buildStateTimeLine(resultsIterator shim.StateQueryIteratorInterface) bytes.Buffer {

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return buffer
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return buffer
}
