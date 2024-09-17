use `airplanes`;

ALTER TABLE `users` 
  modify column ID INT AUTO_INCREMENT,
  modify email VARCHAR(255) NOT NULL UNIQUE;

ALTER TABLE `users`
  ADD UNIQUE(`Email`);


