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
	Message0bmupktChoice bool `json:"Message_0bmupkt_choice"`
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

func (cc *SmartContract) Message_0bmupkt_Send(ctx contractapi.TransactionContextInterface, instanceID string, Message0bmupktChoice bool) error {

	inst, err := cc.getInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to load instance: %v", err)
	}


	msg, ok := inst.Messages["Message_0bmupkt"]
	if !ok {
		return fmt.Errorf("message Message_0bmupkt not found in instance")
	}

	if msg.State != ENABLED {
		return fmt.Errorf("message Message_0bmupkt is not enabled (current state: %v)", msg.State)
	}

	clientMSP, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}
	if clientMSP != msg.Sender {
		return fmt.Errorf("participant %s is not authorized to send message Message_0bmupkt (expected sender: %s)",
			clientMSP, msg.Sender)
	}


	inst.Memory.Message0bmupktChoice = Message0bmupktChoice

	msg.State = WAITING
	inst.Messages["Message_0bmupkt"] = msg

	ctx.GetStub().SetEvent("Message_0bmupkt_sent:"+instanceID, []byte("Message sent"))



	return cc.putInstance(ctx, inst)
}

func (cc *SmartContract) Message_0bmupkt_Confirm(ctx contractapi.TransactionContextInterface, instanceID string) error {
	inst, err := cc.getInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to load instance: %v", err)
	}

	msg, ok := inst.Messages["Message_0bmupkt"]
	if !ok {
		return fmt.Errorf("message Message_0bmupkt not found in instance")
	}

	if msg.State != WAITING {
		return fmt.Errorf("message Message_0bmupkt is not waiting for confirmation (current state: %v)", msg.State)
	}

	clientMSP, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}
	if clientMSP != msg.Receiver {
		return fmt.Errorf("participant %s is not authorized to confirm message Message_0bmupkt", clientMSP)
	}

	msg.State = COMPLETED
	inst.Messages["Message_0bmupkt"] = msg

	ctx.GetStub().SetEvent("Message_0bmupkt_confirmed:"+instanceID, []byte("Message confirmed"))


{
    elem := inst.Gateways["Gateway_1rrrq3j"]
    elem.State = ENABLED
    inst.Gateways["Gateway_1rrrq3j"] = elem
}



	return cc.putInstance(ctx, inst)
}

func (cc *SmartContract) Gateway_1rrrq3j(ctx contractapi.TransactionContextInterface, instanceID string) error {
	inst, err := cc.getInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to load instance: %v", err)
	}

	gw, ok := inst.Gateways["Gateway_1rrrq3j"]
	if !ok {
		return fmt.Errorf("gateway Gateway_1rrrq3j not found in instance")
	}

	if gw.State != ENABLED {
		return fmt.Errorf("gateway Gateway_1rrrq3j is not enabled (current state: %v)", gw.State)
	}

	gw.State = COMPLETED
	inst.Gateways["Gateway_1rrrq3j"] = gw

	if inst.Memory.Message0bmupktChoice == true {

        {
            elem := inst.Events["Event_01827qy"]
            elem.State = ENABLED
            inst.Events["Event_01827qy"] = elem
        }
	}
	if inst.Memory.Message0bmupktChoice == false {

        {
            elem := inst.Events["Event_0n34w4n"]
            elem.State = ENABLED
            inst.Events["Event_0n34w4n"] = elem
        }
	}


	ctx.GetStub().SetEvent("Gateway_1rrrq3j_executed:"+instanceID, []byte("Gateway executed"))



	return cc.putInstance(ctx, inst)
}

func (cc *SmartContract) Event_1ryxo25(ctx contractapi.TransactionContextInterface, instanceID string) error {
	inst, err := cc.getInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to load instance: %v", err)
	}

	event, ok := inst.Events["Event_1ryxo25"]
	if !ok {
		return fmt.Errorf("event Event_1ryxo25 not found in instance")
	}

	if event.State != ENABLED {
		return fmt.Errorf("event Event_1ryxo25 is not enabled (current state: %v)", event.State)
	}

	event.State = COMPLETED
	inst.Events["Event_1ryxo25"] = event

	ctx.GetStub().SetEvent("Event_1ryxo25_executed:"+instanceID, []byte("Event executed"))


{
    elem := inst.Messages["Message_0bmupkt"]
    elem.State = ENABLED
    inst.Messages["Message_0bmupkt"] = elem
}



	return cc.putInstance(ctx, inst)
}

func (cc *SmartContract) Event_01827qy(ctx contractapi.TransactionContextInterface, instanceID string) error {
	inst, err := cc.getInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to load instance: %v", err)
	}

	event, ok := inst.Events["Event_01827qy"]
	if !ok {
		return fmt.Errorf("event Event_01827qy not found in instance")
	}

	if event.State != ENABLED {
		return fmt.Errorf("event Event_01827qy is not enabled (current state: %v)", event.State)
	}

	event.State = COMPLETED
	inst.Events["Event_01827qy"] = event

	ctx.GetStub().SetEvent("Event_01827qy_executed:"+instanceID, []byte("Event executed"))





	return cc.putInstance(ctx, inst)
}

func (cc *SmartContract) Event_0n34w4n(ctx contractapi.TransactionContextInterface, instanceID string) error {
	inst, err := cc.getInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to load instance: %v", err)
	}

	event, ok := inst.Events["Event_0n34w4n"]
	if !ok {
		return fmt.Errorf("event Event_0n34w4n not found in instance")
	}

	if event.State != ENABLED {
		return fmt.Errorf("event Event_0n34w4n is not enabled (current state: %v)", event.State)
	}

	event.State = COMPLETED
	inst.Events["Event_0n34w4n"] = event

	ctx.GetStub().SetEvent("Event_0n34w4n_executed:"+instanceID, []byte("Event executed"))





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
		"Message_0bmupkt": {ID: "Message_0bmupkt", Sender: "Org1MSP", Receiver: "Org2MSP", State: DISABLED},
		},
		Events: map[string]Event{
		"Event_1ryxo25": {ID: "Event_1ryxo25", State: ENABLED},
		"Event_01827qy": {ID: "Event_01827qy", State: DISABLED},
		"Event_0n34w4n": {ID: "Event_0n34w4n", State: DISABLED},
		},
		Gateways: map[string]Gateway{
		"Gateway_1rrrq3j": {ID: "Gateway_1rrrq3j", State: DISABLED},
		},
		Memory:   BPV{},
	}

	data, err := json.Marshal(inst)
	if err != nil {
		return fmt.Errorf("failed to marshal instance: %v", err)
	}
	return stub.PutState(instanceID, data)
}
