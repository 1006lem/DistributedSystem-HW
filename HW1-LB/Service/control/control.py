# control.py
# [1] control client
# : Continue to receive input from the user interactive
# : Consist of control client code
# : Send 'register/unregister' message to control server (on [LB]:8080)

# [2] healthz server
# : Consist of health check server code (on [Server]:8080)
# : Manage 'health check' message
# ----------------------------------------------------------

import socket
import threading

from Service.log.logger import Logger
from Service.message.message import ControlMessage, HealthCheckMessage

health_check_thread = None
terminate_health_check = False
Log = Logger()


# Open health check server & wait for a client
def start_healthcheck_server(ip, port):
    global health_check_thread
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server_socket.bind((ip, port))
    server_socket.listen(10)

    Log.info(f"Server listening on {ip}:{port}")

    while not terminate_health_check:
        client_socket, addr = server_socket.accept()
        data = client_socket.recv(1024)

        if not data:
            break

        # Parsing data in JSON format to dictionary
        message = HealthCheckMessage.from_json(data.decode('utf-8'))
        Log.info(f"health check: {message.to_json().encode('utf-8')}")
        client_socket.sendall(message.to_json().encode('utf-8'))


# Make a new thread for health check server
def start_healthcheck_thread(port):
    global health_check_thread, terminate_health_check
    terminate_health_check = False
    health_check_thread = threading.Thread(target=start_healthcheck_server, args=('0.0.0.0', port))
    health_check_thread.start()


# Send TCP Message(ex: Control message request) to Specific Server(ip, port)
def send_tcp_message(message, ip, port):
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        try:
            s.connect((ip, port))
            s.sendall(message.encode('utf-8'))
            data = s.recv(1024)
            return data.decode('utf-8')
        except Exception as e:
            Log.error(f"An error occurred: {e}")
            return None


# Continue to receive input(until 'quit') from the user interactive
def control(control_server_ip, control_server_port, health_check_port):
    global health_check_thread, terminate_health_check
    print("if you want to quit, please enter 'quit'")
    while True:
        try:
            user_input = input("Enter a command (e.g., 'register TCP 8080', 'quit'): ")
        except EOFError:
            print("EOFError: Please provide valid input.")
            continue

        if user_input.lower() == "quit":
            print("Exiting the program.")
            terminate_health_check = True
            if health_check_thread:
                health_check_thread.join()
                health_check_thread = None
            break

        try:
            cmd, protocol, port = user_input.split()
            port = int(port)
        except ValueError:
            print("Invalid input. Please use the format 'register/unregister TCP/UDP(TCP/UDP) port'.")
            continue

        if protocol.lower() != "udp" and protocol.lower() != "tcp":
            print("Invalid input. please use 'TCP' or 'UDP' as a protocol")
            continue

        if cmd.lower() == "register":
            start_healthcheck_thread(health_check_port)
            control_message = ControlMessage(cmd="register", protocol=protocol.lower(), port=port)
            tcp_message = control_message.to_json()
            response = send_tcp_message(tcp_message, control_server_ip, control_server_port)
            Log.info(f"Register Response: {response}")

        elif cmd.lower() == "unregister":
            terminate_health_check = True
            if health_check_thread:
                health_check_thread.join()
                health_check_thread = None
            control_message = ControlMessage(cmd="unregister", protocol=protocol.lower(), port=port)
            tcp_message = control_message.to_json()
            response = send_tcp_message(tcp_message, control_server_ip, control_server_port)
            Log.info(f"Unregister Response: {response}")

        else:
            print("Unsupported command. Please use 'register' or 'unregister'.")

