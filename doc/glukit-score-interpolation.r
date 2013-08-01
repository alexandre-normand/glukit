data <- c(xy.coords(0,100), xy.coords(2556288,0), xy.coords(1364832,30), xy.coords(711936,50), xy.coords(498816,60), xy.coords(362880,65), xy.coords(128640,75), xy.coords(88120,84))

data <- xy.coords(list(0, 88120, 128640, 362880, 498816, 711936, 1364832, 2556288), list(100, 84, 75, 65, 60, 50, 30, 0))
f <- splinefun(data)

# From mathematica
s <- 1:2556299; y <- 100 + s * (-0.0002298 + 4.36414*10^-10*s - 8.57647*10^-22*s^3 + 7.30739*10^-28*s^4 - 1.63087*10^-34*s^5)
# The input data based on real values + educated guesses
data <- xy.coords(list(0, 88120, 128640, 362880, 498816, 711936, 1364832, 2556288), list(100, 84, 75, 65, 60, 50, 30, 0))
f <- splinefun(data, method="natural")
ffmm <- splinefun(data, method="fmm")
# This ones seems best because it doesn't go back up again 
fhyman <- splinefun(data, method="hyman")
fmono <- splinefun(data, method="monoH.FC")
plot(f, from=1, to=2556288)
points(list(0, 88120, 128640, 362880, 498816, 711936, 1364832, 2556288), list(100, 84, 75, 65, 60, 50, 30, 0), pch="+")
curve(f, from=1, to=2556288)
curve(ffmm, from=1, to=2556288, col="blue")
curve(fhyman, from=1, to=2556288, col="red")
curve(fmono, from=1, to=2556288, col="green")
lines(s, y, col="purple")
