build/examples:
	rm -rf dist && mkdir -p dist
	go build -o dist/simple ./examples/simple
	gopherjs build -m -o dist/gopherjs.js ./examples/gopherjs
