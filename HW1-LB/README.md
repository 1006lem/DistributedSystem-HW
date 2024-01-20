## L3/L4 Load Balancer

---

### 주요 기능

- **Forwarding**
  - client의 요청을 server에게 릴레이


- Control Channel
  - **Health Check**
    - 서버가 유효한 상태인지 확인

  - **Control Message**
    - server register/unregister

---

### 아키텍처

<p align="center"><img width="600" alt="image" src="https://github.com/1006lem/DistributedSystem-HW/assets/68532437/ddb66790-5a90-4f20-b5f2-60a0576772be"></p>

---



### 실행 방법

### 0. Pre-request
- 도커, 쿠버네티스 설치 
- 도커 컨테이너 이미지 생성 (Dockerfile)
```shell
$ cd pre-request
$ ./install_docker.sh
$ ./install_k8s.sh
$ ./prepare_container_image.sh
```


<br>

---

### 1. example 예제 실행

```shell
kubectl apply -f ./test-yaml/example1.yaml
kubectl apply -f ./test-yaml/example1_client.yaml
```
<br>

---

### 2. 서버 실행
- API 서버
```shell
$ kubectl exec -it api-server-1 -- /bin/sh
  # cd Service/service
  # python main.py
```

```shell
$ kubectl exec -it api-server-1 -- /bin/sh
  # cd Service
  # python main.py
```

- TCP 서버
```shell
$ kubectl exec -it tcp-server-1 -- /bin/sh
  # cd Service/service
  # python main.py
```

```shell
$ kubectl exec -it tcp-server-1 -- /bin/sh
  # cd Service
  # python main.py
```

- UDP 서버 
```shell
$ kubectl exec -it udp-server-1 -- /bin/sh
  # cd Service/service
  # python main.py
```

```shell
$ kubectl exec -it udp-server-1 -- /bin/sh
  # cd Service
  # python main.py
```

<br>

---

### 3. LB 실행
```shell
$ kubectl exec -it load-balancer -- /bin/sh
  # cd LB
  # python main.py
```
<br>

---
### 4. API 서버 등록
- 2에서 실행한 API 서버에서 실행
```shell
$ register tcp [port]
```

<br>

---

### 5. client 실행
- API 클라이언트
```shell
$ kubectl exec -it api-client-1 -- /bin/sh 
  # cd Client
  # python main.py
```

- TCP 클라이언트
```shell
$ kubectl exec -it tco-client-1 -- /bin/sh 
  # cd Client
  # python main.py
```

- UDP 클라이언트
```shell
$ kubectl exec -it udp-client-1 -- /bin/sh 
  # cd Client
  # python main.py
```
<br>
<br>
