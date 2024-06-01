# Teste de carga com lock distribuído usando Go e Redis

Testando desempenho de ação mutamente exclusiva, bloqueando por identificadores diferentes.

Serve como amostra de aplicação para setup rápido.

```bash
# Dependências, apenas no início ou se novas forem instaladas
docker compose run --rm app go install

docker compose up
docker compose exec -it app go run main.go
```

# Profilador

Captura:

```bash
go tool pprof -seconds 10 http://localhost:8080/debug/pprof/profile
```

_(nesse ínterim executamos o endpoint da aplicação)_

```bash
curl http://localhost:8080/dotask/12345
```

Exemplo de comando: `top [n]`