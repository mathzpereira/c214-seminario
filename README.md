# Lista de Contatos

API para lista de contatos, desenvolvida como parte da disciplina C214 - Engenharia de Software.

## Tecnologias Utilizadas

- Go
- Gin (Framework HTTP)
- Testes utilizando o pacote `testing` do Go
- CircleCI para o pipeline de CI/CD

## Instalação

1. **Certifique-se de que o Go está instalado**:
 Você pode verificar a instalação do Go executando `go version` no terminal.

2. **Clone o repositório**:
   ```bash
   git clone https://github.com/mathzpereira/c214-seminario.git
   cd c214-seminario/contact-list-api
   ```

3. **Instale as dependências**:
   ```bash
   go mod tidy
   ```
## Rodando localmente

Já estando dentro da pasta do projeto, inicie o server com

```bash
  go run main.go
```

## Rodando os testes

Para rodar os testes:
```
go test
```

Para rodar os testes e criar um report em html instale o pacote e depois rode os testes com as linhas de comando abaixo:
```
go get github.com/vakenbolt/go-test-report@latest
go install github.com/vakenbolt/go-test-report@latest
go test -json ./tests | go-test-report -o test_report.html
```