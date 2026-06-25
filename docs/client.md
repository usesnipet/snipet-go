# Client

O cliente é o aplicativo que irá interagir com o bot. Ele pode ser um aplicativo mobile, um site, um aplicativo desktop, etc.

## Integração
### Webhook
1. Chama a API do Snipet de autenticação com um token terceiro
2. O Snipet envia uma requisição para o webhook do cliente configurado com o token terceiro
3. O webhook deve retornar um JSON com informações do usuario
4. O Snipet usa o ID do usuario do cliente para identificar o usuario e gera um JWT para o usuario

### JWT
1. Ao logar no app cliente, ele vai enviar uma requisição para a API do Snipet com as informações do usuario
2. A API do Snipet vai verificar se o usuario existe, ou criar um novo usuario se não existir
3. A API do Snipet vai gerar um JWT para o usuario
4. O JWT é retornado
5. No frontend do app cliente, o JWT é usado para fazer as requisições para a API do Snipet