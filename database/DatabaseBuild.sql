	/* drop database if exists pest; */
	/* create database pest; */
	
	use pest;
	
	
	drop table if exists Submissions;	
	drop table if exists Assignments;	
	drop table if exists StudentCourses;
	drop table if exists CourseDescription;
	drop table if exists Users;


	create table Users
	(
		UserID int AUTO_INCREMENT PRIMARY KEY,
		FirstName varchar (30),
		MiddleInitial varchar (1),
		LastName varchar (30),
		Username varchar (35) NOT NULL UNIQUE,
		Email varchar (35) NOT NULL UNIQUE,
		/* change length of password when she provides it */
		Password varchar (20) NOT NULL,
		PrivLevel tinyint NOT NULL,
		LastLogin datetime,
		PwdChangeFlag boolean,
		NumLoginAttempts tinyint
	);

	/* change this number to indicate where user id's start at */
	ALTER TABLE Users AUTO_INCREMENT=10001;


	create table CourseDescription
	(
		CourseName varchar(42) PRIMARY KEY,
		/* size of display name field? */
		CourseDisplayName varchar(40),
		CourseDescription text,
		Instructor int NOT NULL,
		StartDate datetime NOT NULL,
		EndDate datetime NOT NULL,
		SI1 int,
		SI2 int,
		SIGradeFlag boolean,
		SITestCaseFlag boolean,

		FOREIGN KEY (Instructor) REFERENCES Users(UserID)
			ON UPDATE CASCADE
			/* so, deleting an instructor also deletes his/her courses, is this what we want? */
			ON DELETE CASCADE,

		FOREIGN KEY (SI1) REFERENCES Users(UserID)
			ON UPDATE CASCADE
			ON DELETE SET NULL,

		FOREIGN KEY (SI2) REFERENCES Users(UserID)
			ON UPDATE CASCADE
			ON DELETE SET NULL
	);


	create table StudentCourses
	(
		Student int NOT NULL,
		CourseName varchar(42) NOT NULL,

		FOREIGN KEY (Student) REFERENCES Users(UserID)
			ON UPDATE CASCADE
			ON DELETE CASCADE,

		FOREIGN KEY (CourseName) REFERENCES CourseDescription(CourseName)
			ON UPDATE CASCADE
			ON DELETE CASCADE
	);


	create table Assignments
	(
		CourseName varchar(42) NOT NULL,
		/* size for assignment name field? */
		AssignmentName varchar(40) NOT NULL,
		StartDate datetime,
		EndDate datetime,
		/* is this seconds or milliseconds? */
		MaxRuntime int,
		CompilerOptions text,
		NumTestCases int,

		FOREIGN KEY (CourseName) REFERENCES CourseDescription(CourseName)
			ON UPDATE CASCADE
			ON DELETE CASCADE,

		PRIMARY KEY (CourseName, AssignmentName)
	);


	create table Submissions
	(
		CourseName varchar(42),
		/* again, size for assignment name field? */
		AssignmentName varchar(40),
		Student int,
		Grade tinyint,
		comment varchar(240),
		Compile bool,
		Results text,
		SubmissionNumber smallint,

		FOREIGN KEY (CourseName, AssignmentName) REFERENCES Assignments (CourseName, AssignmentName)
			ON UPDATE CASCADE
			ON DELETE CASCADE,

		FOREIGN KEY (Student) REFERENCES Users(UserID)
			ON UPDATE CASCADE
			ON DELETE CASCADE
	);


	/*
		Add some sample data
		Users, classes, assignments, etc.
	*/

	insert into Users values (NULL, "James", "A", "Jerkins", "jajerkins", "jajerkins@una.edu", "1234", 15, NULL, true, 0);
	insert into Users values (NULL, "Mark", "G", "Terwilliger", "mterwilliger", "mterwilliger@una.edu", "PerlRules", 10, NULL, true, 0);
	insert into Users values (NULL, "Patricia", "L", "Roden", "plroden", "plroden@una.edu", "password", 10, NULL, true, 0);
	insert into Users values (NULL, "John", "W", "Doe", "jdoe", "jdoe@una.edu", "password", 1, NULL, true, 0);
	insert into Users values (NULL, "Eileen", "R", "Drass", "edrass", "edrass@una.edu", "Mattie9423!", 5, NULL, true, 0);
	insert into Users values (NULL, "Turner", "J", "Brian", "bmiller1", "bmiller1@una.edu", "Whomp2!", 1, NULL, true, 0);
	insert into Users values (NULL, "Martin", "W", "Nicholas", "nmartin2", "nmartin2@una.edu", "unaLions16!", 1, NULL, false, 0);
	insert into Users values (NULL, "Smith", "E", "Jessica", "jsmith4", "jsmith4@una.edu", "waFFles108!", 0, NULL, true, 0);
	insert into Users values (NULL, "Bradly", "A", "Updike", "bupdike", "bupdike@una.edu", "MyRealPass17!", 1, NULL, true, 0);
	insert into Users values (NULL, "Edward", "J", "Snowden", "esnowden", "esnowden@una.edu", "TryMeNSA13!", 1, NULL, true, 0);
	insert into Users values (NULL, "William", "H", "Gates", "wgates1", "wgates1@una.edu", "I<3Jobs", 1, NULL, false, 0);
	insert into Users values (NULL, "Tyrion", "I", "Lannister", "tlannister", "tlannister@una.edu", "IDrink&Know", 0, NULL, true, 0);
	insert into Users values (NULL, "Abdullah", "M", "Karaman", "akaraman", "akaraman@una.edu", "WhatsApp123!", 1, NULL, true, 0);
	insert into Users values (NULL, "Richard", "J", "Hendricks", "rhendricks", "rhendricks@una.edu", "awkwardMan2017!", 1, NULL, true, 0);
	insert into Users values (NULL, "Erlich", "M", "Bachman", "ebachman", "ebachman@una.edu", "AviatoListen2016!", 1, NULL, false, 0);
	insert into Users values (NULL, "Hannah", "M", "Hopkins", "hhopkins", "hhopkins@una.edu", "ginger12", 5, NULL, true, 0);
	insert into Users values (NULL, "Frank", "L", "Liles", "fliles", "fliles@una.edu", "mysteryd0g2", 1, NULL, true, 0);
	insert into Users values (NULL, "Cody", "C", "Klein", "cklein", "cklein@una.edu", "ilovehannahlol", 1, NULL, true, 0);
	insert into Users value (NULL, "Albert", "K", "Einstein", "aeinstein", "aeinstein@una.edu", "eequalsmcsquare", 15, NULL, true, 0);
	insert into Users values (NULL, "Chance", "T", "Rapper", "crapper", "crapper@una.edu", "n0pr0bl3m", 1, NULL, true, 0);
	insert into Users values (NULL, "Timothy", "T", "Tanner", "ttanner1", "ttanner1@una.edu", "passWoRd", 1, NULL, true, 0);
	insert into Users values (NULL, "Bob", "G", "Dylan", "bdylan1", "bdylan1@una.edu", "SingingForAll", 1, NULL, true, 0);
	insert into Users values (NULL, "Elvis", "A", "Pressley", "epressley1", "epressley1@una.edu", "PiNKCADDy123", 1, NULL, true, 0);
	insert into Users values (NULL, "David", "R", "Johnson", "djohnson2", "djohnson2@una.edu", "DAvvYJ", 1, NULL, true, 0);
	insert into Users values (NULL, "Sheila", "H", "Jordan", "sjordan1", "sjordan1@una.edu", "OleSheila", 1, NULL, true, 0);
	

	insert into CourseDescription values ("TerwilligerCS15501SP17", "Computer Science I", "This is CS 1.", (select UserID from Users where username="mterwilliger"), NOW(), NOW(), NULL, NULL, false, false);
	insert into CourseDescription values ("JerkinsCS15502SP17", "Computer Science I", "This is Dr. Jerkins' CS 1.", (select UserID from Users where username="jajerkins"), NOW(), NOW(), NULL, NULL, false, false);

	insert into StudentCourses values ((select UserID from Users where username="jdoe"), "TerwilligerCS15501SP17");
	insert into StudentCourses values ((select UserID from Users where username = "nmartin2"), "TerwilligerCS15501SP17");
	insert into StudentCourses values ((select UserID from Users where username = "bmiller1"), "JerkinsCS15502SP17");

	insert into Assignments values ("TerwilligerCS15501SP17", "Assignment 1", NOW(), NOW(), 1000, "-Wall -std=c++0x", 2);
	insert into Assignments values ("JerkinsCS15502SP17", "Assignment 0", NOW(), NOW(), 1000, "-Wall -std=c++0x", 2);

	insert into Submissions values ("TerwilligerCS15501SP17", "Assignment 1", (select UserID from Users where username="jdoe"), 90, "Good program.", true, "Some results...", 1);
	insert into Submissions values ("TerwilligerCS15501SP17", "Assignment 1", (select UserID from Users where username = "nmartin2"), 90, "Good program.", true, "Some results...", 1);
	insert into Submissions values ("JerkinsCS15502SP17", "Assignment 0", (select UserID from Users where username="bmiller1"), 80, "Decent program. Work on style.", true, "Some results...", 1);

	/* check tables */
	show tables;