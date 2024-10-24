use `airplanes`;

INSERT INTO `questions`
  (`ID`, `Text`) 
VALUES 
  (1, 'Please rate our aircraft flown on AMONIC Airlines'),
  (2, 'How would you rate our filght attendants'),
  (3, 'How would you rate our infilght entertainment'),
  (4, 'Please rate the ticket price for the trip you are taking')
;


INSERT INTO `question_answers`
  (`QuestionID`, `Value`, `Text`) 
VALUES
  (1, 1, 'Outstanding'),
  (1, 2, 'Very good'),
  (1, 3, 'Good'),
  (1, 4, 'Adequate'),
  (1, 5, 'Needs improvement'),
  (1, 6, 'Poor'),
  (1, 7, "Don't know"),
  
  (2, 1, 'Outstanding'),
  (2, 2, 'Very good'),
  (2, 3, 'Good'),
  (2, 4, 'Adequate'),
  (2, 5, 'Needs improvement'),
  (2, 6, 'Poor'),
  (2, 7, "Don't know"),

  (3, 1, 'Outstanding'),
  (3, 2, 'Very good'),
  (3, 3, 'Good'),
  (3, 4, 'Adequate'),
  (3, 5, 'Needs improvement'),
  (3, 6, 'Poor'),
  (3, 7, "Don't know"),

  (4, 1, 'Outstanding'),
  (4, 2, 'Very good'),
  (4, 3, 'Good'),
  (4, 4, 'Adequate'),
  (4, 5, 'Needs improvement'),
  (4, 6, 'Poor'),
  (4, 7, "Don't know");



INSERT INTO `groups`
  (`ID`, `Name`)
VALUES
  (1, 'Gender'),
  (2, 'Age'),
  (3, 'TravelClass'),
  (4, 'Arrival'),
  (5, 'Departure');



INSERT INTO `group_values`
  (`GroupID`, `Name`)
VALUES
  (1, 'Male'),
  (1, 'Female'),

  (2, '18-24'),
  (2, '25-39'),
  (2, '40-59'),
  (2, '60+'),
  
  (3, 'Economy'),
  (3, 'Business'),
  (3, 'First'),

  (4, 'AUH'),
  (4, 'BAH'),
  (4, 'DOH'),
  (4, 'BYU'),
  (4, 'CAI'),

  (5, 'AUH'),
  (5, 'BAH'),
  (5, 'DOH'),
  (5, 'BYU'),
  (5, 'CAI');

