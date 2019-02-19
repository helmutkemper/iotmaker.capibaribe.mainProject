# Capibaribe

> This project is in test.

This is the preliminary version of my version of the reverse proxy, designed to be simple and fast. 
The initial idea is to allow a single machine to receive multiple small sites, and can host then at a low 
cost.

To run this project:

```zsh

  # I assume you have a file named 'capibaribe-config.yml' in the same directory.
  ./builds/Linux_Ubuntu_V_1_0

```

or

```zsh

  ./builds/Linux_Ubuntu_V_1_0 -f path/you_yaml_config.yml

```

How to configure:

```yaml

version: 1.0

capibaribe:

  affluentRiverName:

    listen: :8081

    ssl:
      enabled: true
      certificate: /etc/nginx/company.com.crt
      certificateKey: /etc/nginx/company.com.key
      sslProtocols:
        - TLSv1
        - TLSv1.1
        - TLSv1.2

    static:
      - filePath: /docker/static
        serverPath: static

    pygocentrus:
      enabled: true
      dontRespond: 0.0
      changeLength: 0.0
      changeContent:
        changeRateMin: 0.01
        changeRateMax: 0.10
        changeBytesMin: 1
        changeBytesMax: 1
        rate: 0.5
      deleteContent: 0.0
      changeHeaders:
        - number: 500
          header:
            - key: Content-Type
              value: application/json
          rate: 0.1

    proxy:
      - ignorePort: true
        host: 127.0.0.1
        # roundRobin -              joga a carga para o próximo servidor
        # lowTimeResponseHeader -   procura o servidor com menor tempo de resposta inicial
        # lowTimeResponseLastByte - procura o servidor com menor tempo de resposta do último byte
        # random -                  joga a carga de forma aleatória
        # overLoad -                joga todas a carga para o primeiro servidor até o número máximo de conexões ser atingido
        # ipv4Hash -                000.000.*: de onde tiver o asterisco vai para um servidor
        # ipv6Hash -                lista de ips
        # hash -                    request_uri
#        bind:
#          - host: 127.0.0.1
#            ignorePort: true
        loadBalancing: roundRobin # roundRobin | lowTimeResponseHeader | lowTimeResponseLastByte | random | overLoad | ipv4Hash | ipv6Hash | hash
        path: /
        maxAttemptToRescueLoop: 10
        healthCheck:
          enabled: true
          # padrão: 5 segundos
          interval: 5000
          # numeros de falhas consecutivas para erro
          fails: 3
          # numeros de ok consecutivos para zerar falhas
          passes: 2
          # caminho de teste
          uri: /some/path
          # se 0 suspende para sempre
          # se !0 suspende por milesegundos
          suspendInterval: 30000
          match:
            status:
              - expReg: expreg
                value: 300
                in:
                  - min: 200
                    max: 299
                notIn:
                  - min: 200
                    max: 299
            header:
              - key: Content-Type
                value: text/html
            body:
              - expreg
        servers:
          - host: http://localhost:3000
            weight: 1
            overLoad: 1000000
          - host: http://localhost:3000
            weight: 1
            overLoad: 1000000
          - host: http://localhost:3000
            weight: 1
            overLoad: 1000000

```
