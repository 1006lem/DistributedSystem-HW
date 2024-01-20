import os
import API.simple_json.server as api_server
import TCP.add_two_num.server as tcp_server
import UDP.add_two_num.server as udp_server

if __name__ == "__main__":
    service_type = os.environ.get('SERVICE_TYPE')
    service_port = int(os.environ.get('SERVICE_PORT'))

    control_server_ip = os.environ.get('CONTROL_SERVER_IP')
    control_server_port = int(os.environ.get('CONTROL_SERVER_PORT', '8080'))
    health_check_port = int(os.environ.get('HEALTH_CHECK_PORT','8080')) # Service에서 돌아가고 있는 health check server의 포트

    if service_type.lower() == "api":
        if service_port is None:
            service_port = 8081 # 기본값
        api_server.app.run('0.0.0.0', port=service_port, debug=True)

    elif service_type.lower() == "tcp":
        if service_port is None:
            service_port = 8082 # 기본값
        tcp_server.tcp_server('0.0.0.0', port=service_port)

    elif service_type.lower() == "udp":
        if service_port is None:
            service_port = 8083 # 기본값
        udp_server.udp_server('0.0.0.0', port=service_port)

