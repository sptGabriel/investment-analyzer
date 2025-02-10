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

para fins da finalização do desáfio, só será possível com esse portfolio-id que é o de alice: 408186c6-b76a-4ad6-8d4a-9ace3762b997.

### Acesso à Documentação Swagger

Se configurado, a rota Swagger estará em:

```bash
http://localhost:3000/docs/v1/investment_analyzer/swagger/index.html
```

(Ajuste a porta conforme seu .env ou docker-compose.yml.)

## Resposta à Pergunta

### Pergunta:
“Considerando um patrimônio inicial de R$ 100.000,00, os arquivos de exemplo (march_2021_trades.csv, march_2021_pricesA.csv e march_2021_pricesB.csv) e o retorno que Alice conseguiu com seus trades, seria melhor para Alice ter comprado 100% no ativo A ou 100% no ativo B no início do dia ao invés de operar ao longo do dia?”


Se formos considera o Exemplo 1:

### Exemplo 1

Alice gostaria de ver o relatório de 1 de março de 2021 até 7 de março de 2021 com uma taxa de amostragem de 10 min.

#### 3 primeiras linhas do relatório

| timestamp           | Patrimônio Total | Rentabilidade Acumulada |
| ------------------- | ---------------- | ----------------------- |
| 2021-03-01 10:00:00 |        100.000,0 |                 0,00000 |
| 2021-03-01 10:10:00 |        100.024,0 |                 0,00024 |
| 2021-03-01 10:20:00 |         99.919,0 |                -0,00081 |

#### 3 últimas linhas do relatório

| timestamp           | Patrimônio Total | Rentabilidade Acumulada |
| ------------------- | ---------------- | ----------------------- |
| 2021-03-07 17:30:00 |         99.575,0 |                -0,00425 |
| 2021-03-07 17:40:00 |         98.972,0 |                -0,01028 |
| 2021-03-07 17:50:00 |         99.397,0 |                -0,00603 |

onde o patrimônio total no final do dia foi de: 99.397,0, a melhor estrátegia é comprar o ativo B no cenário 'buy and hold'.
Ativo B - Patrimônio final: R$99958.11 (Compra a R$23.87, Venda a R$23.75)

Já no exemplo 2:

### Exemplo 2

Alice gostaria de ver o relatório de 5 de março de 2021 até 12 de março de 2021 com uma taxa de amostragem de 10 min.

#### 3 primeiras linhas do relatório

| timestamp           | Patrimônio Total | Rentabilidade Acumulada |
| ------------------- | ---------------- | ----------------------- |
| 2021-03-05 10:00:00 |         99.168,0 |                0,000000 |
| 2021-03-05 10:10:00 |         99.019,0 |               -0,001503 |
| 2021-03-05 10:20:00 |         99.262,0 |                0,000948 |

#### 3 últimas linhas do relatório

| timestamp           | Patrimônio Total | Rentabilidade Acumulada |
| ------------------- | ---------------- | ----------------------- |
| 2021-03-12 17:30:00 |         99.062,0 |               -0,001069 |
| 2021-03-12 17:40:00 |         99.508,0 |                0,003429 |
| 2021-03-12 17:50:00 |         99.317,0 |                0,001503 |

 o patrimônio total no final do dia foi de: 99.317,0
 
  Ativo A - Patrimônio final: R$98999.17 (Compra a R$23.97, Venda a R$23.92)
  Ativo B - Patrimônio final: R$100841.75 (Compra a R$23.87, Venda a R$23.75)

Melhor comprar o ativo B no cenário 'buy and hold'.