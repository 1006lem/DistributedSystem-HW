# forwarding.py
# : Forwards client's request to real Server
# : Manage forwarding server
# -----------------------------------------------

import socket
import threading

from LB.log.logger import Logger
from LB.table.session_table import add_session_table_entry, delete_session_table_entry
Log = Logger()
exit_flag = False


def handle_tcp_client(client_socket, port):
    # Receive message from client
    data = b""
    client_socket.settimeout(30)  # set timeout
    try:
        while True:
            chunk = client_socket.recv(1024)
            data += chunk
            if not chunk or len(chunk) < 1024:
                break
    except socket.timeout:
        Log.error("Timeout while receiving data from client ...")

    if not data:
        return

    Log.info(f"Raw TCP Data: {data}")
    # session_table 추가 -> return되는 IP가 있다면 해당 IP로 포워딩
    #       client_socket, TCP, forwarding_ip, port 추가
    forwarding_ip = add_session_table_entry(client_fd=client_socket, target_protocol="TCP", target_port=port)
    # server list 추가 -> ip에 user_count 증가
    from LB.table.server_list import user_count_up
    user_count_up("tcp", port, forwarding_ip)
    return_from_server = send_to_tcp_server(data, forwarding_ip, port)
    if return_from_server is None:
        from LB.table.server_list import delete_entry_ip
        delete_entry_ip("tcp", port, forwarding_ip)

    client_socket.close()
    # session_table 삭제
    delete_session_table_entry(client_fd=client_socket)
    # server list 추가 -> ip에 user_count 감소
    from LB.table.server_list import user_count_down
    user_count_down("tcp", port, forwarding_ip)


def send_to_tcp_server(payload, ip, port):
    try:
        # Create socket
        client_socket2 = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

        # Connect to the service
        server_addr = (ip, port)
        client_socket2.connect(server_addr)
        Log.info(f"Connected to service on port {port}")

        if payload is None:
            client_socket2.close()
        client_socket2.sendall(payload)

        # Receive result from service
        result_data = b""
        while True:
            chunk = client_socket2.recv(1024)
            result_data += chunk
            #if not chunk or b"Connection: close" in result_data:
            if not chunk or len(chunk) < 1024:
                break
        Log.info(f"Result data: {result_data}")

        return result_data
    except Exception as e:
        Log.error(f"Error: {e}")
    finally:
        # Close the socket
        client_socket2.close()


def send_to_udp_server(payload, ip, port):
    try:
        # Create socket
        client_socket2 = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)  # Change to SOCK_DGRAM for UDP
        server_addr = (ip, port)
        client_socket2.sendto(payload, server_addr)
        # Receive result from service
        result_data = client_socket2.recv(1024)
        Log.info(f"Result data: {result_data}")

        return result_data
    except Exception as e:
        Log.error(f"Error: {e}")
    finally:
        # Close the socket
        client_socket2.close()
        pass


def handle_tcp_connections(server_socket, port):
    global exit_flag

    while not exit_flag:
        try:
            tcp_client_socket, addr = server_socket.accept()
            Log.info(f"Connection accepted from TCP client: {addr}")

            # Create a new thread to handle the client
            client_thread = threading.Thread(target=handle_tcp_client, args=(tcp_client_socket, port))
            client_thread.start()
        except socket.error as e:
            if exit_flag:
                break
            Log.error(f"Error accepting connection: {e}")


def handle_udp_client(client_socket, port):
    try:
        message_data, client_addr = client_socket.recvfrom(1024)  # Ensure you receive 8 bytes

        Log.info(f"Raw UDP Data: {message_data}")
        # session_table 추가 -> return되는 IP가 있다면 해당 IP로 포워딩
        #       client_socket, UDP, forwarding_ip, port 추가
        forwarding_ip = add_session_table_entry(client_fd=client_socket, target_protocol="UDP", target_port=port)
        #         user_count_up("udp", port, forwarding_ip)
        from LB.table.server_list import user_count_up
        user_count_up("udp", port, forwarding_ip)
        return_from_server = send_to_udp_server(message_data, forwarding_ip, port)
        if return_from_server is None:
            from LB.table.server_list import delete_entry_ip
            delete_entry_ip("udp", port, forwarding_ip)

        client_socket.sendto(return_from_server, client_addr)
        # session_table 삭제
        delete_session_table_entry(client_fd=client_socket)
        from LB.table.server_list import user_count_down
        user_count_down("udp", port, forwarding_ip)
    except Exception as e:
        Log.error(f"Error: {e}")
    finally:
        pass


def handle_udp_connections(udp_server_socket, port):
    while not exit_flag:
        # Create a new thread to handle the client
        client_thread = threading.Thread(target=handle_udp_client, args=(udp_server_socket, port))
        client_thread.start()


def open_forwarding_server(protocol, port, ip):
    global exit_flag
    exit_flag = False
    Log.info(f"Opening service .. [{protocol}] {ip}:{port}")

    try:
        if protocol.lower() == "tcp":
            tcp_server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            tcp_server_socket.bind((ip, port))
            tcp_server_socket.listen(1)
            # service-list에 추가 (protocol, service port, LB fd)
            from LB.table.server_list import add_entry_fd
            add_entry_fd(protocol=protocol, port=port, fd=tcp_server_socket)
            threading.Thread(target=handle_tcp_connections, args=(tcp_server_socket, port)).start()
            return None

        elif protocol.lower() == "udp":
            udp_server_socket = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
            udp_server_socket.bind((ip, port))

            # service-list에 추가 (protocol, service port, LB fd)
            from LB.table.server_list import add_entry_fd
            add_entry_fd(protocol=protocol, port=port, fd=udp_server_socket)
            threading.Thread(target=handle_udp_connections, args=(udp_server_socket, port)).start()
            return None

    except socket.error as e:
        if e.errno == 13:  # Permission denied
            error_message = "Error: Permission denied. You may not have sufficient privileges to bind to this port."
            return error_message
        elif e.errno == 22:  # Invalid argument
            error_message = "Error: Invalid argument. Please check the provided arguments."
            return error_message
        elif e.errno == 98:  # Address already in use
            error_message = f"Error: Address {ip}:{port} is already in use. Please choose a different address or port."
            return error_message
        elif e.errno == 99:  # Cannot assign requested address
            error_message = "Error: Cannot assign the requested address. Please check the provided IP address."
            return error_message
        elif e.errno == 111:  # Connection refused
            error_message = "Error: Connection refused. Check if the server is running and reachable."
            return error_message
        elif e.errno == 106:  # Transport endpoint is already connected
            error_message = "Error: Transport endpoint is already connected."
            return error_message
        else:
            error_message = f"Error: {e}"
            return error_message


def close_forwarding_server_socket(server_socket):
    global exit_flag
    exit_flag = True

    if server_socket:
        server_socket.close()
        Log.info("Server Socket closed.")


# Open LB's port (start forwarding service)
def open_port_from_LB(protocol, port):
    if protocol.lower() == "udp" or protocol.lower() == "tcp":
        Log.info(f"Trying to open new {protocol} forwarding server on :{port}")
        # Start forwarding service
        err = open_forwarding_server(protocol, port, "0.0.0.0")
        return err

# Close LB's port (End forwarding service)
def remove_port_from_LB(protocol, port):
    # Service-list에서 socket fd 확인 (protocol, port로 조회)
    from LB.table.server_list import get_fd
    tcp_server_socket = get_fd(protocol, port)
    close_forwarding_server_socket(tcp_server_socket)
    Log.info(f"[Success] close LB's [{protocol}] :{port}")
    return False
