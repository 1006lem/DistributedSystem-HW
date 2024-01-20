import os
from Service.control.control import control

if __name__ == "__main__":
    service_type = os.environ.get('SERVICE_TYPE')
    service_port = int(os.environ.get('SERVICE_PORT'))

    control_server_ip = os.environ.get('CONTROL_SERVER_IP')
    control_server_port = int(os.environ.get('CONTROL_SERVER_PORT', '8080'))
    health_check_port = int(os.environ.get('HEALTH_CHECK_PORT', '8080')) # Service에서 돌아가고 있는 health check server의 포트

    # run control_channel after running Real server
    control(control_server_ip=control_server_ip, control_server_port=control_server_port, health_check_port=health_check_port)

