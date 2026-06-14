package chaincode

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)



type Status int

const (
	DISABLED Status = iota
	ENABLED
	WAITING
	COMPLETED
)



type Message struct {
	ID       string       `json:"id"`
	Sender   string       `json:"sender"`
	Receiver string       `json:"receiver"`
	State    Status       `json:"state"`
}


type Gateway struct {
	ID    string       `json:"id"`
	State Status       `json:"state"`
}


type Event struct {
	ID    string       `json:"id"`
	State Status       `json:"state"`
}



type BPV struct {
	Message1sg5ubaMsg2 string `json:"Message_1sg5uba_msg2"`
	Message17b0k54Msg1 string `json:"Message_17b0k54_msg1"`
}


type SmartContract struct {
	contractapi.Contract
}


type ProcessInstance struct {
	ID        string                 `json:"id"`
	StartedBy string                 `json:"startedBy"`

	Messages  map[string]Message     `json:"messages"`
	Events    map[string]Event       `json:"events"`
	Gateways  map[string]Gateway     `json:"gateways"`
	Memory    BPV                    `json:"memory"`
}



func (cc *SmartContract) getInstance(ctx contractapi.TransactionContextInterface, instanceID string) (*ProcessInstance, error) {
	data, err := ctx.GetStub().GetState(instanceID)
	if err != nil { 
		return nil, fmt.Errorf("failed to read instance: %v", err) 
	}
	if data == nil { 
		return nil, fmt.Errorf("instance %s not found", instanceID) 
	}

	var inst ProcessInstance
	if err := json.Unmarshal(data, &inst); err != nil { 
		return nil, fmt.Errorf("failed to unmarshal instance: %v", err) 
	}
	return &inst, nil
}

func (cc *SmartContract) putInstance(ctx contractapi.TransactionContextInterface, inst *ProcessInstance) error {
	data, err := json.Marshal(inst)
	if err != nil { 
		return fmt.Errorf("failed to marshal instance: %v", err) 
	}
	return ctx.GetStub().PutState(inst.ID, data)
}

func (cc *SmartContract) IsCompleted(inst *ProcessInstance, elementID string) bool {
	if msg, ok := inst.Messages[elementID]; ok { 
		return msg.State == COMPLETED 
	}
	if evt, ok := inst.Events[elementID]; ok { 
		return evt.State == COMPLETED 
	}
	if gtw, ok := inst.Gateways[elementID]; ok { 
		return gtw.State == COMPLETED 
	}
	return false
}


func (cc *SmartContract) GetInstanceState(
    ctx contractapi.TransactionContextInterface,
    instanceID string,
) (*ProcessInstance, error) {
    return cc.getInstance(ctx, instanceID)
}

func (cc *SmartContract) Message_17b0k54_Send(ctx contractapi.TransactionContextInterface, instanceID string, Message17b0k54Msg1 string) error {

	inst, err := cc.getInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to load instance: %v", err)
	}


	msg, ok := inst.Messages["Message_17b0k54"]
	if !ok {
		return fmt.Errorf("message Message_17b0k54 not found in instance")
	}

	if msg.State != ENABLED {
		return fmt.Errorf("message Message_17b0k54 is not enabled (current state: %v)", msg.State)
	}

	clientMSP, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}
	if clientMSP != msg.Sender {
		return fmt.Errorf("participant %s is not authorized to send message Message_17b0k54 (expected sender: %s)",
			clientMSP, msg.Sender)
	}


	inst.Memory.Message17b0k54Msg1 = Message17b0k54Msg1

	msg.State = WAITING
	inst.Messages["Message_17b0k54"] = msg

	ctx.GetStub().SetEvent("Message_17b0k54_sent:"+instanceID, []byte("Message sent"))



	return cc.putInstance(ctx, inst)
}

func (cc *SmartContract) Message_17b0k54_Confirm(ctx contractapi.TransactionContextInterface, instanceID string) error {
	inst, err := cc.getInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to load instance: %v", err)
	}

	msg, ok := inst.Messages["Message_17b0k54"]
	if !ok {
		return fmt.Errorf("message Message_17b0k54 not found in instance")
	}

	if msg.State != WAITING {
		return fmt.Errorf("message Message_17b0k54 is not waiting for confirmation (current state: %v)", msg.State)
	}

	clientMSP, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}
	if clientMSP != msg.Receiver {
		return fmt.Errorf("participant %s is not authorized to confirm message Message_17b0k54", clientMSP)
	}

	msg.State = COMPLETED
	inst.Messages["Message_17b0k54"] = msg

	ctx.GetStub().SetEvent("Message_17b0k54_confirmed:"+instanceID, []byte("Message confirmed"))


{
    elem := inst.Gateways["Gateway_1od6sj0"]
    elem.State = ENABLED
    inst.Gateways["Gateway_1od6sj0"] = elem
}



	return cc.putInstance(ctx, inst)
}

func (cc *SmartContract) Message_1sg5uba_Send(ctx contractapi.TransactionContextInterface, instanceID string, Message1sg5ubaMsg2 string) error {

	inst, err := cc.getInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to load instance: %v", err)
	}


	msg, ok := inst.Messages["Message_1sg5uba"]
	if !ok {
		return fmt.Errorf("message Message_1sg5uba not found in instance")
	}

	if msg.State != ENABLED {
		return fmt.Errorf("message Message_1sg5uba is not enabled (current state: %v)", msg.State)
	}

	clientMSP, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}
	if clientMSP != msg.Sender {
		return fmt.Errorf("participant %s is not authorized to send message Message_1sg5uba (expected sender: %s)",
			clientMSP, msg.Sender)
	}


	inst.Memory.Message1sg5ubaMsg2 = Message1sg5ubaMsg2

	msg.State = WAITING
	inst.Messages["Message_1sg5uba"] = msg

	ctx.GetStub().SetEvent("Message_1sg5uba_sent:"+instanceID, []byte("Message sent"))



	return cc.putInstance(ctx, inst)
}

func (cc *SmartContract) Message_1sg5uba_Confirm(ctx contractapi.TransactionContextInterface, instanceID string) error {
	inst, err := cc.getInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to load instance: %v", err)
	}

	msg, ok := inst.Messages["Message_1sg5uba"]
	if !ok {
		return fmt.Errorf("message Message_1sg5uba not found in instance")
	}

	if msg.State != WAITING {
		return fmt.Errorf("message Message_1sg5uba is not waiting for confirmation (current state: %v)", msg.State)
	}

	clientMSP, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}
	if clientMSP != msg.Receiver {
		return fmt.Errorf("participant %s is not authorized to confirm message Message_1sg5uba", clientMSP)
	}

	msg.State = COMPLETED
	inst.Messages["Message_1sg5uba"] = msg

	ctx.GetStub().SetEvent("Message_1sg5uba_confirmed:"+instanceID, []byte("Message confirmed"))


{
    elem := inst.Gateways["Gateway_1od6sj0"]
    elem.State = ENABLED
    inst.Gateways["Gateway_1od6sj0"] = elem
}



	return cc.putInstance(ctx, inst)
}

func (cc *SmartContract) Gateway_03677iw(ctx contractapi.TransactionContextInterface, instanceID string) error {
	inst, err := cc.getInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to load instance: %v", err)
	}

	gw, ok := inst.Gateways["Gateway_03677iw"]
	if !ok {
		return fmt.Errorf("gateway Gateway_03677iw not found in instance")
	}

	if gw.State != ENABLED {
		return fmt.Errorf("gateway Gateway_03677iw is not enabled (current state: %v)", gw.State)
	}

	gw.State = COMPLETED
	inst.Gateways["Gateway_03677iw"] = gw


{
    elem := inst.Messages["Message_17b0k54"]
    elem.State = ENABLED
    inst.Messages["Message_17b0k54"] = elem
}

{
    elem := inst.Messages["Message_1sg5uba"]
    elem.State = ENABLED
    inst.Messages["Message_1sg5uba"] = elem
}


	ctx.GetStub().SetEvent("Gateway_03677iw_executed:"+instanceID, []byte("Gateway executed"))



	return cc.putInstance(ctx, inst)
}

func (cc *SmartContract) Gateway_1od6sj0(ctx contractapi.TransactionContextInterface, instanceID string) error {
	inst, err := cc.getInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to load instance: %v", err)
	}

	gw, ok := inst.Gateways["Gateway_1od6sj0"]
	if !ok {
		return fmt.Errorf("gateway Gateway_1od6sj0 not found in instance")
	}

	if gw.State != ENABLED {
		return fmt.Errorf("gateway Gateway_1od6sj0 is not enabled (current state: %v)", gw.State)
	}

	

	if inst.Messages["Message_1sg5uba"].State != COMPLETED { return nil }
	if inst.Messages["Message_17b0k54"].State != COMPLETED { return nil }
	
	gw.State = COMPLETED
	inst.Gateways["Gateway_1od6sj0"] = gw

{
    elem := inst.Events["Event_1a12slo"]
    elem.State = ENABLED
    inst.Events["Event_1a12slo"] = elem
}


	ctx.GetStub().SetEvent("Gateway_1od6sj0_executed:"+instanceID, []byte("Gateway executed"))



	return cc.putInstance(ctx, inst)
}

func (cc *SmartContract) Event_10ga7gq(ctx contractapi.TransactionContextInterface, instanceID string) error {
	inst, err := cc.getInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to load instance: %v", err)
	}

	event, ok := inst.Events["Event_10ga7gq"]
	if !ok {
		return fmt.Errorf("event Event_10ga7gq not found in instance")
	}

	if event.State != ENABLED {
		return fmt.Errorf("event Event_10ga7gq is not enabled (current state: %v)", event.State)
	}

	event.State = COMPLETED
	inst.Events["Event_10ga7gq"] = event

	ctx.GetStub().SetEvent("Event_10ga7gq_executed:"+instanceID, []byte("Event executed"))


{
    elem := inst.Gateways["Gateway_03677iw"]
    elem.State = ENABLED
    inst.Gateways["Gateway_03677iw"] = elem
}



	return cc.putInstance(ctx, inst)
}

func (cc *SmartContract) Event_1a12slo(ctx contractapi.TransactionContextInterface, instanceID string) error {
	inst, err := cc.getInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to load instance: %v", err)
	}

	event, ok := inst.Events["Event_1a12slo"]
	if !ok {
		return fmt.Errorf("event Event_1a12slo not found in instance")
	}

	if event.State != ENABLED {
		return fmt.Errorf("event Event_1a12slo is not enabled (current state: %v)", event.State)
	}

	event.State = COMPLETED
	inst.Events["Event_1a12slo"] = event

	ctx.GetStub().SetEvent("Event_1a12slo_executed:"+instanceID, []byte("Event executed"))





	return cc.putInstance(ctx, inst)
}



func (cc *SmartContract) CreateInstance(ctx contractapi.TransactionContextInterface, instanceID string) error {
	stub := ctx.GetStub()

	existing, err := stub.GetState(instanceID)
	if err != nil {
		return fmt.Errorf("failed to check instance: %v", err)
	}
	if existing != nil {
		return fmt.Errorf("instance %s already exists", instanceID)
	}

	initiator, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client identity: %v", err)
	}

	inst := ProcessInstance{
		ID:        instanceID,
		StartedBy: initiator,

		Messages: map[string]Message{
		"Message_1sg5uba": {ID: "Message_1sg5uba", Sender: "Org1MSP", Receiver: "Org2MSP", State: DISABLED},
		"Message_17b0k54": {ID: "Message_17b0k54", Sender: "Org1MSP", Receiver: "Org2MSP", State: DISABLED},
		},
		Events: map[string]Event{
		"Event_10ga7gq": {ID: "Event_10ga7gq", State: ENABLED},
		"Event_1a12slo": {ID: "Event_1a12slo", State: DISABLED},
		},
		Gateways: map[string]Gateway{
		"Gateway_03677iw": {ID: "Gateway_03677iw", State: DISABLED},
		"Gateway_1od6sj0": {ID: "Gateway_1od6sj0", State: DISABLED},
		},
		Memory:   BPV{},
	}

	data, err := json.Marshal(inst)
	if err != nil {
		return fmt.Errorf("failed to marshal instance: %v", err)
	}
	return stub.PutState(instanceID, data)
}
