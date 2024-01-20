# client.py
# : Simple client code for TCP communication
# : Send 2 nums, and get sum of 2 nums
# -----------------------------------------------

import socket
import struct
import sys

class Message:
    def __init__(self, number1, number2):
        self.number1 = number1
        self.number2 = number2

def tcp_call(server_ip, server_port):
    # Create socket
    client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    # Connect to the server
    server_addr = (server_ip, server_port)
    client_socket.connect(server_addr)
    print(f"Connected to server on port {server_port}")

    # Prepare message
    client_message = Message(1, 2)

    # Send message to server
    message_data = struct.pack('ii', client_message.number1, client_message.number2)
    client_socket.sendall(message_data)
    print(f"Numbers sent to server: {client_message.number1} and {client_message.number2}")

    # Receive result from server
    result_data = client_socket.recv(4)
    sum_result = struct.unpack('i', result_data)[0]
    print(f"Raw Sum received from server: {result_data}")
    print(f"Sum received from server: {sum_result}")

    # Close socket
    client_socket.close()

