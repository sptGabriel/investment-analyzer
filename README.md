# Investment-Analyzer

Este repositório contém a solução para o desafio técnico, focado na análise de rentabilidade de investimentos de Alice ao longo de um período específico.
## Sumário

- [Contexto do Desafio](#contexto-do-desafio)
- [Arquivos de Dados](#arquivos-de-dados)
- [Especificações Principais](#especificações-principais)
  - [Cálculo do Patrimônio Total](#cálculo-do-patrimônio-total)
  - [Cálculo da Rentabilidade Acumulada](#cálculo-da-rentabilidade-acumulada)
- [Requisitos e Dependências](#requisitos-e-dependências)
  - [Go](#go)
  - [Ferramentas de Terceiros](#ferramentas-de-terceiros)
  - [Docker/Docker Compose](#dockerdocker-compose)
- [Instalação e Configuração](#instalação-e-configuração)
  - [Clonando o Repositório](#clonando-o-repositório)
  - [Variáveis de Ambiente](#variáveis-de-ambiente)
  - [Buildfile](#buildfile)
- [Comandos Principais via Makefile](#comandos-principais-via-makefile)
- [Execução da Solução](#execução-da-solução)
  - [Execução via binário (local)](#execução-via-binário-local)
  - [Execução via Docker Compose](#execução-via-docker-compose)
- [Exemplo de Uso via cURL](#exemplo-de-uso-via-curl)
- [Acesso à Documentação Swagger](#acesso-à-documentação-swagger)
- [Resposta à Pergunta](#resposta-à-pergunta)

---

## Contexto do Desafio

Alice tem comprado e vendido dois ativos (A e B) ao longo do mês de março de 2021. Ela deseja verificar se sua estratégia de múltiplas operações diárias tem resultado em boa rentabilidade comparada a uma simples compra única e manutenção de um dos ativos.

O objetivo é **gerar relatórios** que mostrem:

1. **Patrimônio Total** em determinados intervalos de tempo (por exemplo, a cada 10 minutos).
2. **Rentabilidade Acumulada** ao longo do período especificado.

### Arquivos de Dados

- **march_2021_trades.csv**: Lista de todas as operações (BUY/SELL) realizadas por Alice.
- **march_2021_pricesA.csv** e **march_2021_pricesB.csv**: Arquivos de preços dos ativos A e B, amostrados minuto a minuto.

---

## Especificações Principais

### Cálculo do Patrimônio Total
No instante \( t \), o patrimônio total é calculado como:  
\[
  \text{Patrimônio Total}(t) = \text{Dinheiro em Caixa}(t) + 
  \bigl[\text{PreçoA}(t) \times \text{UnidadesA}(t)\bigr] + 
  \bigl[\text{PreçoB}(t) \times \text{UnidadesB}(t)\bigr]
\]
Alice começa o mês de março com R\$ 100.000,00 em caixa.

### Cálculo da Rentabilidade Acumulada
A rentabilidade acumulada até o instante \( t \) é:  
\[
  \text{Rent}(t) = \frac{\text{Patrimônio Total}(t)}{\text{Patrimônio Total}(t_0)} - 1
\]
onde \( t_0 \) é o instante inicial (por exemplo, o início do período de análise).

---

## Requisitos e Dependências

### Go
- Versão **1.23.6** ou superior (verificado via `buildfile.yaml`)

### Ferramentas de Terceiros
- **Make** para executar os alvos do Makefile.
- **Ferramentas Go**: `moq`, `gotest`, `swaggo`, `staticcheck`, `govulncheck`, `gci`, etc.
- **yq** (para parsing do `buildfile.yaml`). O Makefile tenta instalar automaticamente via `go install` se você não tiver localmente.

### Docker/Docker Compose
- Necessário caso deseje executar a aplicação em contêiner.  
- Caso não vá utilizar contêineres, basta ter Go instalado localmente e executar o binário.

---

## Instalação e Configuração

### Clonando o Repositório

```bash
git clone https://github.com/usuario/investment-analyzer.git
cd investment-analyzer
```

### Comandos Principais via Makefile

```bash
make setup
```

Verifica dependências essenciais (Go, Docker, yq etc.).
Instala ferramentas (linters, geradores, etc.).
Executa docker compose pull para puxar imagens definidas no docker-compose.yml.

```bash
make build
```

Compila o projeto e gera binários a partir do buildfile.yaml. Ele:

Garante que o yq esteja instalado (caso contrário, instala).
Lê o nome do aplicativo e as pastas de origem.
Produz o binário em ./bin.


```bash
make docker-up
```

Para e remove quaisquer contêineres em execução localmente.
Sobe os serviços definidos em docker-compose.yml.


## Execução da Solução

### Execução via binário (local)

```bash
make setup
make build
./bin/investment-analyzer-api
```

### Execução via Docker Compose

```bash
make setup
make build
make docker-up
```

## Exemplo de Uso via cURL

Para obter um relatório de rentabilidade e patrimônio (intervalo de 10 minutos) de 1 de março de 2021 até 7 de março de 2021, por exemplo:

```bash
curl --location 'localhost:3000/api/v1/investment_analyzer/portfolios/408186c6-b76a-4ad6-8d4a-9ace3762b997/reports' \
  --header 'Content-Type: application/json' \
  --data '{
      "start_date": "2021-03-01 10:00:00",
      "end_date":   "2021-03-07 17:50:00",
      "interval":   "10m"
  }'
```

### Acesso à Documentação Swagger

Se configurado, a rota Swagger estará em:

```bash
http://localhost:3000/docs/v1/investment_analyzer/swagger/index.html
```

(Ajuste a porta conforme seu .env ou docker-compose.yml.)

## Resposta à Pergunta

### Pergunta:
“Considerando um patrimônio inicial de R$ 100.000,00, os arquivos de exemplo (march_2021_trades.csv, march_2021_pricesA.csv e march_2021_pricesB.csv) e o retorno que Alice conseguiu com seus trades, seria melhor para Alice ter comprado 100% no ativo A ou 100% no ativo B no início do dia ao invés de operar ao longo do dia?”


