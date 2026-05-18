from jinja2 import Environment, FileSystemLoader, Template
from typing import Dict, Any, Optional
from pathlib import Path


class FunctionBuilder:

    def __init__(self, template_dir: Optional[str] = None):
        if template_dir is None:
            template_dir = Path(__file__).parent / "templates"

        self.env = Environment(
            loader=FileSystemLoader(template_dir),
            trim_blocks=True,
            lstrip_blocks=True
        )

    def build_message_send(self,
                           message_id: str,
                           participant_id: str,
                           parameters: list[dict],
                           validation_code: str = "",
                           next_hook: str = "") -> str:
        template = self.env.get_template("message_send.j2")
        return template.render(
            message_id=message_id,
            participant_id=participant_id,
            parameters=parameters,
            validation_code=validation_code,
            next_hook=next_hook
        )

    def build_message_confirm(self,
                              message_id: str,
                              participant_id: str,
                              next_state_code: str,
                              next_hook: str = "") -> str:
        template = self.env.get_template("message_confirm.j2")
        return template.render(
            message_id=message_id,
            participant_id=participant_id,
            next_state_code=next_state_code,
            next_hook=next_hook
        )

    def build_gateway_split(self,
                            gateway_id: str,
                            gateway_type: str,
                            branches: list[dict],
                            next_hook: str = "") -> str:

        template = self.env.get_template("gateway_split.j2")
        return template.render(
            gateway_id=gateway_id,
            gateway_type=gateway_type,
            branches=branches,
            next_hook=next_hook
        )

    def build_gateway_merge(self,
                            gateway_id: str,
                            gateway_type: str,
                            incoming_count: int,
                            next_state_code: str,
                            next_hook: str = "") -> str:

        template = self.env.get_template("gateway_merge.j2")
        return template.render(
            gateway_id=gateway_id,
            gateway_type=gateway_type,
            incoming_count=incoming_count,
            next_state_code=next_state_code,
            next_hook=next_hook
        )

    def build_event_handler(self,
                            event_id: str,
                            event_type: str,
                            next_state_code: str,
                            next_hook: str = "") -> str:

        template = self.env.get_template("event_handler.j2")
        return template.render(
            event_id=event_id,
            event_type=event_type,
            next_state_code=next_state_code,
            next_hook=next_hook
        )