import xml.etree.ElementTree as ET
from typing import List, Optional, Dict, Any

from .graph import (
    NodeType, EdgeType,
    Participant, Message,
    StartEvent, EndEvent,
    ChoreographyTask, ExclusiveGateway, ParallelGateway,
    MessageFlow, SequenceFlow,
    Choreography
)

BPMN_NS = "http://www.omg.org/spec/BPMN/20100524/MODEL"
BPMN_TAG = lambda tag: f"{{{BPMN_NS}}}{tag}"


class BPMNParser:

    def __init__(self, xml_content: str):
        self.xml_content = xml_content
        self.choreography = Choreography()

    def parse(self) -> Choreography:
        root = ET.fromstring(self.xml_content)
        choreo_elem = root.find(BPMN_TAG('choreography'))
        if choreo_elem is None:
            choreo_elem = root.find(BPMN_TAG('collaboration'))

        if choreo_elem is None:
            raise ValueError("BPMN file must contain either <choreography> or <collaboration> element")

        self._parse_messages(root)
        self._parse_participants(choreo_elem)
        self._parse_nodes(choreo_elem)
        self._parse_flows(choreo_elem)
        self.choreography.finalize()

        return self.choreography

    def _parse_messages(self, root: ET.Element) -> None:
        for msg_elem in root.findall(BPMN_TAG('message')):
            msg_id = msg_elem.get('id', '')
            name = msg_elem.get('name', '')

            doc_elem = msg_elem.find(BPMN_TAG('documentation'))
            doc_text = doc_elem.text.strip() if doc_elem is not None and doc_elem.text else "{}"

            message = Message(id=msg_id, name=name, documentation=doc_text)
            self.choreography.add_message(message)

    def _parse_participants(self, choreo_elem: ET.Element) -> None:
        for p_elem in choreo_elem.findall(BPMN_TAG('participant')):
            p_id = p_elem.get('id', '')
            p_name = p_elem.get('name', '')
            participant = Participant(id=p_id, name=p_name)
            self.choreography.add_participant(participant)

    def _parse_nodes(self, choreo_elem: ET.Element) -> None:
        for el in choreo_elem.findall(BPMN_TAG('startEvent')):
            outgoing_elems = el.findall(BPMN_TAG('outgoing'))
            outgoing = outgoing_elems[0].text if outgoing_elems else ""
            node = StartEvent(
                id=el.get('id', ''),
                name=el.get('name', ''),
                outgoing=outgoing
            )
            self.choreography.add_node(node)

        for el in choreo_elem.findall(BPMN_TAG('endEvent')):
            incoming_elems = el.findall(BPMN_TAG('incoming'))
            incoming = incoming_elems[0].text if incoming_elems else ""
            node = EndEvent(
                id=el.get('id', ''),
                name=el.get('name', ''),
                incoming=incoming
            )
            self.choreography.add_node(node)

        for el in choreo_elem.findall(BPMN_TAG('choreographyTask')):
            participants = [
                p_ref.text for p_ref in el.findall(BPMN_TAG('participantRef')) if p_ref.text
            ]
            message_flows = [
                mf_ref.text for mf_ref in el.findall(BPMN_TAG('messageFlowRef')) if mf_ref.text
            ]

            incoming_elems = el.findall(BPMN_TAG('incoming'))
            outgoing_elems = el.findall(BPMN_TAG('outgoing'))
            incoming = incoming_elems[0].text if incoming_elems else ""
            outgoing = outgoing_elems[0].text if outgoing_elems else ""

            node = ChoreographyTask(
                id=el.get('id', ''),
                name=el.get('name', ''),
                incoming=incoming,
                outgoing=outgoing,
                participants=participants,
                init_participant=el.get('initiatingParticipantRef', ''),
                message_flows=message_flows
            )
            self.choreography.add_node(node)

        self._parse_gateway(choreo_elem, ExclusiveGateway, 'exclusiveGateway')
        self._parse_gateway(choreo_elem, ParallelGateway, 'parallelGateway')

    def _parse_gateway(self, choreo_elem: ET.Element, gateway_class, tag_name: str) -> None:
        for el in choreo_elem.findall(BPMN_TAG(tag_name)):
            incomings = [i.text for i in el.findall(BPMN_TAG('incoming')) if i.text]
            outgoings = [o.text for o in el.findall(BPMN_TAG('outgoing')) if o.text]
            node = gateway_class(
                id=el.get('id', ''),
                name=el.get('name', ''),
                incomings=incomings,
                outgoings=outgoings
            )
            self.choreography.add_node(node)

    def _parse_flows(self, choreo_elem: ET.Element) -> None:
        for el in choreo_elem.findall(BPMN_TAG('sequenceFlow')):
            flow = SequenceFlow(
                id=el.get('id', ''),
                name=el.get('name', ''),
                source=el.get('sourceRef', ''),
                target=el.get('targetRef', '')
            )
            self.choreography.add_edge(flow)

        for el in choreo_elem.findall(BPMN_TAG('messageFlow')):
            flow = MessageFlow(
                id=el.get('id', ''),
                name=el.get('name', ''),
                source=el.get('sourceRef', ''),
                target=el.get('targetRef', ''),
                message=el.get('messageRef', '')
            )
            self.choreography.add_edge(flow)