# dns-json-go
DNS JSON parser

## dns to json

```
msg := new(dns.Msg)
...
dnsMsg := &dnsjson.DNSMsg{msg}
respJSON, err := dnsMsg.DNS2JSON()
respStr, err := json.Marshal(respJSON)
log.Println(string(jsonByte))
```

## json to dns
```
jsonByte = ...
...
jsonMsg := &jsondns.JSONMsg{}
json.Unmarshal(jsonByte, jsonMsg)
log.Println(jsonMsg)
```