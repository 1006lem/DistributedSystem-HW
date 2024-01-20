# client.py
# : Simple client code for UDP communication
# : Send 2 nums, and get sum of 2 nums
# -----------------------------------------------

import socket
import struct

TIMEOUT = 5  # seconds


class Message:
    def __init__(self, number1, number2):
        self.number1 = number1
        self.number2 = number2


def udp_call(server_ip, server_port):
    # Create socket
    client_socket = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)  # Change to SOCK_DGRAM for UDP
    client_socket.settimeout(TIMEOUT)

    # Prepare message
    # client_message = Message(int(sys.argv[1]), int(sys.argv[2]))
    # client_message = Message(1, 2)
    messages = [Message(1, 2), Message(3, 5)]

    for client_message in messages:
        # Send message to server
        message_data = struct.pack('ii', client_message.number1, client_message.number2)
        server_addr = (server_ip, server_port)
        client_socket.sendto(message_data, server_addr)
        print(f"Numbers sent to server: {client_message.number1} and {client_message.number2}")

        try:
            # Receive result from server
            result_data, server_addr = client_socket.recvfrom(4)
            sum_result = struct.unpack('i', result_data)[0]
            print(f"Raw Sum received from server: {result_data}")
            print(f"Sum received from server: {sum_result}")
        except socket.timeout:
            print(f"Timeout reached. Communication with the server has ended.")
            break

    # Close socket
    client_socket.close()
