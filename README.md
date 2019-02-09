# reverseProxy

This project is not finished, so it is not in a free repository

```yaml

ReverseProxy:
  Config:                                       # Configuração principal
    listenAndServer: :9090                      # Servidor local. Formato: server:port 
    outputConfig: true                          # Imprime a configuração inicial na saída padrão. Formato: true/false 
    staticServer: true                          # Habilita um servidor de arquivos no servidor principal. Formato: true/false
    staticFolder: /docker/static                # Caminho da pasta do servidor de arquivos

  Proxy:                                        # Servidor proxy. Redireciona o endereço para um novo servidor
   
    blog:                                       # Nome a ser gravado no log
      Host: blog.localhost:8888                 # Endereço de entrada
      Server:                                   # Lista de servidores
      - docker 1 - ok:http://locahost:2368      # nome a ser gravado no log:endereço
      - docker 2 - error:http://locahost:2369   # nome a ser gravado no log:endereço
      - docker 3 - error:http://locahost:2370   # nome a ser gravado no log:endereço

    blog_2:                                     # Nome a ser gravado no log
      Host: blog2.localhost:8888                # Endereço de entrada
      Server:                                   # Lista de servidores
      - docker 1 - ok:http://locahost:2368      # nome a ser gravado no log:endereço
      - docker 2 - error:http://locahost:2369   # nome a ser gravado no log:endereço
      - docker 3 - error:http://locahost:2370   # nome a ser gravado no log:endereço

```