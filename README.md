# snmptrap2activemq
a sample of snmp-trapping vmware alarm to activemq.
## Usage:
```# ./trap2mq --help
Usage of ./trap2mq:
  -mqpass string
    	activemq pass (default "")
  -mqsite string
    	activemq site [ip:port] (default "0.0.0.0:61613")
  -mquser string
    	activemq user (default "")
  -queue string
    	queue name, i.e. /queue/myq (default "/queue/myq")
  -testtimes int
    	the count of message will be read from queue for test (default 10)
  -trapd string
    	snmp trapd [ip:port] (default "0.0.0.0:162")
```

## Startup activemq container for test:

`docker run -p 61616:61616 -p 8161:8161 -p 61613:61613 rmohr/activemq`

## VMware SnmpTrap setting.
From vSphere Web Client: Click vCenter Instance -> Monitor -> Alarm -> Configuration.
