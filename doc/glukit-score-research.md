glukitScore(average, days) = if average > 83 then (average - 83) * 2 * 288 * days else (83 - average) * 1 * 288 * days

Low value average = 69
Lower than normal score = glukitScore(Low value average, 7) => 28224
High value average = 110
Higher than normal score = glukitScore(High value average, 7) => 108864
Rating of 75 percent score = Lower than normal score + Higher than normal score => 137088

Generally good average = 128
Naive glukit score estimate = glukitScore(Generally good average, 14) => 362880

Personal Glukit score from end of March = 128640
Personal Glukit score as of July 20th = 88120
