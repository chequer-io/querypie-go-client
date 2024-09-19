# QueryPie Client for Operation

QueryPie Client for Operation 은 QueryPie 운영을 위해 사용할 수 있는 cli client 를 제공합니다.
Command Line Interface 로 작동하며, 자동화 프로그램을 위해 사용할 수 있으며, 운영자가 터미널에서
정보를 조회하거나 정책 설정을 관리하는 용도로 사용할 수 있습니다.

## 설치

### 1. 설치 파일 다운로드 및 빌드

git repository 를 git-clone 하여, 필요한 파일을 내려받습니다.
```bash
$ git clone https://github.com/chequer-io/querypie-go-client.git
```

빌드를 수행하기 위해서는, Golang 최신 버전이 필요합니다. 빌드와 작동을 테스트한 버전은 go1.23.0 입니다.
```bash
$ go version  
go1.23.0 darwin/arm64
```

다운로드 받은 디렉토리로 이동하여, 빌드를 수행합니다.
```bash
$ cd querypie-go-client
$ make build
```

빌드가 성공하면, `qpc` 실행파일이 생성됩니다.

### 2. YAML 설정파일 작성

`.querypie-client.yaml`이라는 설정파일이 실행시점의 current working directory 에 놓여 있어야 합니다.
`qpc` 실행파일과 동일한 디렉토리에 `.querypie-client.yaml` 파일을 생성하고, `./qpc` 처럼 실행하는 방식을
이용하면 편리합니다.

#### QueryPie 서버 주소와 API Key 설정

QueryPie Client for Operation 을 사용하기 위해서는, QueryPie 서버의 주소와 API Token 을 설정해야 합니다.
QueryPie 서버에서 API Token 을 발급받은 후, UUID 형태의 Token 을 준비해 주세요.

아래 예시와 같이 QueryPie 서버의 주소와 API Token 을 설정합니다. yaml 파일에는 하나 이상의 QueryPie 서버를 
설정해 둘 수 있습니다. 연결에 사용할 서버는 `default: true` 로 설정합니다.

#### Local database 설정

`sqlite3-data-source` 설정값을 입력하여야 합니다. Local database 는 sqlite3 를 사용하며, Local database 는
QueryPie 서버로부터 내려받은 정보를 저장하는 용도로 사용됩니다.

Local database 는 sqlite3 를 사용하며, 예시 설정에서는 `var` 디렉토리에 저장됩니다.

Local database 는 QueryPie 서버로부터 내려받은 정보를 저장하는 용도이며, 이용자가 삭제하여도 무방합니다.
Local database 를 초기화하려면, `var/resources.db` 파일을 삭제하면 됩니다. 

`qpc`는 실행 시점에, 설정된 경로에 Local database 를 생성하거나, 읽어들입니다. 설정값을 변경하여, 다른 경로에 저장할 수 있습니다.

#### .querypie-client.yaml 파일 예시

```.querypie-client.yaml
# My config of qpc
querypie-servers:
  - name: "internal.querypie.io"
    url: https://internal.querypie.io/
    token: ap5bbb2a2f-****-****-****-************
    default: true
sqlite3-data-source: ./var/resources.db
```

### 3. QueryPie 서버 연결을 확인하여 설정 완료 

아래와 같이 설정된 QueryPie 서버의 연결을 확인합니다. STATUS OK 로 표시되면, API Token 을 이용해, 서버로부터
정상적인 API 응답을 얻었다는 것을 의미합니다.
```bash
$ ./qpc config querypie
NAME                            BASE_URL                              ACCESS_TOKEN                            STATUS
internal.querypie.io[*]         https://internal.querypie.io          ap5bbb2a2f-****-****-****-************  OK
perf.querypie-stage.com         https://perf.querypie-stage.com       your_access_*****_****                  FAIL
$
```

## 사용법

QueryPie Client for Operation 는 QueryPie 서버로부터 정보를 조회하거나, 정책을 설정하는 명령어를
제공합니다. 접근제어 정책을 설정하는 경우, 관련된 정보를 미리 내려받아 local database 에 저장한 이후,
local database 에서 정보를 조회하여, 정책을 설정하는 방식으로 작동합니다.

### 1. 서버로부터 이용자, DAC 정보 내려받기

서버로부터 이용자, DAC 정보를 내려받아 Local database 에 저장하는 기능은 멱등성(Idempotence)을 갖습니다.
명령을 여러 차례 실행하여, 갱신된 정보를 내려받을 수 있습니다.

```bash
$ ./qpc user fetch-all # 서버로부터 이용자 정보를 내려받아 Local database 에 저장합니다.
$ ./qpc dac fetch-all privileges # 서버로부터 DAC Privilege 정보를 내려받아 저장합니다.
$ ./qpc dac fetch-all detailed-connections # 서버로부터 DAC Connection 상세정보를 내려받아 저장합니다.
```

```bash
$ ./qpc user ls # Local database 에 저장된 이용자 정보를 조회합니다.
$ ./qpc dac ls privileges # Local database 에 저장된 DAC Privilege 정보를 조회합니다.
$ ./qpc dac ls detailed-connections # DAC Connection 상세정보를 조회합니다.
```

서버로부터 DAC Connection 상세정보를 내려받아 저장한 후에는, DAC Cluster 정보를 조회할 수 있습니다.
DAC 의 연결대상인 Connection 는 1개 이상의 Cluster 로 구성되어 있습니다. DAC Access Control 권한 부여는
기본적으로 Cluster 단위로 이루어집니다. 다만, Connection 내에 Cluster 가 1개만 설정된 경우, Connection 식별자로
권한을 부여하는 편의 기능을 `qpc dac grant` 명령어에서 제공합니다.

```bash
$ ./qpc dac ls clusters # DAC Cluster 정보를 조회합니다.
```

### 2. DAC Access Control 설정하기

Connection 또는 Cluster 에 대해, User 에게 Privilege 를 부여합니다. `qpc dac grant` 명령을 사용하는
경우, Local database 에 user, privilege, detailed-connection 정보가 미리 저장되어 있어야 합니다.
```bash
$ ./qpc dac grant <user> <privilege> <cluster|connection> # User 에게 Privilege 를 부여합니다.
```
- \<user\> 는 이용자의 login_id, email, 또는 uuid 를 사용합니다.
- \<privilege\> 는 Privilege 의 name, 또는 uuid 를 사용합니다.
- \<cluster\> 는 Cluster 의 host:port, cloud_identifier, 또는 uuid 를 사용합니다.
- \<connection\> 은 Connection 의 name, 또는 uuid 를 사용합니다.

\<cluster\> 또는 \<connection\> 을 지정하는 경우, Local database 를 조회하여 Cluster 하나를 특정합니다.
둘 이상의 Cluster 가 검색되는 경우, Validation 과정에서 오류 응답을 보여줍니다. 이 경우, Connection 정보를
적절히 구성하였는지, 검토하여 주세요. 동일한 Name 으로 둘 이상의 Connection 이 존재하는 경우, 다른 이용자에게
혼동을 줄 수 있습니다. 동일한 Name 을 사용하는 것이 필요한 경우, Connection 의 UUID 를 사용하여 권한 부여 대상을
지정할 수 있습니다.

지정된 이용자, Cluster 에게 이미 권한이 부여된 경우, `QPS-00003` 오류가 발생합니다.
이미 권한이 부여된 경우, 다른 Privilege 특권을 부여하려는 경우, `--force` 옵션을 사용할 수 있습니다.
```bash
$ ./qpc dac grant <user> <privilege> <cluster|connection> --force
```

## 버전과 구현 상태

### v0.1 - DAC Access Control 의 Grant 기능을 제공합니다.

- DAC Cluster 또는 Connection 에 대해, User 에게 Privilege 를 부여합니다.
- QueryPie 서버로부터, 자산/정책 내려받기 기능
  - User 정보를 local database 에 내려받습니다.
  - DAC Connection, Privilege, Access Control 정보를 local database 에 내려받습니다.

## 구현할 기능

- [x] QueryPie 서버를 지정하여 접속하는 client 실행파일
- [ ] DAC - DB Connection 생성
- [x] DAC - DB Connection 또는 cluster, privilege, user 를 지정하여 접근 권한 부여
- [ ] DAC - table, user 를 지정하여 접근 권한 부여
- [ ] DAC - DB Connection, table 에 Sensitive data 설정 적용
- [ ] DAC - Ledger Policy 적용

## 의견 주고받기

- 고객사, 파트너사
  - QueryPie 와 공유하는 Slack Channel 을 이용해 주세요.
- QueryPie 내부
  - #team-tpm 채널에서 QueryPie Client for Operation 에 대해 질문하고 의견을 남겨주세요.
