use `airplanes`;

ALTER TABLE `users` modify column ID INT AUTO_INCREMENT;

ALTER TABLE `users` ADD UNIQUE(`Email`);
