# server_list.py
# : Define a server list common to all files
# : Consist of Server List & Server List manage methods

# : server list example (comment) start -------------
'''
data_dict = {
    ("TCP", 80): {
        "ips": [
            {"ip": "192.168.0.1", "user_count": 1, "unhealthy_count": 0},
            {"ip": "192.168.0.2", "user_count": 2, "unhealthy_count": 2}
        ],
        "fd": ,
        "reference_count": 2
        "round_robin": 0
    },
    ("UDP", 53): {
        "ips": [
            {"ip": "192.168.0.3", "user_count": 2, "unhealthy_count": 0},
            {"ip": "192.168.0.4", "user_count": 5, "unhealthy_count": 1}
        ],
        "fd": ,
        "reference_count": 2
        "round_robin": 0
    },
    ...
}

'''
# : server list example (comment) end ---------------
# ----------------------------------------------------------

import random
import threading

from LB.forwarding.forwarding import remove_port_from_LB, open_port_from_LB
from LB.log.logger import Logger

server_list = {}
Log = Logger()


# Get whole Server list dictionary
def get_server_list():
    return server_list


# Check if input server is in Server list dictionary
def is_server_living(protocol, port, ip):
    key = (protocol.lower(), port)

    if key in server_list:
        protocol_dict = server_list[key]

        if "ips" in protocol_dict:
            ips_list = protocol_dict["ips"]

            for entry in ips_list:
                if entry["ip"] == ip:
                    return True
    return False


# Get unhealthy count of server with specific port, protocol, IP in Server list dictionary
def get_unhealthy_count(protocol, port, ip):
    key = (protocol.lower(), port)

    if key in server_list:
        protocol_dict = server_list[key]

        if "ips" in protocol_dict:
            ips_list = protocol_dict["ips"]

            for entry in ips_list:
                if entry["ip"] == ip:
                    return entry["unhealthy_count"]
    return None

def unhealthy_count_up(protocol, port, ip):
    key = (protocol.lower(), port)

    if key in server_list:
        protocol_dict = server_list[key]

        if "ips" in protocol_dict:
            ips_list = protocol_dict["ips"]

            for entry in ips_list:
                if entry["ip"] == ip:
                    unhealthy_count = entry["unhealthy_count"]
                    entry["unhealthy_count"] = unhealthy_count + 1
                return True
    return False


# Get {fd} info from Server list dictionary
def get_fd(protocol, port):
    key = (protocol.lower(), port)
    if key in server_list:
        return server_list[key].get("fd")
    else:
        return None


# Get {reference count} info from Server list dictionary
def get_reference_count(protocol, port):
    key = (protocol.lower(), port)

    if key in server_list:
        protocol_dict = server_list[key]
        return protocol_dict["reference_count"]

    return None


# Get {user count} info from Server list dictionary
def get_user_count(protocol, port, ip):
    key = (protocol.lower(), port)

    if key in server_list:
        protocol_dict = server_list[key]

        if "ips" in protocol_dict:
            ips_list = protocol_dict["ips"]

            for entry in ips_list:
                if entry["ip"] == ip:
                    return entry["user_count"]

    return None


# Set {user count} UP from Server list dictionary
def user_count_up(protocol, port, ip):
    key = (protocol.lower(), port)

    if key in server_list:
        protocol_data = server_list[key]

        if "ips" in protocol_data:
            ips_list = protocol_data["ips"]

            for entry in ips_list:
                if entry["ip"] == ip:
                    # update user_count
                    count = entry["user_count"]
                    entry["user_count"] = count + 1
                    return True
    return False


# Set {reference count} UP from Server list dictionary
def reference_count_up(protocol, port):
    key = (protocol.lower(), port)

    if key in server_list:
        protocol_dict = server_list[key]
        reference_count = protocol_dict["reference_count"]
        protocol_dict["reference_count"] = reference_count + 1
    return None


# Set {user count} DOWN from Server list dictionary
def user_count_down(protocol, port, ip):
    key = (protocol.lower(), port)

    if key in server_list:
        protocol_data = server_list[key]

        if "ips" in protocol_data:
            ips_list = protocol_data["ips"]

            for entry in ips_list:
                if entry["ip"] == ip:
                    # update user_count
                    count = entry["user_count"]
                    entry["user_count"] = count - 1
                    return True
    return False


# Set {reference count} UP from Server list dictionary
def reference_count_down(protocol, port):
    key = (protocol.lower(), port)

    if key in server_list:
        protocol_dict = server_list[key]
        reference_count = protocol_dict["reference_count"]
        protocol_dict["reference_count"] = reference_count - 1
        if (reference_count - 1) == 0:
            delete_entry_fd(protocol, port)
    return None


# Update new entry(with new server_pyte=PROTOCOL, PORT) into Server list dictionary
def add_entry_fd(protocol, port, fd):
    if can_add_entry(protocol, port) == False:
        return "cannot add entry"
    key = (protocol.lower(), port)

    if key in server_list:
        protocol_data = server_list[key]
        if "fd" in protocol_data:
            return f"you already has fd [{protocol}] :{port}"
        else:
            protocol_data["fd"] = fd
            return None

    else:
        server_list[key] = {"fd": fd}
        return None


# Update new IP into Server list dictionary
def add_entry_ip(protocol, port, ip):
    if can_add_entry(protocol, port) == False:
        return False
    key = (protocol.lower(), port)
    print(server_list)

    if key in server_list:
        protocol_data = server_list[key]

        if "ips" in protocol_data:
            ips_list = protocol_data["ips"]

            for entry in ips_list:
                if entry["ip"] == ip:
                    return f"You already registered [{protocol}] {ip}:{port}"

            # new IP entry
            ips_list.append({"ip": ip, "user_count": 0, "unhealthy_count": 0})
            reference_count = protocol_data["reference_count"]
            protocol_data["reference_count"] = reference_count + 1

            Log.info(f"new Server [{protocol}] {ip}:{port} successfully registered")
            return None
    else:
        # 키가 존재하지 않으면 새로운 키와 값을 추가
        # LB입장에서 처음 보는 procotol, port조합에 대해 -> (1) fd open
        # open fd first
        open_port_from_LB(protocol.lower(), port)  # server_list에 fd 추가  (open_forwarding_server 함수 참고)
        # LB입장에서 처음 보는 procotol, port조합에 대해 -> (2) user_count , ip, round_robin, reference count 설정
        server_list[key]["ips"] = [({"ip": ip, "user_count": 0, "unhealthy_count": 0})]
        server_list[key]["reference_count"] = 1
        server_list[key]["round_robin"] = 0
        Log.info(server_list)
        return None


# Delete entry(with server_pyte=PROTOCOL, PORT) from Server list dictionary
def delete_entry_fd(protocol, port):
    key = (protocol.lower(), port)

    if key in server_list:
        protocol_data = server_list[key]

        if "ips" in protocol_data:
            ips_list = protocol_data["ips"]
            if len(ips_list) != 0:
                ips = []
                for ip in ips_list:
                    ips.append(ip)
                return f"ip {ips} is registerd, so you cannot delete entry [{protocol}] :{port} and fd"
            else:
                return remove_port_from_LB(protocol, port)
    else:
        return remove_port_from_LB(protocol, port)


# Delete IP from Server list dictionary
def delete_entry_ip(protocol, port, ip):
    key = (protocol.lower(), port)

    if key in server_list:
        protocol_data = server_list[key]

        if "ips" in protocol_data:
            ips_list = protocol_data["ips"]

            for index, entry in enumerate(ips_list):
                if entry["ip"] == ip:
                    del server_list[(protocol, port)]["ips"][index]
                    Log.info(f"Data [{protocol}] {ip}:{port} removed Success..")
                    reference_count_down(protocol, port)
                    # delete fd if ips_list is empty
                    if len(ips_list) == 0:
                        delete_entry_fd(protocol, port)
                    Log.info(server_list)

                    return None
            Log.error(f"cannot delete ip {ip} NOT in server list")
            return f"cannot delete ip {ip} NOT in server list"
    else:
        Log.error(f"cannot delete ip: [{protocol}] :{port} NOT in server list")
        return f"cannot delete ip: [{protocol}] :{port} NOT in server list"


# Check if a server can be deleted from Server list dictionary
def can_delete_entry_ip(protocol, port, ip):
    key = (protocol.lower(), port)

    if key in server_list:
        protocol_data = server_list[key]

        if "ips" in protocol_data:
            ips_list = protocol_data["ips"]

            for index, entry in enumerate(ips_list):
                if entry["ip"] == ip:
                    delete_entry_ip_thread = threading.Thread(target=delete_entry_ip, args=(protocol, port, ip))
                    delete_entry_ip_thread.start()
                    return None
            Log.error(f"cannot delete ip {ip} NOT in server list")
            return f"cannot delete ip {ip} NOT in server list"
    else:
        Log.error(f"cannot delete ip: [{protocol}] :{port} NOT in server list")
        return f"cannot delete ip: [{protocol}] :{port} NOT in server list"


# Check if a server can be added into Server list dictionary
def can_add_entry(target_protocol, target_port):
    # (1) If (protocol,port) already exists in the existing server list                 ->  Can be added
    # (2) If the protocol exists but the port does not in the existing server list      ->  Can be added
    # (3) If the port exists but the protocol does not in the existing server list      ->  ERROR
    # (4) If neither the port nor the protocol exists in the existing server list       ->  Can be added

    target_protocol = target_protocol.lower()

    # Check the protocol associated with the port
    matching_protocols = []

    for key, protocol_dict in server_list.items():
        current_port, current_protocol = key
        if current_port == target_port:
            matching_protocols.append(current_protocol)

    for protocol in matching_protocols:
        if protocol != target_protocol:
            Log.error(f"port [{target_port}] already registered at [{target_protocol}] :{target_port}")
            return f"port [{target_port}] already registered at [{target_protocol}] :{target_port}"
    return None
