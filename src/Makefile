up:
	@go run cmd/logogenerator/main.go
	@mv -f amazing_logo.png static/images/
	@rm -rf assets.zip
	@zip -r assets.zip static templates
	@docker compose up -d

down:
	@docker compose down
	@rm -rf assets.zip
