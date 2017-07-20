package main

import (
    "errors"
    "fmt"
	"encoding/json"
	"strconv"
    "github.com/hyperledger/fabric/core/chaincode/shim"
)

// ----- Events ----- //
type Event struct {
    ObjectType string        `json:"docType"` //field for couchdb
    Id       string          `json:"id"`      //the fieldtags are needed to keep case from bouncing around
    Amount      int           `json:"size"`    //size in mm of marble
    Owner      Owner `json:"owner"`
}

// ----- Owners ----- //
type Owner struct {
    ObjectType string `json:"docType"`     //field for couchdb
    Id         string `json:"id"`
    Username   string `json:"username"`
}

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func main() {
    err := shim.Start(new(SimpleChaincode))
    if err != nil {
        fmt.Printf("Error starting Simple chaincode: %s", err)
    }
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting 1")
    }

    err := stub.PutState("hello_world", []byte(args[0]))
    if err != nil {
        return nil, err
    }

    return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("invoke is running " + function)

    // Handle different functions
    if function == "init" {
        return t.Init(stub, "init", args)
    } else if function == "write" {
        return t.write(stub, args)
    }
    fmt.Println("invoke did not find func: " + function)

    return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("query is running " + function)

    // Handle different functions
    if function == "read" { //read a variable
        return t.read(stub, args)
    }
    fmt.Println("query did not find func: " + function)

    return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var key, value string
    var err error
    fmt.Println("running write()")

    if len(args) != 2 {
        return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
    }

    key = args[0] //rename for funsies
    value = args[1]
    err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
    if err != nil {
        return nil, err
    }
    return nil, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var key, jsonResp string
    var err error

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }

    key = args[0]
    valAsbytes, err := stub.GetState(key)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
    }

    return valAsbytes, nil
}

func get_event(stub shim.ChaincodeStubInterface, id string) (Event, error) {
    var event Event
    eventAsBytes, err := stub.GetState(id)                  //getState retreives a key/value from the ledger
    if err != nil {                                          //this seems to always succeed, even if key didn't exist
        return event, errors.New("Failed to find marble - " + id)
    }
    json.Unmarshal(eventAsBytes, &event)                   //un stringify it aka JSON.parse()

    if event.Id != id {                                     //test if marble is actually here or just nil
        return event, errors.New("Event does not exist - " + id)
    }

    return event, nil
}

// ============================================================================================================================
// Get Owner - get the owner asset from ledger
// ============================================================================================================================
func get_owner(stub shim.ChaincodeStubInterface, id string) (Owner, error) {
    var owner Owner
    ownerAsBytes, err := stub.GetState(id)                     //getState retreives a key/value from the ledger
    if err != nil {                                            //this seems to always succeed, even if key didn't exist
        return owner, errors.New("Failed to get owner - " + id)
    }
    json.Unmarshal(ownerAsBytes, &owner)                       //un stringify it aka JSON.parse()

    if len(owner.Username) == 0 {                              //test if owner is actually here or just nil
        return owner, errors.New("Owner does not exist - " + id + ", '" + owner.Username + "' '" + "'")
    }
    
    return owner, nil
}

// ========================================================
// Input Sanitation - dumb input checking, look for empty strings
// ========================================================
func sanitize_arguments(strs []string) error{
    for i, val:= range strs {
        if len(val) <= 0 {
            return errors.New("Argument " + strconv.Itoa(i) + " must be a non-empty string")
        }
        if len(val) > 32 {
            return errors.New("Argument " + strconv.Itoa(i) + " must be <= 32 characters")
        }
    }
    return nil
}
