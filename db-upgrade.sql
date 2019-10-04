-- TODO: 
-- 
-- Deactivate Accounts With LastLogin > 90 Days Ago
--   Send the deactivated users an email before it happens? (30 days before)
-- 
-- mysql zauth -e "SELECT FirstName, LastName, LastLogin FROM Users WHERE LastLogin < (CURDATE() - INTERVAL 60 DAY)";
-- 
-- Require Group for Redmine Login
-- Require Group for Git Login
--
-- Config options for: Full domain/URL, Website Name, Reply Email, Email Name, ...
--
-- Allow admins to SEND a password reset link to a given user? Or put on page?


ALTER TABLE Users 
	DROP COLUMN UserID,
	DROP COLUMN GroupID,
	DROP COLUMN HomeDirectory,
	MODIFY PasswordHash varchar(300) NOT NULL DEFAULT '-',
	MODIFY LastLogin datetime NOT NULL DEFAULT '0001-01-01 00:00:00',
	ADD COLUMN PasswordSet datetime NOT NULL DEFAULT '0001-01-01 00:00:00' AFTER PasswordHash;

-- Disable Shared Accounts
UPDATE Users 
	SET Disabled=1 
	WHERE FirstName IN ('Test', 'Map User 1', 'Map User 2', 'Map User 3', 'SVN', 'LDAP', 'Test', 'Synthetic Teammate', 'Demo', 'BOINC User', 'Air Force User', 'Redmine Account');