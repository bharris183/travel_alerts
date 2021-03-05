CREATE DATABASE IF NOT EXISTS threat_alerts CHARACTER SET latin1 COLLATE latin1_swedish_ci;

CREATE  TABLE IF NOT EXISTS threat_alerts.RSS_THREATS (
  COUNTRY_CODE CHAR(2) NOT NULL,
  THREAT_LEVEL INT NOT NULL ,
  TITLE VARCHAR(255) ,
  LINK VARCHAR(255) ,
  DESCRIPTION VARCHAR(255) , 
  PUB_DATE DATE ,
  PRIMARY KEY (COUNTRY_CODE) )
ENGINE = InnoDB;

CREATE USER 'sa_threatalerts'@'localhost' IDENTIFIED BY 'sw0rdfish';
grant ALL PRIVILEGES on threat_alerts.* to 'root'@'%' identified by 'sw0rdfish';
grant SELECT, INSERT on threat_alerts.* to 'sa_threatalerts'@'localhost' identified by 'sw0rdfish';
FLUSH PRIVILEGES;