# Complete List of Commands

## User Management
- [x] qpc user fetch-all
- [x] qpc user ls
- [ ] qpc upsert-user <loginid> --email=<email> --name=<name>
- [ ] qpc user activate <user>
- [ ] qpc user deactivate <user>
- [ ] qpc user delete <user>
- [ ] qpc user reset-password <user>
- [ ] qpc user describe <user>

## DAC Management
- [x] qpc dac fetch-all connections
- [x] qpc dac fetch-all detailed-connections
- [x] qpc dac fetch-all access-controls
- [x] qpc dac fetch-all privileges
- [x] qpc dac ls connections
- [x] qpc dac ls detailed-connections
- [x] qpc dac ls access-controls
- [x] qpc dac ls privileges
- [x] qpc dac ls clusters
- [x] qpc dac grant <user> <privilege> <cluster|connection>

## DAC Troubleshooting
- [x] qpc dac fetch-by-uuid connection <uuid>
- [x] qpc dac find-by-uuid connection <uuid>

## SAC Server Management
- [ ] qpc sac fetch-all
- [ ] qpc sac fetch <name>
- [ ] qpc sac upsert <name> --hostname=<hostname> --port=<port> ...
- [ ] qpc sac delete <name>
- [ ] qpc sac describe <name>
- [ ] qpc sac upsert-server-group <name> --tag name1=value1 --tag name2=value2
- [ ] qpc sac upsert-server-group <name> --tag=name1=value1 --tag=name2=value2
- [ ] sac sac upsert-server-group <name> --tag="{name1=value1, name2=value2, name3=value3}"

## Note
- 위의 명령어는 현재 구현되어 있지 않습니다. 개발 진행에 따라, 명령어가 추가되거나 변경될 수 있습니다.
