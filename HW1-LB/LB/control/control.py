# control.py
# [1] control server
# : Consist of control server code (on [LB]:8080)
# : Manage 'register/unregister' message

# [2] healthz client
# : Consist of health check client code
# : Manage 'health check' message
# -------------------------------------------

import os
import time
import socket
import threading
import schedule

from LB.log.logger import Logger
from LB.message.message import HealthCheckMessage, ControlMessage, ControlMessageResponse
from LB.table.server_list import delete_entry_ip, add_entry_ip, can_delete_entry_ip, is_server_living, \
    get_unhealthy_count, unhealthy_count_up

Log = Logger()


# Handle requests(register/unregister) from Control Channel client
def handle_client(client_socket):
    client_address = client_socket.getpeername()
    client_ip = client_address[0]

    while True:
        try:
            data = client_socket.recv(1024)
            if not data:
                break

            # Parsing JSON to Dictionary
            message = ControlMessage.from_json(data.decode('utf-8'))
            Log.info(f"message from control client {client_ip}: {message.cmd}, {message.protocol}, {message.port}")

            response = process_command(message, client_ip)
            Log.info(f"response to control client {client_ip}: {response.to_json().encode('utf-8')}")
            client_socket.sendall(response.to_json().encode('utf-8'))

            if message.cmd == "unregister" and response.ack == "Successful":
                # Delete specific server from server table
                delete_entry_ip(message.protocol, message.port, message.ip)

        except Exception as e:
            print(f"Error: {e}")
            Log.error(f"{e}")
            print("err in handle_client")
            break
    client_socket.close()


# Send TCP Message(ex: Control message response) to Specific Server(ip, port)
def send_tcp_message(message, ip, port, timeout=100):
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.settimeout(timeout)
        try:
            s.connect((ip, port))
            s.sendall(message.encode('utf-8'))
            data = s.recv(1024)
            Log.info(f"data")
            return data.decode('utf-8')
        except socket.timeout:
            Log.error(f"Timeout occurred while sending TCP message to {ip}:{port}")
            return None
        except Exception as e:
            Log.error(f"An error occurred: {e}")
            # Shut down port sending tcp message
            return Exception


# Send health check message to the server (every 1 min.) & delete if from Server List if it is unhealthy
def health_check_process(server_ip, protocol, healthz_port, server_port):
    global job

    def health_check():
        if not is_server_living(protocol, server_port, server_ip):
            raise schedule.CancelJob

        healthcheck_message = HealthCheckMessage(cmd="hello")
        tcp_message = healthcheck_message.to_json()

        health_check_bound = int(os.environ.get('HEALTH_CHECK_BOUND', '5'))

        response = send_tcp_message(tcp_message, server_ip, healthz_port, timeout=10)
        Log.info(f"healthz response from [{protocol}] {server_ip}:{healthz_port} {response}")

        response_message = HealthCheckMessage.from_json(response).cmd

        # Log.info(f"response:{HealthCheckMessage}")
        # Log.info(f"HealthCheckMessage.from_json(response).cmd:{response_message}")

        unhealthy_count = get_unhealthy_count(protocol, server_port, server_ip)
        if response_message != "hello":
            unhealthy_count += 1
            unhealthy_count_up(protocol, server_port, server_ip)
            Log.info(f"unhealthy_check_count UP: {unhealthy_count}")
        if unhealthy_count >= health_check_bound:
            # Delete IP from Server List
            Log.warning(
                f"[{protocol}] {server_ip}:{server_port} unhealthy_check_count {unhealthy_count} >= UNHEALTHY_BOUND")
            delete_entry_ip(protocol, server_port, server_ip)
            # Exception Occurred on Job
            raise schedule.CancelJob

    # Event Registration : Run health_Check() func every 1 min.
    job = schedule.every(1).minutes.do(health_check)
    schedule.every(1).minutes.do(health_check)
    try:
        while True:
            # Check registered events & run
            schedule.run_pending()
            time.sleep(1)
    except Exception as e:
    #     Log.info("Job canceled successfully.")
    #     if is_server_living(protocol, server_port, server_ip):
    #         delete_entry_ip(protocol, server_port, server_ip)
        if isinstance(e, schedule.CancelJob):
            Log.info("Job canceled successfully.")
            if is_server_living(protocol, server_port, server_ip):
                delete_entry_ip(protocol, server_port, server_ip)
        else:
            Log.info("Exception..")
            Log.error(e)
            Log.info("Job canceled successfully.")
            if is_server_living(protocol, server_port, server_ip):
                delete_entry_ip(protocol, server_port, server_ip)

# Handle Control message(register/unregister) received from Control Channel client
def process_command(message, client_ip):
    global job
    if message.cmd is None or message.protocol is None or message.port is None:
        return ControlMessageResponse("failed", "Insufficient argument")

    if message.cmd == "register":
        err_msg = add_entry_ip(protocol=message.protocol, port=message.port, ip=client_ip)
        if err_msg is not None:
            return ControlMessageResponse("failed", err_msg)
        else:
            # Start the health check process for the registered client
            health_check_port = int(os.environ.get('HEALTH_CHECK_PORT', '8080'))

            client_process_thread = threading.Thread(target=health_check_process, args=(
                client_ip, message.protocol, health_check_port, message.port))
            client_process_thread.start()

            # Respond to the client indicating successful registration
            return ControlMessageResponse("successful", "")

    elif message.cmd == "unregister":
        # check if the server can be unregistered
        err_msg = can_delete_entry_ip(protocol=message.protocol, port=message.port, ip=client_ip)
        if err_msg is not None:
            return ControlMessageResponse("failed", err_msg)
        else:
            return ControlMessageResponse("Successful", "")

    else:
        return ControlMessageResponse("failed", "Unsupported command")


# Open control server & waiting for clients
def start_server(port):
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server_socket.bind(('0.0.0.0', port))
    server_socket.listen(5)

    Log.info(f"Control Server listening on :{port}")

    while True:
        client_socket, addr = server_socket.accept()
        Log.info(f"Accepted connection from {addr}")

        client_process_thread = threading.Thread(target=handle_client, args=(client_socket,))
        client_process_thread.start()
