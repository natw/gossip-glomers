default: single

TESTER := ../maelstrom/maelstrom
BIN := broadcast

single: $(BIN)
	$(TESTER) test -w broadcast --bin $(BIN) --node-count 1 --time-limit 20 --rate 10

multi: $(BIN)
	$(TESTER) test -w broadcast --bin $(BIN) --node-count 5 --time-limit 20 --rate 10

$(BIN): *.go ../go.sum
	go build .
