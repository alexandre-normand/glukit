## General info on the Dexcom Platinum G4:
From [the user manual](http://www.dexcom.com/sites/dexcom.com/files/LBL-011119_Rev_07_User's_Guide,_G4_US.pdf)
Range displayed: between 40 and 400
In reality, the CGM records values outside of that range but the receiver displays "low" or "high" instead. I'm going with the assumption that we might have people
that go over the limits but not by that much. 

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
Score improvement between March and now = Personal Glukit score as of July 20th - Personal Glukit score from end of March => -40520

Person with a very bad control (shifting from 50% at 40 and 50% at 400): 
Maximum high = 400
Minimum low = 40
Bad bouncy guy score = glukitScore(Maximum high, 7) + glukitScore(Minimum low, 7) => 1364832

Person who is high all the time, I'm hoping nobody does that badly but let's set a score of 0 to this value:
Worst guy ever score = glukitScore(Maximum high, 14) => 2556288

Person with a bad average but still within the notion of high and low limits by the ADA:
Bad average high = 200
Bad average low = 50
Bad average score = glukitScore(Bad average high, 10) + glukitScore(Bad average low, 4) => 711936

Person that is maybe within an OK/high control:
Average person high = 165
Average person low = 60
Average person score = glukitScore(Average person high, 10) + glukitScore(Average person low, 4) => 498816

Perfect score = 0

user facing score(Perfect score) = 100
user facing score(Worst guy ever score) = 0
user facing score(Bad bouncy guy score) = 30
user facing score(Bad average score) = 50
user facing score(Average person score) = 60
user facing score(Naive glukit score estimate) = 65
user facing score(Personal Glukit score from end of March) = 75
user facing score(Personal Glukit score as of July 20th) = 83

##From wolfram alpha
[Calculated from values above](http://www.wolframalpha.com/input/?i=InterpolatingPolynomial%28%7B%7B0%2C+100%7D%2C+%7B88120%2C+83%7D%2C+%7B137088%2C+75%7D%2C+%7B362880%2C+65%7D%2C+%7B500000%2C+60%7D%2C+%7B1000000%2C+50%7D%2C+%7B5000000%2C+30%7D%2C%7B10000000%2C+0%7D%7D%2C+x%29):

user facing(s) = 100 + s * (-0.0002298 + 4.36414*10^-10*s - 8.57647*10^-22*s^3 + 7.30739*10^-28*s^4 - 1.63087*10^-34*s^5)

### Testing some values
This would be a pretty bad score:
user facing(711936) => 49.6792139
This would be a perfect score:
user facing(11110) => 97.5007766
This is my current score (not excellent but still ok):
user facing(90000) => 82.8009114

