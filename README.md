# teste-carga-go-redis-lock
# Teste de carga com lock distribuído usando Go e Redis

Testando desempenho de ação mutamente exclusiva, bloqueando por identificadores diferentes.

Serve como amostra de aplicação para setup rápido.

```bash
# Dependências, apenas no início ou se novas forem instaladas
docker compose run --rm app go install

docker compose up
docker compose exec -it app go run main.go
```
