## Simple Distributed System


### 주요 기능

- **primary based local-write**
    - config.json 파일을 통한 replica 정보 관리
    - PM으로 client의 요청 릴레이
    - local 저장소 동기화


- **primary based remote-write**
    - config.json 파일을 통한 replica 정보 관리
    - 모든 서버에서 client 요청 처리
    - client의 요청에 따라 PM 갱신
    - 모든 서버에서 PM 추적
    - local 저장소 동기화


---


### 실행 방법

### 0. Pre-request
- golang 설치
```shell
$ ./pre_requirement.sh
```

<br>

---
### 1. 서버 실행
remote wirte 또는 local write 실행

### 1-1. remote write 실행

#### (1) PM
```shell
cd ./server
go run main.go ../config-file/remote-write/config2.json
```

#### (2) Replica 1
```shell
cd ./server
go run main.go ../config-file/remote-write/config1.json
```

#### (3) Replica 2
```shell
cd ./server
go run main.go ../config-file/remote-write/config3.json
```


<br>

### 1-2. local write 실행


```shell
cd ./server
go run main.go ../config-file/local-write/config2.json
```

#### (2) Replica 1
```shell
cd ./server
go run main.go ../config-file/local-write/config1.json
```

#### (3) Replica 2
```shell
cd ./server
go run main.go ../config-file/local-write/config3.json
```
<br>

---
### 2. Client 실행 

```shell
cd ./client

go run main.go post
go run main.go get

go run main.go patch
go run main.go get 1
go run main.go get 2
go run main.go get 3

go run main.go put
go run main.go get 

go run main.go delete
go run main.go get 
```

<br>
<br>
