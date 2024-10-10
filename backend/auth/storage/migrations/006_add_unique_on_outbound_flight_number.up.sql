use `airplanes`;

ALTER TABLE schedules
ADD CONSTRAINT constraint_date_and_flight_number UNIQUE (Date, FlightNumber);