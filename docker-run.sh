#!/usr/bin/env sh

docker network create derennia
docker run -d --name elasticsearch --net derennia -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" elasticsearch:7.14.0
docker run -d --name kibana --net derennia -p 5601:5601 kibana:7.14.0
docker run -d --name prometheus --hostname prometheus -v /Users/gmedina/Projects/kibanatest/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml -p 9090:9090 prom/prometheus
docker run -d --name grafana --hostname grafana -p 7070:3000 grafana/grafana
