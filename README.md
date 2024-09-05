# Desafio 4 do curso de Go da dio.me

## Criando uma api rest em Go para catálogo de bolos

- adiciona informação de novos bolos (post)
- atualiza informações de bolos existentes (put)
- atualiza informações de comentários e votações de bolos existentes (patch)
- retorna o catalogo completo dos bolos (get)
- retorna informações de um bolo específico no catálogo (get)

## Observações

- Carregando o endereço localhost:8086 no browser será carregado uma referencia da api (não é swagger é uma página estática)
- A pasta shell conta com scripts para teste utilizando curl, basta executar por exemplo ./shell/get-all-data.sh ** (executar em bash para windows pode-se instalar o cmder ou usar o wsl)