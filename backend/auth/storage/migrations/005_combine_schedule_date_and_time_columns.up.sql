use `airplanes`;

-- Step 1: Add a new DATETIME column
ALTER TABLE `schedules`
-- ADD COLUMN Outbound timestamp(6);
ADD COLUMN Outbound DATETIME;

-- Step 2: Update the new column with combined values from date and time columns
UPDATE `schedules`
SET Outbound = CONCAT(Date, ' ', Time);

ALTER TABLE `schedules` 
MODIFY COLUMN Outbound DATETIME NOT NULL;