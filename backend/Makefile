.PHONY: build dev dev-css dev-server css clean

build: css
	go build -o ayomendaki-admin .

dev: css
	@echo "Development mode: run these in separate terminals:"
	@echo "  1. npx @tailwindcss/cli -i web/static/css/input.css -o web/static/css/tailwind.css --watch"
	@echo "  2. go run . --dev"
	@echo ""
	@echo "Starting dev server..."
	go run . --dev

css:
	npx @tailwindcss/cli -i web/static/css/input.css -o web/static/css/tailwind.css

clean:
	rm -f ayomendaki-admin
	rm -f tmp/ayomendaki-admin
	rm -rf data/
