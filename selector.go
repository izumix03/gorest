package gorest

// Get sets method
func (cli *client) Get() TerminalOperator {
	cli.method = get
	return cli
}

// Post sets method
func (cli *client) Post() TerminalOperator {
	cli.method = post
	return cli
}
