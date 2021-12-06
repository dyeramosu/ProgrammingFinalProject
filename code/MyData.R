#Deepika Yeramosu
#Programming for Scientists
#Final Project

#install.packages("ggplot2")
library(ggplot2)

setwd("/Users/deepikayeramosu/Go/src/finalproject/")

mydata <- read.table("finalMap", header=F, row.names = 1, nrows = 7, fileEncoding="UTF8")
#get number of cols (number of boards/steps in CA algorithm)
mycols <- ncol(mydata)
colnames(mydata) <- c(1:(mycols))
rownames(mydata) <- as.character(rownames(mydata))

mydata <- t(mydata)
mydata <- mydata[ , order(colnames(mydata))]
 
mydata1 <- data.frame(time = c(1:(mycols)))
mydata1$compFreq = mydata[, 1]
mydata1$enzFreq = mydata[, 2]
mydata1$inhibCompFreq = mydata[, 3]
mydata1$inhibFreq = mydata[, 4]
mydata1$prodFreq = mydata[, 5]
mydata1$subFreq = mydata[, 6]

View(mydata1)

ggplot(mydata1, aes(x=time)) + 
  geom_line(aes(y = prodFreq), color = "purple") + 
  geom_line(aes(y = subFreq), color="red") +
  geom_line(aes(y = compFreq), color = "cyan") +
  geom_line(aes(y = enzFreq), color = "light green") +
  geom_line(aes(y = inhibFreq), color = "blue") +
  geom_line(aes(y = inhibCompFreq), color = "black") +
  ylab("frequencies") + ggtitle("Michaelis-Menten Using Cellular Automaton")
