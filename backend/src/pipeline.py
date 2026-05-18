from pathlib import Path
from typing import Dict, Optional, List

from src.bpmn.parser import BPMNParser
from src.analysis import SemanticAnalyzer
from src.fsm import FSMGenerator
from src.codegen.assembler import ChaincodeAssembler
from src.codegen.function_builder import FunctionBuilder


class BPMNTranslator:
    def __init__(self, bpmn_path: str, participant_map: Dict[str, str] = None):
        self.bpmn_path = Path(bpmn_path)
        self.participant_map = participant_map or {}
        self.choreography = None
        self.analyzer = SemanticAnalyzer()
        self.fsm_gen = None
        self.assembler = ChaincodeAssembler()
        self.builder = FunctionBuilder()

        self.parameters: Dict = {}
        self.conditions: Dict = {}
        self.element_logic: Dict = {}

    def _get_msp_for_participant(self, participant_id: str) -> str:
        if not participant_id:
            return "Org1MSP"
        return self.participant_map.get(participant_id, "Org1MSP")

    def load(self) -> None:
        with open(self.bpmn_path, "r", encoding="utf-8") as f:
            xml_content = f.read()
        parser = BPMNParser(xml_content)
        self.choreography = parser.parse()

    def analyze(self) -> None:
        self.parameters, self.conditions = self.analyzer.analyze(self.choreography)
        self.fsm_gen = FSMGenerator(self.choreography, self.analyzer)
        self.element_logic = self.fsm_gen.generate()

    def _generate_state_struct(self) -> None:
        fields = self.analyzer.get_go_struct_fields()
        if fields:
            self.assembler.set_state_memory_fields(fields)

    def _generate_message_functions(self) -> None:
        for task in self.choreography.tasks.values():
            if not task.init_message_flow:
                continue

            message_id = task.init_message_flow.message.id

            params = []
            for scoped_name, info in self.parameters.items():
                if info["source_message_id"] == message_id:
                    params.append({
                        "name": info["original_name"],
                        "scoped": self._to_go_field_name(scoped_name),
                        "go_type": info["go_type"]
                    })

            validation = self._generate_validation_for_params(params)
            logic = self.element_logic.get(task.id, {})
            next_code = logic.get("next_state_code", "")

            send_code = self.builder.build_message_send(
                message_id=message_id,
                participant_id=task.init_participant.id if task.init_participant else "",
                parameters=params,
                validation_code=validation,
                next_hook=""
            )
            self.assembler.add_function(send_code)

            confirm_code = self.builder.build_message_confirm(
                message_id=message_id,
                participant_id="",
                next_state_code=next_code,
                next_hook=""
            )
            self.assembler.add_function(confirm_code)

    def _generate_validation_for_params(self, params: List[Dict]) -> str:
        return ""

    def _generate_gateway_functions(self) -> None:
        type_map = {
            "exclusiveGateway": "ExclusiveGateway",
            "parallelGateway": "ParallelGateway"
        }
        for gw in self.choreography.gateways.values():
            logic = self.element_logic.get(gw.id, {})

            if "branches" in logic and logic["branches"]:
                code = self.builder.build_gateway_split(
                    gateway_id=gw.id,
                    gateway_type=type_map.get(gw.type.value, gw.type.value.title()),
                    branches=logic["branches"],
                    next_hook=""
                )
            else:
                incoming_count = len(gw.incomings) if hasattr(gw, 'incomings') else 1
                code = self.builder.build_gateway_merge(
                    gateway_id=gw.id,
                    gateway_type=type_map.get(gw.type.value, gw.type.value.title()),
                    incoming_count=incoming_count,
                    next_state_code=logic.get("next_state_code", ""),
                    next_hook=""
                )
            self.assembler.add_function(code)

    def _generate_event_functions(self) -> None:
        for event in self.choreography.events.values():
            logic = self.element_logic.get(event.id, {})
            code = self.builder.build_event_handler(
                event_id=event.id,
                event_type=event.type.value.capitalize(),
                next_state_code=logic.get("next_state_code", ""),
                next_hook=""
            )
            self.assembler.add_function(code)

    def _generate_create_instance(self) -> None:
        messages = []
        for mf in self.choreography.edges:
            if hasattr(mf, 'message') and mf.message:
                sender_msp = self._get_msp_for_participant(
                    mf.source.id if mf.source else ""
                )
                receiver_msp = self._get_msp_for_participant(
                    mf.target.id if mf.target else ""
                )
                messages.append({
                    "id": mf.message.id,
                    "sender": sender_msp,
                    "receiver": receiver_msp,
                    "state": "DISABLED"
                })

        gateways = list(self.choreography.gateways.keys())
        start_event = next((e.id for e in self.choreography.events.values()
                            if e.type.value == "startEvent"), "")
        end_events = [e.id for e in self.choreography.events.values()
                      if e.type.value == "endEvent"]

        self.assembler.add_create_instance_function(
            start_event_id=start_event,
            end_events=end_events,
            messages=messages,
            gateways=gateways
        )

    @staticmethod
    def _to_go_field_name(scoped_name: str) -> str:
        parts = scoped_name.split("_")
        return "".join(part.capitalize() for part in parts)

    def generate(self, output_path: str = "contract/chaincode.go") -> str:
        self.assembler.add_contract_definition()
        self._generate_state_struct()
        self.assembler.add_process_instance_struct()
        self.assembler.add_instance_management_functions()
        self._generate_create_instance()
        self._generate_message_functions()
        self._generate_gateway_functions()
        self._generate_event_functions()

        chaincode = self.assembler.assemble()
        output_file = Path(output_path)
        output_file.parent.mkdir(parents=True, exist_ok=True)
        with open(output_file, "w", encoding="utf-8") as f:
            f.write(chaincode)

        return chaincode

    def run(self, output_path: str = "contract/chaincode.go",
            participant_map: Dict[str, str] = None) -> str:
        if participant_map:
            self.participant_map = participant_map
        self.load()
        self.analyze()
        return self.generate(output_path)