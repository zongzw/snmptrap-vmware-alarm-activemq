# snmptrap2activemq
a sample of snmp-trapping vmware alarm to activemq.

```# ./trap2mq --help
Usage of ./trap2mq:
  -mqpass string
    	activemq pass (default "2yxALph4")
  -mqsite string
    	activemq site [ip:port] (default "0.0.0.0:61613")
  -mquser string
    	activemq user (default "activemq_admin")
  -queue string
    	queue name, i.e. /queue/Remedy (default "/queue/Remedy")
  -testtimes int
    	the count of message will be read from queue for test (default 10)
  -trapd string
    	snmp trapd [ip:port] (default "0.0.0.0:162")```
