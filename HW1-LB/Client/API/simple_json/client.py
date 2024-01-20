# client.py
# : Simple client code for API communication
# -----------------------------------------------

import requests

def api_call(server_ip, server_port):
    # API 서버의 주소
    api_url = f"http://{server_ip}:{server_port}"

    # GET 요청 보내기
    response = requests.get(api_url, stream=True)
    print("API Response:")
    print(response.json())  # JSON 응답 출력
