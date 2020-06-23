/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;


-- Groups is now a MySQL reserved keyword, so we must rename our table
RENAME Table `Groups` TO UserGroups;

-- Change UserGroups.UserID FK to delete rows when the User is deleted
ALTER TABLE `User2Group` DROP FOREIGN KEY `User2Group_ibfk_1`;
ALTER TABLE `User2Group` ADD CONSTRAINT `User2Group_ibfk_1` FOREIGN KEY (`UserID`) REFERENCES `Users` (`ID`) ON DELETE CASCADE;

-- Change UserGroups.GroupID FK to delete rows when the Group is deleted
ALTER TABLE `User2Group` DROP FOREIGN KEY `User2Group_ibfk_2`;
ALTER TABLE `User2Group` ADD CONSTRAINT `User2Group_ibfk_2` FOREIGN KEY (`GroupID`) REFERENCES `UserGroups` (`ID`) ON DELETE CASCADE;

/*!40101 SET character_set_client = @saved_cs_client */;
