use `airplanes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;

DROP TABLE IF EXISTS `surveys`;
CREATE TABLE surveys (
  ID INT NOT NULL AUTO_INCREMENT,
  Date date NOT NULL,
  RespondentsCount int DEFAULT 0,

  PRIMARY KEY (`ID`)
);


DROP TABLE IF EXISTS `questions`;
CREATE TABLE questions (
  ID INT NOT NULL AUTO_INCREMENT,
  Text varchar(255) NOT NULL,
  
  PRIMARY KEY (`ID`)
);


DROP TABLE IF EXISTS `question_answers`;
CREATE TABLE question_answers (
  ID INT NOT NULL AUTO_INCREMENT,
  QuestionID INT NOT NULL,
  Value INT NOT NULL DEFAULT 0,
  Text varchar(255) NOT NULL,

  PRIMARY KEY (`ID`),
  UNIQUE(`QuestionID`,`Value`),
  CONSTRAINT `FK_QuestionAnswers_Questions` FOREIGN KEY (`QuestionID`) REFERENCES `questions` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION
);



DROP TABLE IF EXISTS `survey_questions`;
CREATE TABLE survey_questions (
  SurveyID INT NOT NULL,
  QuestionID INT NOT NULL,
  
  PRIMARY KEY (`SurveyID`, `QuestionID`),
  CONSTRAINT `FK_SurveyQuestions_Surveys` FOREIGN KEY (`SurveyID`) REFERENCES `surveys` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_SurveyQuestions_Questions` FOREIGN KEY (`QuestionID`) REFERENCES `questions` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION
);



DROP TABLE IF EXISTS `groups`;
CREATE TABLE `groups` (
  ID INT NOT NULL AUTO_INCREMENT,
  Name VARCHAR(255) NOT NULL UNIQUE,
  PRIMARY KEY (`ID`)
);

DROP TABLE IF EXISTS `group_values`;
CREATE TABLE group_values (
  ID INT NOT NULL AUTO_INCREMENT,
  GroupID INT NOT NULL,  
  Name varchar(255) NOT NULL,

  PRIMARY KEY (`ID`),
  CONSTRAINT `FK_GroupValues_Groups` FOREIGN KEY (`GroupID`) REFERENCES `groups` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION
);



DROP TABLE IF EXISTS `respondents`;
CREATE TABLE respondents (
  ID INT NOT NULL AUTO_INCREMENT,
  SurveyID INT NOT NULL,  

  PRIMARY KEY (`ID`),
  CONSTRAINT `FK_Respondents_Surveys` FOREIGN KEY (`SurveyID`) REFERENCES `surveys` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION
);

DROP TABLE IF EXISTS `respondent_group_values`;
CREATE TABLE `respondent_group_values` (
  RespondentID INT NOT NULL,
  GroupValueID INT NOT NULL,
  
  UNIQUE (`RespondentID`,`GroupValueID`),
  CONSTRAINT `FK_ResondentGroupValues_Resondent` FOREIGN KEY (`RespondentID`) REFERENCES `respondents` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_ResondentGroupValues_GroupValues` FOREIGN KEY (`GroupValueID`) REFERENCES `group_values` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION
);


DROP TABLE IF EXISTS `respondent_answers`;
CREATE TABLE respondent_answers (
  RespondentID INT NOT NULL,
  QuestionAnswerID INT NOT NULL,

  UNIQUE (`RespondentID`,`QuestionAnswerID`),
  CONSTRAINT `FK_RespondentsAnswers_Respondents` FOREIGN KEY (`RespondentID`) REFERENCES `respondents` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_RespondentsAnswers_QuestionAnswers` FOREIGN KEY (`QuestionAnswerID`) REFERENCES `question_answers` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION
);

