package main

type JsonMessageLogger struct{}

func (logger *JsonMessageLogger) LogMessage(mes Message) error {
	return nil
}

func (logger *JsonMessageLogger) GetName() string {
	return "json"
}

type CsvMessageLogger struct{}

func (logger *CsvMessageLogger) LogMessage(mes Message) error {
	return nil
}

func (logger *CsvMessageLogger) GetName() string {
	return "json"
}

type ConsoleMessageLogger struct{}

func (logger *ConsoleMessageLogger) LogMessage(mes Message) error {
	return nil
}

func (logger *ConsoleMessageLogger) GetName() string {
	return "json"
}

type SqliteMessageLogger struct{}

func (logger *SqliteMessageLogger) LogMessage(mes Message) error {
	return nil
}

func (logger *SqliteMessageLogger) GetName() string {
	return "sqlite"
}
