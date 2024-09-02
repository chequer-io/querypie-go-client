# QueryPie Client for Operation

QueryPie Client for Operation 은 쿼리파이 운영을 위해 사용할 수 있는 cli client 를 제공합니다.
Command Line Interface 로 작동하며, 자동화 프로그램을 위해 사용할 수 있으며, 운영자가 터미널에서
정보를 조회하거나 정책 설정을 관리하는 용도로 사용할 수 있습니다.

## 버전과 구현 상태

아직 초기 Draft 상태이며, 릴리즈 버전을 갖고 있지 않습니다.

## 구현할 기능

- [ ] 쿼리파이 서버를 지정하여 접속하는 client 실행파일
- [ ] DAC - DB Connection 생성
- [ ] DAC - 특정 connection, table 에 Sensitive data 설정 적용
- [ ] DAC - connection, table, user 를 지정하여 접근 권한 부여
- [ ] DAC - Ledger Policy 적용

## 초기 디자인 아이디어

- 이용자는 db hostname 또는 connection name, schema name, table name 을
  argument 로 사용하여, 자산 등록, 정책 설정하고자 합니다.
- 일 단위 배치 작업으로 실행하고자 합니다.
- 1차적으로 자산 생성, 권한 부여, 정책 적용하는 기능을 우선적으로 구현합니다. 
  이후, 자산 삭제, 권한 회수, 적용 취소하는 기능을 구현할 수 있습니다.
- External API 를 활용하여 효율적으로 작동하기 위하여, 등록된 자산을 로컬 db 에 저장합니다.

## 의견 주고받기

- #team-tpm 채널에서 QueryPie Client for Operation 에 대해 질문하고 의견을 남겨주세요.
