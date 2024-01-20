# server.py
# : Simple client code for UDP communication
# : Receive 2 nums, and send sum of 2 nums
# -----------------------------------------------

import socket
import struct
import threading

class Message:
    def __init__(self, number1, number2):
        self.number1 = number1
        self.number2 = number2

def handle_client(data, client_address, udp_socket):
    # Unpack received data into a Message object
    received_message = struct.unpack('ii', data)
    print(f"Received numbers from {client_address}: {received_message}")

    # Calculate the sum
    sum_result = received_message[0] + received_message[1]

    # Send the result back to the client
    result_data = struct.pack('i', sum_result)
    udp_socket.sendto(result_data, client_address)
    print(f"Sum sent to {client_address}: {sum_result}")

def client_handler(udp_socket):
    while True:
        data, client_address = udp_socket.recvfrom(1024)
        print(f"Data received from {client_address}: {data}")

        handle_client(data, client_address, udp_socket)

def udp_server(ip, port):
    udp_socket = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    udp_socket.bind((ip, port))

    print(f"UDP Server listening on port {port}...")

    # Create a thread for handling clients
    handler_thread = threading.Thread(target=client_handler, args=(udp_socket,))
    handler_thread.start()
