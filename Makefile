include .env
export

### Копирования переменных окружения
copy:
	cp .env-example .env

### Пинг LLM
ping-llm:
	@if [ -z "$$MODEL_PORT" ]; then \
		echo "MODEL_PORT не задан (экспортируй или положи в .env)"; \
		exit 1; \
	fi
	curl -I "http://localhost:$(MODEL_PORT)"

### Сборка проекта в бинарник
build:
	go build -o eng-util .