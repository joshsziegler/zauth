/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;

-- ALTER TABLE Users
-- 	DROP COLUMN UserID,
-- 	DROP COLUMN GroupID,
-- 	DROP COLUMN HomeDirectory,
-- 	MODIFY PasswordHash varchar(300) NOT NULL DEFAULT '-',
-- 	MODIFY LastLogin datetime NOT NULL DEFAULT '0001-01-01 00:00:00',
-- 	ADD COLUMN PasswordSet datetime NOT NULL DEFAULT '0001-01-01 00:00:00' AFTER PasswordHash;

-- Groups is now a MySQL reserved keyword, so we must rename our table
RENAME Table `Groups` TO UserGroups;

-- Change UserGroups.UserID FK to delete rows when the User is deleted
ALTER TABLE `User2Group` DROP FOREIGN KEY `User2Group_ibfk_1`;
ALTER TABLE `User2Group` ADD CONSTRAINT `User2Group_ibfk_1` FOREIGN KEY (`UserID`) REFERENCES `Users` (`ID`) ON DELETE CASCADE;

-- Change UserGroups.GroupID FK to delete rows when the Group is deleted
ALTER TABLE `User2Group` DROP FOREIGN KEY `User2Group_ibfk_2`;
ALTER TABLE `User2Group` ADD CONSTRAINT `User2Group_ibfk_2` FOREIGN KEY (`GroupID`) REFERENCES `UserGroups` (`ID`) ON DELETE CASCADE;

SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `Folders`;
CREATE TABLE `Folders` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `FolderID` int(11) DEFAULT NULL, -- NULL folder means it's in the base of the project/Group
  `Name` varchar(255) NOT NULL,
  `CreatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `CreatedBy` int(11), -- UserID
  `UpdatedAt` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `UpdatedBy` int(11),
  PRIMARY KEY (`ID`),
  UNIQUE KEY `UniqName` (`Name`, `FolderID`),
  FULLTEXT(`Name`),
   CONSTRAINT `FolderFk` FOREIGN KEY (`FolderID`) REFERENCES `Folders` (`ID`) ON DELETE CASCADE,
  CONSTRAINT `CreatedByFk1` FOREIGN KEY (`CreatedBy`) REFERENCES `Users` (`ID`) ON DELETE SET NULL,
  CONSTRAINT `UpdatedByFk1` FOREIGN KEY (`UpdatedBy`) REFERENCES `Users` (`ID`) ON DELETE SET NULL
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `Files`;
CREATE TABLE `Files` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `Name` varchar(255) NOT NULL,
  `GroupID` int(11) NOT NULL, -- Should this be LDAP GroupID or a ProjectID?
  `FolderID` int(11) DEFAULT NULL, -- NULL folder means it's in the base of the project/Group
  `DiskFilename` varchar(255) NOT NULL,
  `FileSize` bigint(20) NOT NULL,
  `Digest` varchar(64) NOT NULL,
  -- MimeType ??
  `CreatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `CreatedBy` int(11), -- UserID
  `UpdatedAt` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `UpdatedBy` int(11),
  PRIMARY KEY (`ID`),
  UNIQUE KEY `UniqName` (`Name`, `FolderID`),
  FULLTEXT(`Name`),
  CONSTRAINT `GroupFk` FOREIGN KEY (`GroupID`) REFERENCES `UserGroups` (`ID`) ON DELETE CASCADE,
  CONSTRAINT `FolderFk2` FOREIGN KEY (`FolderID`) REFERENCES `Folders` (`ID`) ON DELETE CASCADE,
  CONSTRAINT `CreatedByFk2` FOREIGN KEY (`CreatedBy`) REFERENCES `Users` (`ID`) ON DELETE SET NULL,
  CONSTRAINT `UpdatedByFk2` FOREIGN KEY (`UpdatedBy`) REFERENCES `Users` (`ID`) ON DELETE SET NULL
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;
/*!40101 SET character_set_client = @saved_cs_client */;
