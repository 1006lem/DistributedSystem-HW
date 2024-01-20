# server.py
# : Simple client code for TCP communication
# : Receive 2 nums, and send sum of 2 nums
# -----------------------------------------------

import socket
import struct
import threading

class Message:
    def __init__(self, number1, number2):
        self.number1 = number1
        self.number2 = number2

def handle_client(client_socket):
    try:
        # Receive message from client
        message_data = client_socket.recv(8)  # Ensure you receive 8 bytes
        if len(message_data) != 8:
            print("Error: Insufficient data received.")
            print(f"len(message_data): {len(message_data)}")
            return

        # Unpack received data into a Message object
        received_message = struct.unpack('ii', message_data)
        print(f"Received numbers from client: {received_message}")

        # Calculate the sum
        sum_result = received_message[0] + received_message[1]

        # Send the result back to the client
        result_data = struct.pack('i', sum_result)
        client_socket.sendall(result_data)
        print(f"Sum sent to client: {sum_result}")

    except Exception as e:
        print(f"Error in handle_client: {e}")

    finally:
        # Close socket
        client_socket.close()

def tcp_server(ip, port):
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server_socket.bind((ip, port))
    server_socket.listen(1)

    print(f"Server listening on port {port}...")

    while True:
        client_socket, addr = server_socket.accept()
        print(f"Connection accepted from {addr}")

        # Create a thread for handling clients
        handler_thread = threading.Thread(target=handle_client, args=(client_socket,))
        handler_thread.start()

