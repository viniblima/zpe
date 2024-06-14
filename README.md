# ZPE

## Visão Geral

Este projeto é uma aplicação em Go que serve de controle/autenticação de usuários. Ele foi desenvolvido para servir como teste para a ZPE.

## Estrutura do Projeto

- cmd/ - Arquivos de comando da aplicação
- controllers/ - Arquivos de código que estão diretamente ligados à aplicação principal, acessados via endpoints
- database/ - Arquivos de código interno do projeto que inicia e instancia o banco de dados
- routes/ - Definições de API e documentos
- handlers/ - Arquivos de código com funções auxiliares que podem ser acessadas em outros packages
- middlewares/ - Arquivos de código com funções que podem ser executadas antes de própriamente acessar os controllers em um endpoint. Validador de token, por exemplo
- models/ - Arquivos de modelos para o banco de dados

## Requisitos

- Go versão 1.16 ou superior
- Docker

## Instalação

1. Clone o repositório:

   ```sh
   git clone https://github.com/viniblima/zpe.git
   cd zpe
   ```

2. Suba o container:
   ```sh
   docker-compose up
   ```

## Sobre

Neste projeto, foi observado a necessidade de aprimorar alguns pontos. Visando a segurança em uma API que seria posta em ambiente de produção, foi-se criado um endpoint de login e algumas rotas (GET /roles, por exemplo) foram protegidas com um token JWT, ainda que não apontado nos requisitos do teste.
