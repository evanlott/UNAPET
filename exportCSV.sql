SELECT FirstName, MiddleInitial, LastName, Username INTO OUTFILE '/home/hannah/test.csv'
FIELDS TERMINATED BY ',' OPTIONALLY ENCLOSED BY '"'
LINES TERMINATED BY '\n'
FROM Users WHERE UserID IN (SELECT Student FROM StudentCourses WHERE CourseName="JerkinsCS15502SP17")
