.PHONY: build

build:
	pnpm build
	go build -o pmail_telegram_push -ldflags "-s -w"

install:
	cp -v pmail_telegram_push ./plugins

clean:
	rm -vf pmail_telegram_push ./plugins/pmail_telegram_push
