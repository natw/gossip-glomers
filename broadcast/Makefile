default: single

single: broadcast
	../maelstrom/maelstrom test -w broadcast --bin broadcast --node-count 1 --time-limit 20 --rate 10

multi: broadcast
	../maelstrom/maelstrom test -w broadcast --bin broadcast --node-count 1 --time-limit 20 --rate 10

broadcast: *.go ../go.sum
	go build .
