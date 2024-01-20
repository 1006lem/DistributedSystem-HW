import os

from Client.API.simple_json.client import api_call
from Client.TCP.add_two_num.client import tcp_call
from Client.UDP.add_two_num.client import udp_call

def client_request(server_type, server_ip, server_port):
    if server_type.lower() == "api":
        api_call(server_ip, server_port)
    elif server_type.lower() == "tcp":
        tcp_call(server_ip, server_port)
    elif server_type.lower() == "udp":
        udp_call(server_ip, server_port)
    else:
        "[ERROR] server type should be API/TCP/UDP"

if __name__ == "__main__":
    server_type = os.environ.get('SERVER_TYPE')
    server_ip = os.environ.get('SERVER_IP')
    server_port = int(os.environ.get('SERVER_PORT'))

    client_request(server_type, server_ip, server_port)



