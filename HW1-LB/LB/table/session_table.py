# session_table.py
# : Define a session table common to all files
# : Consist of Session Table & Session Table manage methods

# : session table example (comment) start ------------
'''
data_dict = {
    "3332": {                     # client_fd
        "target_protocol": "TCP",
        "target_port": 8080,
        "forwarding_ip": "10.10.0.1"
    },
    "3032": {                     # client_fd
        "target_protocol": "TCP",
        "target_port": 8080,
        "forwarding_ip": "10.10.0.1"
    },
    ...
}
'''
# : session table example (comment) end --------------
# ------------------------------------------------------------

import os

from LB.forwarding.forwarding_rule import select_least_count_ip, random_choice, round_robin
from LB.log.logger import Logger

session_table = {}
prev_index = 0
Log = Logger()


# Get {forwarding ip} info from Session table dictionary
def get_forwarding_ip(protocol, port):
    key = protocol.lower()

    # return existing forwarding ip
    if key in session_table:
        return session_table[key]["forwarding_ip"]

    # return new forwarding ip (by routing rule)
    routing_rule = os.environ.get('ROUTING_RULE', 'USER_COUNT')
    if routing_rule.lower() == "round_robin":
        ip = round_robin(protocol, port)
    elif routing_rule.lower() == "random":
        ip = random_choice(protocol, port)
    else:
        ip = select_least_count_ip(protocol, port)
    return ip


# Update new entry into Session table dictionary
def add_session_table_entry(client_fd, target_protocol, target_port):
    key = client_fd
    target_protocol = target_protocol.lower()

    if key in session_table:
        return False
    else:
        forwarding_ip = get_forwarding_ip(target_protocol, target_port)
        session_table[key] = {"target_protocol": target_protocol, "target_port": target_port,
                              "forwarding_ip": forwarding_ip}
        return forwarding_ip


# Delete entry from Session table dictionary
def delete_session_table_entry(client_fd):
    key = client_fd

    # key가 dictionary에 있는지 확인
    if key in session_table:
        del session_table[key]
        Log.info(f"Data [{client_fd}]  removed Success..")
        return True

    else:
        return False
