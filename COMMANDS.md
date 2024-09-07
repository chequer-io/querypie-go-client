# Complete List of Commands

## User Management
- [ ] qpc fetch-all users
- [ ] qpc fetch-user <loginid>
- [ ] qpc upsert-user <loginid> --email=<email> --name=<name>
- [ ] qpc activate-user <loginid>
- [ ] qpc deactivate-user <loginid>
- [ ] qpc delete-user <loginid>
- [ ] qpc reset-password <loginid>
- [ ] qpc describe-user <loginid>

## DAC Connection Management
- [ ] qpc fetch-all dac-connections
- [ ] qpc fetch-dac-connection <name>
- [ ] qpc upsert-dac-connection <name> --hostname=<hostname> --port=<port> ...
- [ ] qpc delete-dac-connection <name>
- [ ] qpc describe-dac-connection <name>

## SAC Server Management
- [ ] qpc fetch-all sac-servers
- [ ] qpc fetch-sac-server <name>
- [ ] qpc upsert-sac-server <name> --hostname=<hostname> --port=<port> ...
- [ ] qpc delete-sac-server <name>
- [ ] qpc describe-sac-server <name>
- [ ] qpc upsert-sac-server-group <name> --tag name1=value1 --tag name2=value2
- [ ] qpc upsert-sac-server-group <name> --tag=name1=value1 --tag=name2=value2
- [ ] qpc upsert-sac-server-group <name> --tag="{name1=value1, name2=value2, name3=value3}"

## Note
- 위의 명령어는 현재 구현되어 있지 않습니다. 개발 진행에 따라, 명령어가 추가되거나 변경될 수 있습니다.
