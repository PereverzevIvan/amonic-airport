use `airplanes`;

-- Команды по изменению колонок в таблице users, 
-- должна быть напрямую перенесена в миграцию 001
-- для избежания проблем с FK

-- ALTER TABLE `tickets` 
--   DROP FOREIGN KEY `FK_Ticket_User`;

-- ALTER TABLE `users` 
--   MODIFY column ID INT AUTO_INCREMENT,
--   MODIFY column Email VARCHAR(255) COLLATE utf8_bin NOT NULL UNIQUE,
--   MODIFY column Password varchar(255) COLLATE utf8_bin NOT NULL;

-- ALTER TABLE `tickets` 
--   ADD CONSTRAINT `FK_Ticket_User` FOREIGN KEY (`UserID`) REFERENCES `users`(`ID`) ON DELETE CASCADE ON UPDATE CASCADE;
