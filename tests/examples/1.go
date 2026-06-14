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
	Message0fw90ysResponse string `json:"Message_0fw90ys_Response"`
	Message1ml5k9rRequest string `json:"Message_1ml5k9r_Request"`
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

func (cc *SmartContract) Message_1ml5k9r_Send(ctx contractapi.TransactionContextInterface, instanceID string, Message1ml5k9rRequest string) error {

	inst, err := cc.getInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to load instance: %v", err)
	}


	msg, ok := inst.Messages["Message_1ml5k9r"]
	if !ok {
		return fmt.Errorf("message Message_1ml5k9r not found in instance")
	}

	if msg.State != ENABLED {
		return fmt.Errorf("message Message_1ml5k9r is not enabled (current state: %v)", msg.State)
	}

	clientMSP, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}
	if clientMSP != msg.Sender {
		return fmt.Errorf("participant %s is not authorized to send message Message_1ml5k9r (expected sender: %s)",
			clientMSP, msg.Sender)
	}


	inst.Memory.Message1ml5k9rRequest = Message1ml5k9rRequest

	msg.State = WAITING
	inst.Messages["Message_1ml5k9r"] = msg

	ctx.GetStub().SetEvent("Message_1ml5k9r_sent:"+instanceID, []byte("Message sent"))



	return cc.putInstance(ctx, inst)
}

func (cc *SmartContract) Message_1ml5k9r_Confirm(ctx contractapi.TransactionContextInterface, instanceID string) error {
	inst, err := cc.getInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to load instance: %v", err)
	}

	msg, ok := inst.Messages["Message_1ml5k9r"]
	if !ok {
		return fmt.Errorf("message Message_1ml5k9r not found in instance")
	}

	if msg.State != WAITING {
		return fmt.Errorf("message Message_1ml5k9r is not waiting for confirmation (current state: %v)", msg.State)
	}

	clientMSP, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}
	if clientMSP != msg.Receiver {
		return fmt.Errorf("participant %s is not authorized to confirm message Message_1ml5k9r", clientMSP)
	}

	msg.State = COMPLETED
	inst.Messages["Message_1ml5k9r"] = msg

	ctx.GetStub().SetEvent("Message_1ml5k9r_confirmed:"+instanceID, []byte("Message confirmed"))


{
    elem := inst.Messages["Message_0fw90ys"]
    elem.State = ENABLED
    inst.Messages["Message_0fw90ys"] = elem
}



	return cc.putInstance(ctx, inst)
}

func (cc *SmartContract) Message_0fw90ys_Send(ctx contractapi.TransactionContextInterface, instanceID string, Message0fw90ysResponse string) error {

	inst, err := cc.getInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to load instance: %v", err)
	}


	msg, ok := inst.Messages["Message_0fw90ys"]
	if !ok {
		return fmt.Errorf("message Message_0fw90ys not found in instance")
	}

	if msg.State != ENABLED {
		return fmt.Errorf("message Message_0fw90ys is not enabled (current state: %v)", msg.State)
	}

	clientMSP, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}
	if clientMSP != msg.Sender {
		return fmt.Errorf("participant %s is not authorized to send message Message_0fw90ys (expected sender: %s)",
			clientMSP, msg.Sender)
	}


	inst.Memory.Message0fw90ysResponse = Message0fw90ysResponse

	msg.State = WAITING
	inst.Messages["Message_0fw90ys"] = msg

	ctx.GetStub().SetEvent("Message_0fw90ys_sent:"+instanceID, []byte("Message sent"))



	return cc.putInstance(ctx, inst)
}

func (cc *SmartContract) Message_0fw90ys_Confirm(ctx contractapi.TransactionContextInterface, instanceID string) error {
	inst, err := cc.getInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to load instance: %v", err)
	}

	msg, ok := inst.Messages["Message_0fw90ys"]
	if !ok {
		return fmt.Errorf("message Message_0fw90ys not found in instance")
	}

	if msg.State != WAITING {
		return fmt.Errorf("message Message_0fw90ys is not waiting for confirmation (current state: %v)", msg.State)
	}

	clientMSP, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}
	if clientMSP != msg.Receiver {
		return fmt.Errorf("participant %s is not authorized to confirm message Message_0fw90ys", clientMSP)
	}

	msg.State = COMPLETED
	inst.Messages["Message_0fw90ys"] = msg

	ctx.GetStub().SetEvent("Message_0fw90ys_confirmed:"+instanceID, []byte("Message confirmed"))


{
    elem := inst.Events["Event_14qay4i"]
    elem.State = ENABLED
    inst.Events["Event_14qay4i"] = elem
}



	return cc.putInstance(ctx, inst)
}

func (cc *SmartContract) Event_0e5bcqi(ctx contractapi.TransactionContextInterface, instanceID string) error {
	inst, err := cc.getInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to load instance: %v", err)
	}

	event, ok := inst.Events["Event_0e5bcqi"]
	if !ok {
		return fmt.Errorf("event Event_0e5bcqi not found in instance")
	}

	if event.State != ENABLED {
		return fmt.Errorf("event Event_0e5bcqi is not enabled (current state: %v)", event.State)
	}

	event.State = COMPLETED
	inst.Events["Event_0e5bcqi"] = event

	ctx.GetStub().SetEvent("Event_0e5bcqi_executed:"+instanceID, []byte("Event executed"))


{
    elem := inst.Messages["Message_1ml5k9r"]
    elem.State = ENABLED
    inst.Messages["Message_1ml5k9r"] = elem
}



	return cc.putInstance(ctx, inst)
}

func (cc *SmartContract) Event_14qay4i(ctx contractapi.TransactionContextInterface, instanceID string) error {
	inst, err := cc.getInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to load instance: %v", err)
	}

	event, ok := inst.Events["Event_14qay4i"]
	if !ok {
		return fmt.Errorf("event Event_14qay4i not found in instance")
	}

	if event.State != ENABLED {
		return fmt.Errorf("event Event_14qay4i is not enabled (current state: %v)", event.State)
	}

	event.State = COMPLETED
	inst.Events["Event_14qay4i"] = event

	ctx.GetStub().SetEvent("Event_14qay4i_executed:"+instanceID, []byte("Event executed"))





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
		"Message_0fw90ys": {ID: "Message_0fw90ys", Sender: "Org2MSP", Receiver: "Org1MSP", State: DISABLED},
		"Message_1ml5k9r": {ID: "Message_1ml5k9r", Sender: "Org1MSP", Receiver: "Org2MSP", State: DISABLED},
		},
		Events: map[string]Event{
		"Event_0e5bcqi": {ID: "Event_0e5bcqi", State: ENABLED},
		"Event_14qay4i": {ID: "Event_14qay4i", State: DISABLED},
		},
		Gateways: map[string]Gateway{

		},
		Memory:   BPV{},
	}

	data, err := json.Marshal(inst)
	if err != nil {
		return fmt.Errorf("failed to marshal instance: %v", err)
	}
	return stub.PutState(instanceID, data)
}
