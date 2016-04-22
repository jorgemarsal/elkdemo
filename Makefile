all: setup

setup:
	rm -rf docker-elk
	git clone https://github.com/deviantony/docker-elk.git
	cp logstash.conf docker-elk/logstash/config/logstash.conf 
	cp elastic-search-template.json docker-elk/logstash/config/

