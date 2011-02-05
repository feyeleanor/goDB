include $(GOROOT)/src/Make.inc

TARG=goDB

GOFILES=\
	goDB.go\
	transaction.go

include $(GOROOT)/src/Make.pkg