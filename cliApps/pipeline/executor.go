package main

type executor interface {
	execute() (string, error)
}
