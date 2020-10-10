CREATE DATABASE ip2proxy;
USE ip2proxy;
CREATE TABLE `ip2proxy_px7`(
                               `ip_from` DECIMAL(39,0) UNSIGNED,
                               `ip_to` DECIMAL(39,0) UNSIGNED,
                               `proxy_type` VARCHAR(3),
                               `country_code` CHAR(2),
                               `country_name` VARCHAR(64),
                               `region_name` VARCHAR(128),
                               `city_name` VARCHAR(128),
                               `isp` VARCHAR(256),
                               `domain` VARCHAR(128),
                               `usage_type` VARCHAR(11),
                               `asn` INT(10),
                               `as` VARCHAR(256),
                               PRIMARY KEY (`ip_from`, `ip_to`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

SET GLOBAL local_infile=1;

LOAD DATA LOCAL
    INFILE '/csv/ipv6/IP2PROXY-LITE-PX7.CSV'
    INTO TABLE
    ip2proxy.ip2proxy_px7
    FIELDS TERMINATED BY ','
    ENCLOSED BY '"'
    LINES TERMINATED BY '\n';

LOAD DATA LOCAL
    INFILE '/csv/ipv4/IP2PROXY-LITE-PX7.IPV6.CSV'
    INTO TABLE
    ip2proxy.ip2proxy_px7
    FIELDS TERMINATED BY ','
    ENCLOSED BY '"'
    LINES TERMINATED BY '\n';