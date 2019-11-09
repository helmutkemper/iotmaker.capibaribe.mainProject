# Capibaribe

Fast and simple reverse proxy

> MVP 1.0 ok

To run this project:

```zsh

  # I assume you have a file named 'capibaribe-config.yml' in the same directory.
  ./builds/Linux_Ubuntu_V_1_0

```

or

```zsh

  ./builds/Linux_Ubuntu_V_1_0 -f path/you_yaml_config.yml

```

## SSL

SSL max and min versions:

| Version                                    | Value                                      |
|--------------------------------------------|--------------------------------------------|
| TLS 1.0                                    | 10                                         |
| TLS 1.1                                    | 11                                         |
| TLS 1.2                                    | 12                                         |
| SSL 3.0                                    | 30                                         |

```yaml

    ssl:
      enabled: true
      version:
        min: 10
        max: 30

```

Certificate and certificate key

Filenames containing a certificate and matching private key for the server must be provided if neither the Server's 
TLSConfig.Certificates nor TLSConfig.GetCertificate are populated. If the certificate is signed by a certificate 
authority, the certFile should be the concatenation of the server's certificate, any intermediates, and the CA's 
certificate.


```yaml

    ssl:
      enabled: true
      certificate: ./company.com.crt
      certificateKey: ./company.com.key

```

CurvePreferences contains the elliptic curves that will be used in an ECDHE handshake, in preference order. If empty, 
the default will be used.

See https://www.iana.org/assignments/tls-parameters/tls-parameters.xml#tls-parameters-8

```yaml

    ssl:
      enabled: true
      curvePreferences:
        - P256
        - P384
        - P521
        - X25519

```

PreferServerCipherSuites controls whether the server selects the client's most preferred ciphersuite, or the server's 
most preferred ciphersuite. If true then the server's preference, as expressed in the order of elements in CipherSuites, 
is used.
	
```yaml

    ssl:
      enabled: true
      preferServerCipherSuites: true

```

CipherSuites is a list of supported cipher suites. If CipherSuites is nil, TLS uses a list of suites supported by the 
implementation.

```yaml

    ssl:
      enabled: true
      cipherSuites:
        - TLS_RSA_WITH_RC4_128_SHA
        - TLS_RSA_WITH_3DES_EDE_CBC_SHA
        - TLS_RSA_WITH_AES_128_CBC_SHA
        - TLS_RSA_WITH_AES_256_CBC_SHA
        - TLS_RSA_WITH_AES_128_CBC_SHA256
        - TLS_RSA_WITH_AES_128_GCM_SHA256
        - TLS_RSA_WITH_AES_256_GCM_SHA384
        - TLS_ECDHE_ECDSA_WITH_RC4_128_SHA
        - TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA
        - TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA
        - TLS_ECDHE_RSA_WITH_RC4_128_SHA
        - TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA
        - TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA
        - TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA
        - TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256
        - TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256
        - TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
        - TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256
        - TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
        - TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384
        - TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305
        - TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305
        - TLS_FALLBACK_SCSV

```

Reads and parses a public/private key pair from a pair of files. The files must contain PEM encoded data. The 
certificate file may contain intermediate certificates following the leaf certificate to form a certificate chain. 

```yaml

    ssl:
      enabled: true
      x509:
        certificate: ./company.com.crt
        certificateKey: ./company.com.key

```

proxy is the reverse proxy functionality, where a service can be pointed to multiple routes within the 
network.

| Key                      | Value                                                                                       |
|--------------------------|---------------------------------------------------------------------------------------------|
| ignorePort               | if true, process all requests to a given address host                                       |
| host                     | address to be monitored                                                                     |
| bind                     | allowed list of addresses                                                                   |
| loadBalancing            | determines how load balancing will be made for each route on the server                     |
| path                     | does not apply to proxy functionality                                                       |
| maxAttemptToRescueLoop   | maximum number of attempts before an error is propagated to the end user                    |
| healthCheck              | checks the integrity of a route at periodic intervals                                       |
| servers                  | list with all available servers and hosts for the main host                                 |



| loadBalancing            | description                                                                                 |
|--------------------------|---------------------------------------------------------------------------------------------|
| roundRobin               | choose the route to be used according to the percentage value based on the parameter weight |
| lowTimeResponseHeader    | search for the route with the lowest initial response time                                  |
| lowTimeResponseLastByte  | search for the route with the lowest response time of the last byte                         |
| random                   | makes a totally random choice                                                               |
| overLoad                 | throws all the load to the first server until the maximum number of connections is reached  |
| ipv4Hash                 | 000.000.*: Where you have the asterisk goes to a server route                               |
| ipv6Hash                 | ips list                                                                                    |
| hash                     | request uri                                                                                 |



| healthCheck              | description                                                                                 |
|--------------------------|---------------------------------------------------------------------------------------------|
| enabled                  |                                                                                             |
| interval                 |                                                                                             |
| fails                    |                                                                                             |
| passes                   |                                                                                             |
| uri                      |                                                                                             |
| suspendInterval          |                                                                                             |
| match                    |                                                                                             |



| match                    | description                                                                                 |
|--------------------------|---------------------------------------------------------------------------------------------|
| status                   |                                                                                             |
| header                   |                                                                                             |
| body                     |                                                                                             |



| status                   | description                                                                                 |
|--------------------------|---------------------------------------------------------------------------------------------|
| expReg                   |                                                                                             |
| value                    |                                                                                             |
| in                       |                                                                                             |
| notIn                    |                                                                                             |



| header                   | description                                                                                 |
|--------------------------|---------------------------------------------------------------------------------------------|
| key                      |                                                                                             |
| value                    |                                                                                             |



| body                     | description                                                                                 |
|--------------------------|---------------------------------------------------------------------------------------------|
| expReg                   |                                                                                             |



```yaml

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
        bind:
          - host: 127.0.0.1
            ignorePort: true
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

```yaml
```

```yaml
```

```yaml
```

```yaml
```

```yaml
```

How to configure:

```yaml

version: 1.0

capibaribe:

  affluentRiverNameProject:

    listen: :8081
    debugServerEnable: false

    ssl:
      enabled: false

      # TLS 1.0 - 10
      # TLS 1.1 - 11
      # TLS 1.2 - 12
      # SSL 3.0 - 30
      version:
        min: 10
        max: 30
      certificate: /etc/nginx/company.com.crt
      certificateKey: /etc/nginx/company.com.key
      curvePreferences:
        - P256
        - P384
        - P521
        - X25519
      preferServerCipherSuites: false

      # A list of cipher suite IDs that are, or have been, implemented by this
      # package.
      #
      # Taken from https://www.iana.org/assignments/tls-parameters/tls-parameters.xml
      cipherSuites:
        - TLS_RSA_WITH_RC4_128_SHA
        - TLS_RSA_WITH_3DES_EDE_CBC_SHA
        - TLS_RSA_WITH_AES_128_CBC_SHA
        - TLS_RSA_WITH_AES_256_CBC_SHA
        - TLS_RSA_WITH_AES_128_CBC_SHA256
        - TLS_RSA_WITH_AES_128_GCM_SHA256
        - TLS_RSA_WITH_AES_256_GCM_SHA384
        - TLS_ECDHE_ECDSA_WITH_RC4_128_SHA
        - TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA
        - TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA
        - TLS_ECDHE_RSA_WITH_RC4_128_SHA
        - TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA
        - TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA
        - TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA
        - TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256
        - TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256
        - TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
        - TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256
        - TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
        - TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384
        - TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305
        - TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305
        - TLS_FALLBACK_SCSV
      x509:
        certificate: /etc/nginx/company.com.crt
        certificateKey: /etc/nginx/company.com.key

    static:
      - filePath: /docker/static
        serverPath: static

    pygocentrus:
      enabled: true
      delay:
        rate: 0.1
        # time em uS
        min: 2000000
        max: 5000000
      dontRespond: 0.1
      changeLength: 0.0
      changeContent:
        changeRateMin: 0.0
        changeRateMax: 0.0
        changeBytesMin: 1
        changeBytesMax: 10
        rate: 0.1
      deleteContent: 0.1
      changeHeaders:
        - number: 500
          header:
            - key: Content-Type
              value: application/json
          rate: 0.0

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
