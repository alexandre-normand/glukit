## General info on the Dexcom Platinum G4:
From [the user manual](http://www.dexcom.com/sites/dexcom.com/files/LBL-011119_Rev_07_User's_Guide,_G4_US.pdf)
Range displayed: between 40 and 400
In reality, the CGM records values outside of that range but the receiver displays "low" or "high" instead. I'm going with the assumption that we might have people
that go over the limits but not by that much. 

glukitScore(average, days) = if average > 83 then (average - 83) * 2 * 288 * days else (83 - average) * 1 * 288 * days
glukit period = 7
glukit score calculation period = glukit period - 1 => 6

Low value average = 69
Lower than normal score = glukitScore(Low value average, glukit score calculation period / 2) => 12096
High value average = 110
Higher than normal score = glukitScore(High value average, glukit score calculation period / 2) => 46656
Rating of 75 percent score = Lower than normal score + Higher than normal score => 58752

Generally good average = 128
Naive glukit score estimate = glukitScore(Generally good average, glukit score calculation period) => 155520

Personal Glukit score from end of March = 48647
Personal Glukit score as of August 15th = 36014
Score improvement between March and now = Personal Glukit score as of July 20th - Personal Glukit score from end of March => Personal Glukit score as of July 20th - 48647

Person with a very bad control (shifting from 50% at 40 and 50% at 400): 
Maximum high = 400
Minimum low = 40
Bad bouncy guy score = glukitScore(Maximum high, glukit score calculation period / 2) + glukitScore(Minimum low, glukit score calculation period / 2) => 584928

Person who is high all the time, I'm hoping nobody does that badly but let's set a score of 0 to this value:
Worst guy ever score = glukitScore(Maximum high, glukit score calculation period) => 1095552

Person with a bad average but still within the notion of high and low limits by the ADA:
Bad average high = 200
Bad average low = 50
Bad average score = glukitScore(Bad average high, glukit score calculation period * 0.7) + glukitScore(Bad average low, glukit score calculation period * 0.3) => 300153.6

Person that is maybe within an OK/high control:
Average person high = 165
Average person low = 60
Average person score = glukitScore(Average person high, glukit score calculation period * 0.7) + glukitScore(Average person low, glukit score calculation period * 0.3) => 210297.6

Perfect score = 0

user facing score(Perfect score) = 100
user facing score(Worst guy ever score) = 0
user facing score(Bad bouncy guy score) = 30
user facing score(Bad average score) = 50
user facing score(Average person score) = 60
user facing score(Naive glukit score estimate) = 65
user facing score(Personal Glukit score from end of March) = 75
user facing score(Personal Glukit score as of July 20th) = 85

##From wolfram alpha
[Calculated from values above](http://www.wolframalpha.com/input/?i=InterpolatingPolynomial%28%7B%7B0%2C+100%7D%2C+%7B88120%2C+83%7D%2C+%7B137088%2C+75%7D%2C+%7B362880%2C+65%7D%2C+%7B500000%2C+60%7D%2C+%7B1000000%2C+50%7D%2C+%7B5000000%2C+30%7D%2C%7B10000000%2C+0%7D%7D%2C+x%29):

## From Eureka
#user facing(s) = 100 + s * (-0.0002298 + 4.36414*10^-10*s - 8.57647*10^-22*s^3 + 7.30739*10^-28*s^4 - 1.63087*10^-34*s^5)
user facing(s) = 100 + 1.043e-9*s^2 + 6.517e-22*s^4 - 0.0003676*s - 1.434e-15*s^3

### Testing some values
This would be a pretty bad score:
user facing(711936) => 16.9071707
This would be a perfect score:
user facing(11110) => 96.0427471
This is my current score (not excellent but still ok):
user facing(48000) => 84.6031426
user facing(0) => 100
