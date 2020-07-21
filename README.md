# dns-json-go
DNS JSON parser

##dns to json

```
msg := new(dns.Msg)
...
dnsMsg := &DNSMsg{msg}
jsonByte, _ := dnsMsg.DNS2JSON()
log.Println(string(jsonByte))
```

##json to dns
```
jsonByte = ...
...
jsonMsg := &jsondns.JSONMsg{}
json.Unmarshal(jsonByte, jsonMsg)
log.Println(jsonMsg)
```