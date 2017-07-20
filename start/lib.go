package main

import (
    "encoding/json"
    "errors"
    "strconv"

    "github.com/hyperledger/fabric/core/chaincode/shim"
)

// ============================================================================================================================
// Get Marble - get a marble asset from ledger
// ============================================================================================================================
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
        return owner, errors.New("Owner does not exist - " + id + ", '" + owner.Username + "' '" + owner.Company + "'")
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
