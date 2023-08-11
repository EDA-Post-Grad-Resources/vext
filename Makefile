.PHONY: seed

seed:
	go run ./seed/seed.go

.PHONY: dev

dev: 
	@echo "starting dev server + tailwindcss watcher"
	@(go run main.go &)
	@(tailwindcss -i ./styles/input.css -o ./public/output.css --watch &)
	@wait

.PHONY: stop

stop:
	@echo "Stopping any running instances of the application"
	@pkill -f main || true