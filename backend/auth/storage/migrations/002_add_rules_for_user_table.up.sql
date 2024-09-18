use `airplanes`;

ALTER TABLE `users` 
  modify column ID INT AUTO_INCREMENT,
  modify column Email VARCHAR(255) COLLATE utf8_bin NOT NULL UNIQUE;


