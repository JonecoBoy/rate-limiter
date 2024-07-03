# Rate Limiter em Go

## Descrição

Este projeto implementa um rate limiter em Go que controla o tráfego de requisições para um servidor web com base no endereço IP ou api token.

## Configuração

1. Clone o repositório.
2. Crie um arquivo `.env` na raiz do projeto com as variáveis de configuração.
3. Execute `docker-compose up` para iniciar o servidor e o Redis.

## Endpoints

- `GET /`: Retorna uma mensagem de boas-vindas.

## Variáveis de Ambiente

- `REDIS_ADDR`: Endereço do Redis.
- `RATE_LIMIT_STRATEGY`: Tipo de estratégia de armazenamento (memory ou redis).
- `RATE_LIMIT_IP`: Limite de requisições por segundo por IP.
- `RATE_LIMIT_TOKEN`: Limite de requisições por segundo por Token.
- `BLOCK_TIME`: Tempo de bloqueio em segundos.

## Testes

Execute os testes com o comando `go test ./...`.
