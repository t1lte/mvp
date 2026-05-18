import re
from typing import Dict, List, Optional, Any
from src.bpmn.graph import Choreography, Message, SequenceFlow

BPMN_TO_GO_TYPES = {
    "boolean": "bool",
    "int": "int",
    "number": "float64",
    "string": "string",
}

class SemanticAnalyzer:
    def __init__(self):
        self.parameters: Dict[str, Dict[str, Any]] = {}
        self.conditions: Dict[str, str] = {}

    def analyze(self, choreography: Choreography):
        self._extract_parameters(choreography)
        self._extract_conditions(choreography)
        return self.parameters, self.conditions

    def _extract_parameters(self, choreography: Choreography):
        for msg_id, message in choreography.messages.items():
            try:
                import json
                doc = json.loads(message.documentation)
                properties = doc.get("properties", {})
                required_list = doc.get("required", [])

                for param_name, param_def in properties.items():
                    bpmn_type = param_def.get("type", "string")
                    go_type = BPMN_TO_GO_TYPES.get(bpmn_type, "string")
                    is_required = param_name in required_list
                    scoped_name = f"{msg_id}_{param_name}"

                    self.parameters[scoped_name] = {
                        "original_name": param_name,
                        "go_type": go_type,
                        "required": is_required,
                        "source_message_id": msg_id
                    }
            except Exception:
                pass

    def _extract_conditions(self, choreography: Choreography):
        for flow in choreography.edges:
            if not isinstance(flow, SequenceFlow):
                continue
            name = flow.name.strip()

            if not name:
                continue

            name = name.replace("&#34;", '"')
            name = name.replace("&#39;", "'")
            name = name.replace("&amp;", "&")

            if " AND " in name.upper() or " OR " in name.upper():
                go_condition = self._parse_complex_condition(name)
                if go_condition:
                    self.conditions[flow.id] = go_condition
                continue

            match = re.match(r"^\s*(\w+)\s*(==|!=|>=|<=|>|<)\s*(.+?)\s*$", name)
            if match:
                var_name, operator, value_str = match.groups()
                found_scoped_name = None
                for scoped in self.parameters:
                    if scoped.endswith(f"_{var_name}"):
                        found_scoped_name = scoped
                        break
                if found_scoped_name:
                    param_info = self.parameters[found_scoped_name]
                    go_field_name = self._to_go_field_name(found_scoped_name)
                    go_value = self._format_value(value_str, param_info["go_type"])
                    self.conditions[flow.id] = f"state.{go_field_name} {operator} {go_value}"

    def _to_go_field_name(self, scoped_name: str) -> str:
        parts = scoped_name.split("_")
        return "".join(part.capitalize() for part in parts)

    def _format_value(self, value_str: str, go_type: str) -> str:
        val = value_str.strip().strip("\"'")
        if go_type == "bool":
            return "true" if val.lower() == "true" else "false"
        elif go_type == "string":
            return f'"{val}"'
        else:
            return val

    def _parse_complex_condition(self, condition_str: str) -> Optional[str]:
        if not condition_str:
            return None

        normalized = re.sub(r'\s+(and|or)\s+', lambda m: f' {m.group(1).upper()} ', condition_str, flags=re.IGNORECASE)

        if ' AND ' in normalized:
            operator = ' && '
            parts = re.split(r'\s+AND\s+', normalized, flags=re.IGNORECASE)
        elif ' OR ' in normalized:
            operator = ' || '
            parts = re.split(r'\s+OR\s+', normalized, flags=re.IGNORECASE)
        else:
            return None

        go_parts = []
        for part in parts:
            part = part.strip()
            match = re.match(r"^\s*(\w+)\s*(==|!=|>=|<=|>|<)\s*(.+?)\s*$", part)
            if match:
                var_name, op, value_str = match.groups()
                found_scoped_name = None
                for scoped in self.parameters:
                    if scoped.lower().endswith(f"_{var_name.lower()}"):
                        found_scoped_name = scoped
                        break

                if found_scoped_name:
                    param_info = self.parameters[found_scoped_name]
                    go_field = self._to_go_field_name(found_scoped_name)
                    go_value = self._format_value(value_str,
                                                  param_info["go_type"])
                    go_parts.append(f"state.{go_field} {op} {go_value}")
                else:
                    return None
            else:
                return None

        return operator.join(go_parts)

    def get_go_struct_fields(self) -> str:
        lines = []
        for scoped_name, info in self.parameters.items():
            go_field = self._to_go_field_name(scoped_name)
            go_type = info["go_type"]
            lines.append(f'\t{go_field} {go_type} `json:"{scoped_name}"`')
        return "\n".join(lines)

    def get_validation_code(self) -> str:
        checks = []
        for scoped_name, info in self.parameters.items():
            if info["required"]:
                go_name = self._to_go_field_name(scoped_name)
                go_type = info["go_type"]
                if go_type == "string":
                    checks.append(f'\tif state.{go_name} == "" {{ return fmt.Errorf("missing required parameter: {scoped_name}") }}')
                elif go_type in ["int", "float64"]:
                    checks.append(f'\tif state.{go_name} == 0 {{ return fmt.Errorf("missing required parameter: {scoped_name}") }}')
        return "\n".join(checks)