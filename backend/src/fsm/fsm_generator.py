from typing import Dict, List, Any, Optional
from src.bpmn.graph import (
    Choreography, NodeType, Element, SequenceFlow,
    ChoreographyTask, ExclusiveGateway, ParallelGateway
)
from src.analysis import SemanticAnalyzer


class FSMGenerator:
    def __init__(self, choreography: Choreography, analyzer: SemanticAnalyzer):
        self.choreography = choreography
        self.analyzer = analyzer
        self.element_logic: Dict[str, Dict[str, Any]] = {}

    def generate(self) -> Dict[str, Dict[str, Any]]:
        for node in self.choreography.nodes:
            logic = {}
            if node.type == NodeType.START_EVENT:
                logic = self._gen_start_event(node)
            elif node.type == NodeType.END_EVENT:
                logic = self._gen_end_event(node)
            elif node.type == NodeType.CHOREOGRAPHY_TASK:
                logic = self._gen_choreography_task(node)
            elif node.type == NodeType.EXCLUSIVE_GATEWAY:
                logic = self._gen_exclusive_gateway(node)
            elif node.type == NodeType.PARALLEL_GATEWAY:
                logic = self._gen_parallel_gateway(node)

            self.element_logic[node.id] = logic
        return self.element_logic

    def _get_outgoing_flows(self, node: Element) -> List[SequenceFlow]:
        return [
            f for f in self.choreography.edges
            if isinstance(f, SequenceFlow) and f.source == node
        ]

    def _get_incoming_flows(self, node: Element) -> List[SequenceFlow]:
        return [
            f for f in self.choreography.edges
            if isinstance(f, SequenceFlow) and f.target == node
        ]

    def _resolve_target_id(self, target_type: str, target_element: Element) -> str:
        if target_type == "choreographyTask" and target_element:
            try:
                msg_id = target_element.init_message_flow.message.id
                if msg_id:
                    return msg_id
            except AttributeError:
                pass
        return target_element.id if target_element else ""

    def _gen_enable_code(self, target_id: str, target_type: Optional[str] = None,
                         target_element: Optional[Element] = None) -> str:
        final_id = self._resolve_target_id(target_type, target_element)

        if target_type == "choreographyTask" or "message" in str(target_type).lower():
            container = "Messages"
        elif "gateway" in str(target_type).lower():
            container = "Gateways"
        else:
            container = "Events"

        return f"""
{{
    elem := inst.{container}["{final_id}"]
    elem.State = ENABLED
    inst.{container}["{final_id}"] = elem
}}"""

    def _gen_start_event(self, event: Element) -> Dict[str, Any]:
        flows = self._get_outgoing_flows(event)
        if not flows:
            return {"next_state_code": "", "pre_hook": ""}

        activation_codes = []
        for f in flows:
            code = self._gen_enable_code(f.target.id, f.target.type.value, f.target)
            activation_codes.append(code)

        return {
            "next_state_code": "\n".join(activation_codes),
            "pre_hook": ""
        }

    def _gen_end_event(self, event: Element) -> Dict[str, Any]:
        return {
            "next_state_code": "",
            "pre_hook": ""
        }

    def _gen_choreography_task(self, task: ChoreographyTask) -> Dict[str, Any]:
        next_flows = self._get_outgoing_flows(task)
        if not next_flows:
            return {"next_state_code": "", "pre_hook": ""}

        activation_codes = []
        for f in next_flows:
            code = self._gen_enable_code(f.target.id, f.target.type.value, f.target)
            activation_codes.append(code)

        return {
            "next_state_code": "\n".join(activation_codes),
            "pre_hook": ""
        }

    def _gen_exclusive_gateway(self, gateway: ExclusiveGateway) -> Dict[str, Any]:
        out_flows = self._get_outgoing_flows(gateway)

        if len(out_flows) == 1:
            flow = out_flows[0]
            return {
                "branches": [],
                "next_state_code": self._gen_enable_code(flow.target.id, flow.target.type.value, flow.target),
                "is_simple": True
            }

        branches = []
        for flow in out_flows:
            condition = self.analyzer.conditions.get(flow.id, "true")

            condition = condition.replace("state.", "inst.Memory.")

            activation_code = self._gen_enable_code(flow.target.id, flow.target.type.value, flow.target)

            branches.append({
                "target": flow.target.id,
                "condition": condition,
                "activation_code": activation_code
            })

        return {
            "branches": branches,
            "next_state_code": "",
            "is_simple": False
        }

    def _gen_parallel_gateway(self, gateway: ParallelGateway) -> Dict[str, Any]:
        in_flows = self._get_incoming_flows(gateway)
        out_flows = self._get_outgoing_flows(gateway)

        if len(in_flows) <= 1:
            branches = []
            for flow in out_flows:
                code = self._gen_enable_code(flow.target.id, flow.target.type.value, flow.target)
                branches.append({
                    "target": flow.target.id,
                    "activation_code": code
                })

            return {
                "branches": branches,
                "next_state_code": "",
                "is_merge": False
            }

        else:
            checks = []
            for flow in in_flows:
                source_id = self._resolve_target_id(
                    flow.source.type.value if hasattr(flow.source, 'type') else None,
                    flow.source
                )

                container = "Messages" if source_id.startswith("Message") else \
                    "Events" if source_id.startswith("Event") else "Gateways"

                checks.append(f'\tif inst.{container}["{source_id}"].State != COMPLETED {{ return nil }}')

            activations = []
            for flow in out_flows:
                code = self._gen_enable_code(flow.target.id, flow.target.type.value, flow.target)
                activations.append(code)

            full_logic = "\n".join(checks) + "\n" + "\n".join(activations)

            return {
                "branches": [{"activation_code": full_logic}],
                "next_state_code": "",
                "is_merge": True
            }
