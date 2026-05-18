from enum import Enum
from typing import List, Optional, Dict, Any, Union


class NodeType(Enum):
    PARTICIPANT = "participant"
    MESSAGE = "message"
    START_EVENT = "startEvent"
    END_EVENT = "endEvent"
    CHOREOGRAPHY_TASK = "choreographyTask"
    EXCLUSIVE_GATEWAY = "exclusiveGateway"
    PARALLEL_GATEWAY = "parallelGateway"


class EdgeType(Enum):
    MESSAGE_FLOW = "messageFlow"
    SEQUENCE_FLOW = "sequenceFlow"



class Element:
    _type: Union[NodeType, EdgeType]
    _properties: List[str] = ["id", "name", "type"]
    _object_properties: List[str] = []

    def __init__(self, id: str, name: str = ""):
        self._id: str = id
        self._name: str = name
        self._graph_ref: Optional["Choreography"] = None

    @property
    def id(self) -> str:
        return self._id

    @property
    def name(self) -> str:
        return self._name

    @property
    def type(self) -> Union[NodeType, EdgeType]:
        return self._type


    def deferred_init(self, graph: "Choreography") -> None:
        self._graph_ref = graph
        for attr_name in self._object_properties:
            raw_value = getattr(self, f"_{attr_name}", None)
            if raw_value is None:
                continue
            if isinstance(raw_value, list):
                resolved = []
                for item in raw_value:
                    if isinstance(item, dict) and "id" in item:
                        elem = graph.get_element_with_id(item["id"])
                        resolved.append({"id": item["id"], "element": elem})
                    else:
                        resolved.append(item)
                setattr(self, f"_{attr_name}", resolved)
            elif isinstance(raw_value, dict) and "id" in raw_value:
                elem = graph.get_element_with_id(raw_value["id"])
                setattr(self, f"_{attr_name}", {"id": raw_value["id"], "element": elem})

    def _get_linked_element(self, attr_name: str) -> Optional["Element"]:
        value = getattr(self, f"_{attr_name}", None)
        if isinstance(value, dict) and "element" in value:
            return value["element"]
        return None

    def _get_linked_elements(self, attr_name: str) -> List["Element"]:
        value = getattr(self, f"_{attr_name}", None)
        if isinstance(value, list):
            return [
                item["element"] for item in value
                if isinstance(item, dict) and "element" in item and item["element"] is not None
            ]
        return []


class Participant(Element):
    _type = NodeType.PARTICIPANT
    def __init__(self, id: str, name: str = ""):
        super().__init__(id, name)


class Message(Element):
    _type = NodeType.MESSAGE
    _properties = ["id", "name", "type", "documentation"]
    def __init__(self, id: str, name: str = "", documentation: str = "{}"):
        super().__init__(id, name)
        self._documentation: str = documentation

    @property
    def documentation(self) -> str:
        return self._documentation


class StartEvent(Element):
    _type = NodeType.START_EVENT
    _properties = ["id", "name", "type", "outgoing"]
    _object_properties = ["outgoing"]
    def __init__(self, id: str, name: str = "", outgoing: str = ""):
        super().__init__(id, name)
        self._outgoing: Dict[str, Any] = {"id": outgoing, "element": None}

    @property
    def outgoing(self) -> Optional["Element"]:
        return self._get_linked_element("outgoing")


class EndEvent(Element):
    _type = NodeType.END_EVENT
    _properties = ["id", "name", "type", "incoming"]
    _object_properties = ["incoming"]
    def __init__(self, id: str, name: str = "", incoming: str = ""):
        super().__init__(id, name)
        self._incoming: Dict[str, Any] = {"id": incoming, "element": None}

    @property
    def incoming(self) -> Optional["Element"]:
        return self._get_linked_element("incoming")


class ChoreographyTask(Element):
    _type = NodeType.CHOREOGRAPHY_TASK
    _properties = ["id", "name", "type", "incoming", "outgoing", "participants", "init_participant", "message_flows"]
    _object_properties = ["incoming", "outgoing", "participants", "init_participant", "message_flows"]

    def __init__(
            self, id: str, name: str = "", incoming: str = "", outgoing: str = "",
            participants: Optional[List[str]] = None, init_participant: str = "",
            message_flows: Optional[List[str]] = None
    ):
        super().__init__(id, name)
        self._incoming: Dict[str, Any] = {"id": incoming, "element": None}
        self._outgoing: Dict[str, Any] = {"id": outgoing, "element": None}
        self._participants: List[Dict[str, Any]] = [{"id": p, "element": None} for p in (participants or [])]
        self._init_participant: Dict[str, Any] = {"id": init_participant, "element": None}
        self._message_flows: List[Dict[str, Any]] = [{"id": mf, "element": None} for mf in (message_flows or [])]

    @property
    def incoming(self) -> Optional["Element"]:
        return self._get_linked_element("incoming")

    @property
    def outgoing(self) -> Optional["Element"]:
        return self._get_linked_element("outgoing")

    @property
    def participants(self) -> List["Participant"]:
        elems = self._get_linked_elements("participants")
        return [e for e in elems if isinstance(e, Participant)]

    @property
    def init_participant(self) -> Optional["Participant"]:
        elem = self._get_linked_element("init_participant")
        return elem if isinstance(elem, Participant) else None

    @property
    def message_flows(self) -> List["MessageFlow"]:
        elems = self._get_linked_elements("message_flows")
        return [e for e in elems if isinstance(e, MessageFlow)]

    @property
    def init_message_flow(self) -> Optional["MessageFlow"]:
        init_p = self.init_participant
        for mf in self.message_flows:
            if mf.source == init_p:
                return mf
        return None


class ExclusiveGateway(Element):
    _type = NodeType.EXCLUSIVE_GATEWAY
    _properties = ["id", "name", "type", "incomings", "outgoings"]
    _object_properties = ["incomings", "outgoings"]
    def __init__(self, id: str, name: str = "", incomings: Optional[List[str]] = None, outgoings: Optional[List[str]] = None):
        super().__init__(id, name)
        self._incomings: List[Dict[str, Any]] = [{"id": inc, "element": None} for inc in (incomings or [])]
        self._outgoings: List[Dict[str, Any]] = [{"id": out, "element": None} for out in (outgoings or [])]

    @property
    def incomings(self) -> List["Element"]:
        return self._get_linked_elements("incomings")

    @property
    def outgoings(self) -> List["Element"]:
        return self._get_linked_elements("outgoings")


class ParallelGateway(Element):
    _type = NodeType.PARALLEL_GATEWAY
    _properties = ["id", "name", "type", "incomings", "outgoings"]
    _object_properties = ["incomings", "outgoings"]
    def __init__(self, id: str, name: str = "", incomings: Optional[List[str]] = None, outgoings: Optional[List[str]] = None):
        super().__init__(id, name)
        self._incomings: List[Dict[str, Any]] = [{"id": inc, "element": None} for inc in (incomings or [])]
        self._outgoings: List[Dict[str, Any]] = [{"id": out, "element": None} for out in (outgoings or [])]

    @property
    def incomings(self) -> List["Element"]:
        return self._get_linked_elements("incomings")

    @property
    def outgoings(self) -> List["Element"]:
        return self._get_linked_elements("outgoings")




class MessageFlow(Element):
    _type = EdgeType.MESSAGE_FLOW
    _properties = ["id", "name", "type", "source", "target", "message"]
    _object_properties = ["source", "target", "message"]
    def __init__(self, id: str, name: str = "", source: str = "", target: str = "", message: str = ""):
        super().__init__(id, name)
        self._source: Dict[str, Any] = {"id": source, "element": None}
        self._target: Dict[str, Any] = {"id": target, "element": None}
        self._message: Dict[str, Any] = {"id": message, "element": None}

    @property
    def source(self) -> Optional["Element"]:
        return self._get_linked_element("source")

    @property
    def target(self) -> Optional["Element"]:
        return self._get_linked_element("target")

    @property
    def message(self) -> Optional["Message"]:
        elem = self._get_linked_element("message")
        return elem if isinstance(elem, Message) else None


class SequenceFlow(Element):
    _type = EdgeType.SEQUENCE_FLOW
    _properties = ["id", "name", "type", "source", "target"]
    _object_properties = ["source", "target"]
    def __init__(self, id: str, name: str = "", source: str = "", target: str = ""):
        super().__init__(id, name)
        self._source: Dict[str, Any] = {"id": source, "element": None}
        self._target: Dict[str, Any] = {"id": target, "element": None}

    @property
    def source(self) -> Optional["Element"]:
        return self._get_linked_element("source")

    @property
    def target(self) -> Optional["Element"]:
        return self._get_linked_element("target")


class Choreography:
    def __init__(self):
        self.nodes: List[Element] = []
        self.edges: List[Element] = []
        self._id_map: Dict[str, Element] = {}
        self.participants: Dict[str, Participant] = {}
        self.messages: Dict[str, Message] = {}
        self.tasks: Dict[str, ChoreographyTask] = {}
        self.gateways: Dict[str, Union[ExclusiveGateway, ParallelGateway]] = {}
        self.events: Dict[str, Union[StartEvent, EndEvent]] = {}

    def get_element_with_id(self, element_id: str) -> Optional[Element]:
        return self._id_map.get(element_id)

    def add_node(self, node: Element) -> None:
        self.nodes.append(node)
        self._id_map[node.id] = node
        if isinstance(node, Participant):
            self.participants[node.id] = node
        elif isinstance(node, Message):
            self.messages[node.id] = node
        elif isinstance(node, ChoreographyTask):
            self.tasks[node.id] = node
        elif isinstance(node, (ExclusiveGateway, ParallelGateway)):
            self.gateways[node.id] = node
        elif isinstance(node, (StartEvent, EndEvent)):
            self.events[node.id] = node

    def add_edge(self, edge: Element) -> None:
        self.edges.append(edge)
        self._id_map[edge.id] = edge

    def add_participant(self, participant: Participant) -> None:
        self.add_node(participant)

    def add_message(self, message: Message) -> None:
        self.add_node(message)

    def finalize(self) -> None:
        for node in self.nodes:
            node.deferred_init(self)
        for edge in self.edges:
            edge.deferred_init(self)