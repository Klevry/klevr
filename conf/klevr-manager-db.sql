CREATE DATABASE  IF NOT EXISTS `klevr` /*!40100 DEFAULT CHARACTER SET utf8 */;
USE `klevr`;
-- MariaDB dump 10.17  Distrib 10.4.13-MariaDB, for osx10.15 (x86_64)
--
-- Host: 127.0.0.1    Database: klevr
-- ------------------------------------------------------
-- Server version	10.5.4-MariaDB-1:10.5.4+maria~focal

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `AGENTS`
--

-- DROP TABLE IF EXISTS `AGENTS`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `AGENTS` (
  `ID` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '에이전트 ID',
  `AGENT_KEY` varchar(45) DEFAULT NULL COMMENT '에이전트 고유 식별 키 (향후 unique 키로 만들어야 함 - 현재는 개발중)',
  `GROUP_ID` bigint(20) unsigned NOT NULL COMMENT '에이전트 그룹 ID',
  `IS_ACTIVE` tinyint(1) NOT NULL DEFAULT 0 COMMENT '에이전트 활성 여부',
  `LAST_ALIVE_CHECK_TIME` timestamp NULL DEFAULT NULL COMMENT '마스터를 통해 마지막 에이전트 생존 확인 시간',
  `LAST_ACCESS_TIME` timestamp NULL DEFAULT NULL COMMENT '마지막 agent가 매니저에 직접 엑세스한 타임스탬프',
  `IP` varchar(45) DEFAULT NULL COMMENT '에이전트 IP (동일 에이전트 그룹간 통신 가능한 IP)',
  `PORT` int(11) DEFAULT NULL,
  `HMAC_KEY` varchar(45) DEFAULT NULL COMMENT '전송 데이터 위변조 검사용 키',
  `ENC_KEY` varchar(45) DEFAULT NULL COMMENT '전송구간 데이터 암호화 키',
  `CREATED_AT` timestamp NULL DEFAULT current_timestamp(),
  `UPDATED_AT` timestamp NULL DEFAULT current_timestamp(),
  `DELETED_AT` timestamp NULL DEFAULT NULL,
  `CPU` int(11) DEFAULT NULL,
  `MEMORY` int(11) DEFAULT NULL,
  `DISK` int(11) DEFAULT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8 COMMENT='Agent 테이블\n전체 agent 정보를 관리';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `AGENT_GROUPS`
--

-- DROP TABLE IF EXISTS `AGENT_GROUPS`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `AGENT_GROUPS` (
  `ID` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '그룹 ID',
  `GROUP_NAME` varchar(200) NOT NULL COMMENT '그룹 이름',
  `USER_ID` bigint(20) unsigned NOT NULL COMMENT '사용자 ID',
  `PLATFORM` varchar(45) NOT NULL COMMENT '플랫폼(baremetal, k8s, aws 등)',
  `CREATED_AT` timestamp NULL DEFAULT current_timestamp(),
  `UPDATED_AT` timestamp NULL DEFAULT current_timestamp(),
  `DELETED_AT` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `UNQ_01` (`USER_ID`,`PLATFORM`,`GROUP_NAME`),
  KEY `IDX_01` (`USER_ID`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='agent 그룹 관리 테이블';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `API_AUTHENTICATIONS`
--

-- DROP TABLE IF EXISTS `API_AUTHENTICATIONS`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `API_AUTHENTICATIONS` (
  `API_KEY` varchar(64) NOT NULL COMMENT 'API key',
  `GROUP_ID` bigint(20) unsigned NOT NULL COMMENT '그룹 ID',
  `CREATED_AT` timestamp NULL DEFAULT current_timestamp(),
  `UPDATED_AT` timestamp NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`API_KEY`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='인증 테이블';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `PRIMARY_AGENTS`
--

-- DROP TABLE IF EXISTS `PRIMARY_AGENTS`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `PRIMARY_AGENTS` (
  `GROUP_ID` bigint(20) unsigned NOT NULL COMMENT '에이전트의 group ID',
  `AGENT_ID` bigint(20) unsigned NOT NULL COMMENT '에이전트 ID',
  `LAST_ACCESS_TIME` timestamp NULL DEFAULT NULL COMMENT '에이전트가 마지막으로 엑세스한 타임스탬프',
  `CREATED_AT` timestamp NULL DEFAULT current_timestamp(),
  `UPDATED_AT` timestamp NULL DEFAULT current_timestamp(),
  `DELETED_AT` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`GROUP_ID`,`AGENT_ID`),
  UNIQUE KEY `GROUP_ID_UNIQUE` (`GROUP_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Primary 에이전트 관리 테이블\n그룹별 하나의 primary 에이전트를 관리한다.';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `TASK_LOCK`
--

-- DROP TABLE IF EXISTS `TASK_LOCK`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `TASK_LOCK` (
  `TASK` varchar(45) NOT NULL COMMENT 'lock을 잡은 task 명',
  `INSTANCE_ID` varchar(45) NOT NULL COMMENT 'lock을 잡은 인스턴스 ID',
  `LOCK_DATE` timestamp NULL DEFAULT NULL COMMENT 'lock이 걸린 일시',
  PRIMARY KEY (`TASK`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='배치성 작업을 위해 lock을 관리하는 테이블';
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2020-07-30 11:19:00

CREATE USER 'klevr'@'%' identified by 'klevr';

GRANT ALL PRIVILEGES ON klevr.* to `klevr`@`%`;