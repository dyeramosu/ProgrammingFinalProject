#Deepika Yeramosu
#Programming for Scientists
#Final Project

#install.packages("GillespieSSA")
library(GillespieSSA)

setwd("/Users/deepikayeramosu/Go/src/finalproject/")

#parms are the p values 
  #k1 corresponds to S + E --> SE
  #k2 corresponds to SE --> S + E
  #k3 corresponds to SE --> S + P
parms <- c(k1=6.67, k2=1.0, k3=1.0)
#x0 is a vector that initializes the starting values of each molecule
x0 <- c(S=300, E=100, ES=0, P=0)
#a is a vector that defines the propensities of the reactions 
a <- c("k1*S*E", "k2*ES", "k3*ES")
#nu is the change matrix 
nu <- matrix(c(-1, -1, +1, 0, +1, +1, -1, 0, 0, +1, -1, +1), nrow=4, byrow=F)

out<-ssa(x0, a, nu, parms, tf=100, method = ssa.d(), simName = "Michaelis-Menten Using Gillespie")
ssa.plot(out)
