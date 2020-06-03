-- MySQL dump 10.13  Distrib 8.0.20, for Linux (x86_64)
--
-- Host: localhost    Database: zauth
-- ------------------------------------------------------
-- Server version	8.0.20-0ubuntu0.20.04.1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `Files`
--

DROP TABLE IF EXISTS `Files`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Files` (
  `ID` int NOT NULL AUTO_INCREMENT,
  `Name` varchar(255) NOT NULL,
  `GroupID` int NOT NULL,
  `FolderID` int DEFAULT NULL,
  `DiskFilename` varchar(255) NOT NULL,
  `FileSize` bigint NOT NULL,
  `Digest` varchar(64) NOT NULL,
  `CreatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `CreatedBy` int NOT NULL,
  `UpdatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `UpdatedBy` int NOT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `UniqName` (`Name`,`FolderID`),
  KEY `GroupFk` (`GroupID`),
  KEY `FolderFk2` (`FolderID`),
  KEY `CreatedByFk2` (`CreatedBy`),
  KEY `UpdatedByFk2` (`UpdatedBy`),
  FULLTEXT KEY `Name` (`Name`),
  CONSTRAINT `CreatedByFk2` FOREIGN KEY (`CreatedBy`) REFERENCES `Users` (`ID`),
  CONSTRAINT `FolderFk2` FOREIGN KEY (`FolderID`) REFERENCES `Folders` (`ID`),
  CONSTRAINT `GroupFk` FOREIGN KEY (`GroupID`) REFERENCES `UserGroups` (`ID`),
  CONSTRAINT `UpdatedByFk2` FOREIGN KEY (`UpdatedBy`) REFERENCES `Users` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Folders`
--

DROP TABLE IF EXISTS `Folders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Folders` (
  `ID` int NOT NULL AUTO_INCREMENT,
  `FolderID` int DEFAULT NULL,
  `Name` varchar(255) NOT NULL,
  `CreatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `CreatedBy` int NOT NULL,
  `UpdatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `UpdatedBy` int NOT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `UniqName` (`Name`,`FolderID`),
  KEY `FolderFk` (`FolderID`),
  KEY `CreatedByFk1` (`CreatedBy`),
  KEY `UpdatedByFk1` (`UpdatedBy`),
  FULLTEXT KEY `Name` (`Name`),
  CONSTRAINT `CreatedByFk1` FOREIGN KEY (`CreatedBy`) REFERENCES `Users` (`ID`),
  CONSTRAINT `FolderFk` FOREIGN KEY (`FolderID`) REFERENCES `Folders` (`ID`),
  CONSTRAINT `UpdatedByFk1` FOREIGN KEY (`UpdatedBy`) REFERENCES `Users` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `User2Group`
--

DROP TABLE IF EXISTS `User2Group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `User2Group` (
  `UserID` int NOT NULL,
  `GroupID` int NOT NULL,
  PRIMARY KEY (`UserID`,`GroupID`),
  KEY `GroupID` (`GroupID`),
  CONSTRAINT `User2Group_ibfk_1` FOREIGN KEY (`UserID`) REFERENCES `Users` (`ID`),
  CONSTRAINT `User2Group_ibfk_2` FOREIGN KEY (`GroupID`) REFERENCES `UserGroups` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `UserGroups`
--

DROP TABLE IF EXISTS `UserGroups`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `UserGroups` (
  `ID` int NOT NULL AUTO_INCREMENT,
  `Name` varchar(200) NOT NULL,
  `Description` text,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `uniq_username` (`Name`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Users`
--

DROP TABLE IF EXISTS `Users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Users` (
  `ID` int NOT NULL AUTO_INCREMENT,
  `Username` varchar(200) NOT NULL,
  `FirstName` varchar(200) NOT NULL,
  `LastName` varchar(200) NOT NULL,
  `Email` varchar(300) NOT NULL,
  `PasswordHash` varchar(300) NOT NULL DEFAULT '-',
  `PasswordSet` datetime NOT NULL DEFAULT '0001-01-01 00:00:00',
  `LastLogin` datetime NOT NULL DEFAULT '0001-01-01 00:00:00',
  `Disabled` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`ID`),
  UNIQUE KEY `uniq_username` (`Username`)
) ENGINE=InnoDB AUTO_INCREMENT=187 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping routines for database 'zauth'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2020-06-03 14:26:13
