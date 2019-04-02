module github.com/robertzml/Glaucus

go 1.12

require (
	github.com/eclipse/paho.mqtt.golang v1.1.1
	github.com/gomodule/redigo v2.0.0+incompatible
	golang.org/x/net v0.0.0-20190320064053-1272bf9dcd53 // indirect
)

replace (
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2 => github.com/golang/crypto v0.0.0-20190308221718-c2843e01d9a2
	golang.org/x/net v0.0.0-20190320064053-1272bf9dcd53 => github.com/golang/net v0.0.0-20190320064053-1272bf9dcd53
	golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a => github.com/golang/sys v0.0.0-20190215142949-d0b11bdaac8a
	golang.org/x/text v0.3.0 => github.com/golang/text v0.3.0
)
