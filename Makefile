all: build

build:
	go build

install: build
	cp seaport /usr/local/bin/

config:
	mkdir -p /etc/seaport/
	cp seaport.yaml /etc/seaport/seaport.yaml

systemd:
	cp seaport.service /etc/systemd/system/
	systemctl daemon-reload
	systemctl enable seaport
	systemctl start seaport
