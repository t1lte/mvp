from typing import List, Dict, Any
from .function_builder import FunctionBuilder


class ChaincodeAssembler:
    def __init__(self):
        self.builder = FunctionBuilder()
        self.imports = self._get_standard_imports()
        self.structs: List[str] = []
        self.functions: List[str] = []
        self.init_functions: List[str] = []
        self.state_memory_fields: str = ""

    def _get_standard_imports(self) -> str:
        return '''import (
	"encoding/json"
	"fmt"
	"time"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)'''

    def set_state_memory_fields(self, fields: str) -> None:
        self.state_memory_fields = fields

    def add_struct_definition(self, struct_name: str, fields: str) -> None:
        struct_def = f"type {struct_name} struct {{\n{fields}\n}}"
        self.structs.append(struct_def)

    def add_contract_definition(self) -> None:
        self.structs.append("type SmartContract struct {\n\tcontractapi.Contract\n}")

    def add_process_instance_struct(self) -> None:
        code = '''
type ProcessInstance struct {
	ID        string                 `json:"id"`
	StartedBy string                 `json:"startedBy"`
	StartedAt int64                  `json:"startedAt"`

	Messages  map[string]Message     `json:"messages"`
	Events    map[string]Event       `json:"events"`
	Gateways  map[string]Gateway     `json:"gateways"`
	Memory    BPV                    `json:"memory"`
}'''
        self.structs.append(code)

    def add_instance_management_functions(self) -> None:
        code = '''
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
}'''
        self.functions.append(code)

    def add_create_instance_function(self,
                                     start_event_id: str,
                                     end_events: List[str],
                                     messages: List[Dict],
                                     gateways: List[str]) -> None:

        messages_init = []
        for m in messages:
            messages_init.append(
                f'\t\t"{m["id"]}": {{ID: "{m["id"]}", Sender: "{m["sender"]}", Receiver: "{m["receiver"]}", State: DISABLED}},')
        messages_code = "\n".join(messages_init) if messages_init else ""
        events_init = []
        if start_event_id:
            events_init.append(f'\t\t"{start_event_id}": {{ID: "{start_event_id}", State: ENABLED}},')
        for ev in end_events:
            events_init.append(f'\t\t"{ev}": {{ID: "{ev}", State: DISABLED}},')
        events_code = "\n".join(events_init) if events_init else ""
        gateways_init = []
        for gw in gateways:
            gateways_init.append(f'\t\t"{gw}": {{ID: "{gw}", State: DISABLED}},')
        gateways_code = "\n".join(gateways_init) if gateways_init else ""

        code = f'''
func (cc *SmartContract) CreateInstance(ctx contractapi.TransactionContextInterface, instanceID string) error {{
	stub := ctx.GetStub()

	existing, err := stub.GetState(instanceID)
	if err != nil {{
		return fmt.Errorf("failed to check instance: %v", err)
	}}
	if existing != nil {{
		return fmt.Errorf("instance %s already exists", instanceID)
	}}

	initiator, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {{
		return fmt.Errorf("failed to get client identity: %v", err)
	}}

	inst := ProcessInstance{{
		ID:        instanceID,
		StartedBy: initiator,
		StartedAt: time.Now().Unix(),

		Messages: map[string]Message{{
{messages_code}
		}},
		Events: map[string]Event{{
{events_code}
		}},
		Gateways: map[string]Gateway{{
{gateways_code}
		}},
		Memory:   BPV{{}},
	}}

	data, err := json.Marshal(inst)
	if err != nil {{
		return fmt.Errorf("failed to marshal instance: %v", err)
	}}
	return stub.PutState(instanceID, data)
}}'''

        self.init_functions.append(code)

    def add_function(self, function_code: str) -> None:
        self.functions.append(function_code)

    def assemble(self, package_name: str = "chaincode") -> str:
        if self.state_memory_fields:
            state_memory_struct = f"type BPV struct {{\n{self.state_memory_fields}\n}}"
        else:
            state_memory_struct = "type BPV struct {}"

        parts = [
            f"package {package_name}\n",
            self.imports,
            "\n\n",
            "type Status int"
            "\n\nconst (",
            "\tDISABLED Status = iota",
            "\tENABLED",
            "\tWAITING",
            "\tCOMPLETED",
            ")",
            "\n\n",
            "type Message struct {",
            "\tID       string       `json:\"id\"`",
            "\tSender   string       `json:\"sender\"`",
            "\tReceiver string       `json:\"receiver\"`",
            "\tState    Status       `json:\"state\"`",
            "}",
            "\n\ntype Gateway struct {",
            "\tID    string       `json:\"id\"`",
            "\tState Status       `json:\"state\"`",
            "}",
            "\n\ntype Event struct {",
            "\tID    string       `json:\"id\"`",
            "\tState Status       `json:\"state\"`",
            "}",
            "\n\n",
            state_memory_struct,
            "\n\n" + "\n\n".join(self.structs),
            "\n\n" + "\n\n".join(self.functions),
            "\n\n" + "\n\n".join(self.init_functions) if self.init_functions else "",
        ]

        return "\n".join(parts)