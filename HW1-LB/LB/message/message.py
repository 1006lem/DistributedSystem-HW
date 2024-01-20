# message.py
# : Define Messages for Server-Client Communication
# :     [Control message] Communication 1: Server(LB)-Client(Server)
# :     [Health check]    Communication 2: Server(Server)-Client(LB)
# --------------------------------------------------------------------

import json


class ControlMessage:
    def __init__(self, cmd, protocol, port):
        self.cmd = cmd
        self.protocol = protocol
        self.port = port

    def to_json(self):
        return json.dumps(self.__dict__)

    @classmethod
    def from_json(cls, json_str):
        data = json.loads(json_str)
        return cls(**data)


class ControlMessageResponse:
    def __init__(self, ack, msg):
        self.ack = ack
        self.msg = msg

    def to_json(self):
        return json.dumps(self.__dict__)

    @classmethod
    def from_json(cls, json_str):
        data = json.loads(json_str)
        return cls(**data)


class HealthCheckMessage:
    def __init__(self, cmd):
        self.cmd = cmd

    def to_json(self):
        return json.dumps(self.__dict__)

    @classmethod
    def from_json(cls, json_str):
        data = json.loads(json_str)
        return cls(**data)
