FROM mysql:5
RUN apt-get update && apt-get install -y unzip
ADD IP2PROXY-LITE-PX7.CSV.zip /zip/IP2PROXY-LITE-PX7.CSV.zip
ADD IP2PROXY-LITE-PX7.IPV6.CSV.zip /zip/IP2PROXY-LITE-PX7.IPV6.CSV.zip
RUN mkdir -p /csv/ipv6
RUN mkdir -p /csv/ipv4
RUN unzip /zip/IP2PROXY-LITE-PX7.CSV.zip -d /csv/ipv6
RUN unzip /zip/IP2PROXY-LITE-PX7.IPV6.CSV.zip -d /csv/ipv4
COPY init.sql /docker-entrypoint-initdb.d/init.sql
