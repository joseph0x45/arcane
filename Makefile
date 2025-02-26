BINARY=arcane.out

build:
	go build -o $(BINARY) .

tailwind-compile:
	npx @tailwindcss/cli -i ./input.css -o ./output.css

tailwind-watch:
	npx @tailwindcss/cli -i ./input.css -o ./output.css --watch
