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
	Message0dlqkdeMsgP2 bool `json:"Message_0dlqkde_msg_p2"`
	Message1e0g1sxMsgP1 bool `json:"Message_1e0g1sx_msg_p1"`
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

func (cc *SmartContract) Message_1e0g1sx_Send(ctx contractapi.TransactionContextInterface, instanceID string, Message1e0g1sxMsgP1 bool) error {

	inst, err := cc.getInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to load instance: %v", err)
	}


	msg, ok := inst.Messages["Message_1e0g1sx"]
	if !ok {
		return fmt.Errorf("message Message_1e0g1sx not found in instance")
	}

	if msg.State != ENABLED {
		return fmt.Errorf("message Message_1e0g1sx is not enabled (current state: %v)", msg.State)
	}

	clientMSP, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}
	if clientMSP != msg.Sender {
		return fmt.Errorf("participant %s is not authorized to send message Message_1e0g1sx (expected sender: %s)",
			clientMSP, msg.Sender)
	}


	inst.Memory.Message1e0g1sxMsgP1 = Message1e0g1sxMsgP1

	msg.State = WAITING
	inst.Messages["Message_1e0g1sx"] = msg

	ctx.GetStub().SetEvent("Message_1e0g1sx_sent:"+instanceID, []byte("Message sent"))



	return cc.putInstance(ctx, inst)
}

func (cc *SmartContract) Message_1e0g1sx_Confirm(ctx contractapi.TransactionContextInterface, instanceID string) error {
	inst, err := cc.getInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to load instance: %v", err)
	}

	msg, ok := inst.Messages["Message_1e0g1sx"]
	if !ok {
		return fmt.Errorf("message Message_1e0g1sx not found in instance")
	}

	if msg.State != WAITING {
		return fmt.Errorf("message Message_1e0g1sx is not waiting for confirmation (current state: %v)", msg.State)
	}

	clientMSP, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}
	if clientMSP != msg.Receiver {
		return fmt.Errorf("participant %s is not authorized to confirm message Message_1e0g1sx", clientMSP)
	}

	msg.State = COMPLETED
	inst.Messages["Message_1e0g1sx"] = msg

	ctx.GetStub().SetEvent("Message_1e0g1sx_confirmed:"+instanceID, []byte("Message confirmed"))


{
    elem := inst.Messages["Message_0dlqkde"]
    elem.State = ENABLED
    inst.Messages["Message_0dlqkde"] = elem
}



	return cc.putInstance(ctx, inst)
}

func (cc *SmartContract) Message_0dlqkde_Send(ctx contractapi.TransactionContextInterface, instanceID string, Message0dlqkdeMsgP2 bool) error {

	inst, err := cc.getInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to load instance: %v", err)
	}


	msg, ok := inst.Messages["Message_0dlqkde"]
	if !ok {
		return fmt.Errorf("message Message_0dlqkde not found in instance")
	}

	if msg.State != ENABLED {
		return fmt.Errorf("message Message_0dlqkde is not enabled (current state: %v)", msg.State)
	}

	clientMSP, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}
	if clientMSP != msg.Sender {
		return fmt.Errorf("participant %s is not authorized to send message Message_0dlqkde (expected sender: %s)",
			clientMSP, msg.Sender)
	}


	inst.Memory.Message0dlqkdeMsgP2 = Message0dlqkdeMsgP2

	msg.State = WAITING
	inst.Messages["Message_0dlqkde"] = msg

	ctx.GetStub().SetEvent("Message_0dlqkde_sent:"+instanceID, []byte("Message sent"))



	return cc.putInstance(ctx, inst)
}

func (cc *SmartContract) Message_0dlqkde_Confirm(ctx contractapi.TransactionContextInterface, instanceID string) error {
	inst, err := cc.getInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to load instance: %v", err)
	}

	msg, ok := inst.Messages["Message_0dlqkde"]
	if !ok {
		return fmt.Errorf("message Message_0dlqkde not found in instance")
	}

	if msg.State != WAITING {
		return fmt.Errorf("message Message_0dlqkde is not waiting for confirmation (current state: %v)", msg.State)
	}

	clientMSP, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}
	if clientMSP != msg.Receiver {
		return fmt.Errorf("participant %s is not authorized to confirm message Message_0dlqkde", clientMSP)
	}

	msg.State = COMPLETED
	inst.Messages["Message_0dlqkde"] = msg

	ctx.GetStub().SetEvent("Message_0dlqkde_confirmed:"+instanceID, []byte("Message confirmed"))


{
    elem := inst.Gateways["Gateway_0a8le7r"]
    elem.State = ENABLED
    inst.Gateways["Gateway_0a8le7r"] = elem
}



	return cc.putInstance(ctx, inst)
}

func (cc *SmartContract) Gateway_0a8le7r(ctx contractapi.TransactionContextInterface, instanceID string) error {
	inst, err := cc.getInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to load instance: %v", err)
	}

	gw, ok := inst.Gateways["Gateway_0a8le7r"]
	if !ok {
		return fmt.Errorf("gateway Gateway_0a8le7r not found in instance")
	}

	if gw.State != ENABLED {
		return fmt.Errorf("gateway Gateway_0a8le7r is not enabled (current state: %v)", gw.State)
	}

	gw.State = COMPLETED
	inst.Gateways["Gateway_0a8le7r"] = gw

	if inst.Memory.Message1e0g1sxMsgP1 == true && inst.Memory.Message0dlqkdeMsgP2 == true {

        {
            elem := inst.Events["Event_19k95ew"]
            elem.State = ENABLED
            inst.Events["Event_19k95ew"] = elem
        }
	}
	if inst.Memory.Message1e0g1sxMsgP1 == false || inst.Memory.Message0dlqkdeMsgP2 == false {

        {
            elem := inst.Events["Event_0pul0ed"]
            elem.State = ENABLED
            inst.Events["Event_0pul0ed"] = elem
        }
	}


	ctx.GetStub().SetEvent("Gateway_0a8le7r_executed:"+instanceID, []byte("Gateway executed"))



	return cc.putInstance(ctx, inst)
}

func (cc *SmartContract) Event_1bcobgg(ctx contractapi.TransactionContextInterface, instanceID string) error {
	inst, err := cc.getInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to load instance: %v", err)
	}

	event, ok := inst.Events["Event_1bcobgg"]
	if !ok {
		return fmt.Errorf("event Event_1bcobgg not found in instance")
	}

	if event.State != ENABLED {
		return fmt.Errorf("event Event_1bcobgg is not enabled (current state: %v)", event.State)
	}

	event.State = COMPLETED
	inst.Events["Event_1bcobgg"] = event

	ctx.GetStub().SetEvent("Event_1bcobgg_executed:"+instanceID, []byte("Event executed"))


{
    elem := inst.Messages["Message_1e0g1sx"]
    elem.State = ENABLED
    inst.Messages["Message_1e0g1sx"] = elem
}



	return cc.putInstance(ctx, inst)
}

func (cc *SmartContract) Event_19k95ew(ctx contractapi.TransactionContextInterface, instanceID string) error {
	inst, err := cc.getInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to load instance: %v", err)
	}

	event, ok := inst.Events["Event_19k95ew"]
	if !ok {
		return fmt.Errorf("event Event_19k95ew not found in instance")
	}

	if event.State != ENABLED {
		return fmt.Errorf("event Event_19k95ew is not enabled (current state: %v)", event.State)
	}

	event.State = COMPLETED
	inst.Events["Event_19k95ew"] = event

	ctx.GetStub().SetEvent("Event_19k95ew_executed:"+instanceID, []byte("Event executed"))





	return cc.putInstance(ctx, inst)
}

func (cc *SmartContract) Event_0pul0ed(ctx contractapi.TransactionContextInterface, instanceID string) error {
	inst, err := cc.getInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to load instance: %v", err)
	}

	event, ok := inst.Events["Event_0pul0ed"]
	if !ok {
		return fmt.Errorf("event Event_0pul0ed not found in instance")
	}

	if event.State != ENABLED {
		return fmt.Errorf("event Event_0pul0ed is not enabled (current state: %v)", event.State)
	}

	event.State = COMPLETED
	inst.Events["Event_0pul0ed"] = event

	ctx.GetStub().SetEvent("Event_0pul0ed_executed:"+instanceID, []byte("Event executed"))





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
		"Message_0dlqkde": {ID: "Message_0dlqkde", Sender: "Org2MSP", Receiver: "Org1MSP", State: DISABLED},
		"Message_1e0g1sx": {ID: "Message_1e0g1sx", Sender: "Org1MSP", Receiver: "Org2MSP", State: DISABLED},
		},
		Events: map[string]Event{
		"Event_1bcobgg": {ID: "Event_1bcobgg", State: ENABLED},
		"Event_19k95ew": {ID: "Event_19k95ew", State: DISABLED},
		"Event_0pul0ed": {ID: "Event_0pul0ed", State: DISABLED},
		},
		Gateways: map[string]Gateway{
		"Gateway_0a8le7r": {ID: "Gateway_0a8le7r", State: DISABLED},
		},
		Memory:   BPV{},
	}

	data, err := json.Marshal(inst)
	if err != nil {
		return fmt.Errorf("failed to marshal instance: %v", err)
	}
	return stub.PutState(instanceID, data)
}
