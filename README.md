# Capibaribe

> This project is in test.

This is the preliminary version of my version of the reverse proxy, designed to be simple and fast. 
The initial idea is to allow a single machine to receive multiple small sites, and can host then at a low cost.

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

version: '1.0'

reverseProxy:

  config:                                         # Server configuration
    listenAndServer: :8080                        # Local server address and port. Format: server:port
    outputConfig: true                            # Print the configuration before server start. Format: true/false
    staticServer: true                            # Enable a static files server with files contained into same server machine. Format: true/false
    staticFolder:                                 # Static file server config
      - folder: /docker/static/                   # Files folders to server
        serverPath: static                        # Main server path. Example: 'yourdomain.com/static'
      - folder: /docker/another_static/           # Files folders to server
        serverPath: another_static                # Main server path. Example: 'yourdomain.com/static'

  proxy:                                          # Servidor proxy. Redireciona o endereço para um novo servidor

    # To the example below, run the docker command to enable ghost blog at port 2368
    # docker run -d --name ghost-blog-1 -p 2368:2368 ghost
    blog:                                         # Name from server 1 log file
      host: blog.localhost:8080                   # Income host address
      server:                                     # Alternatives servers list
        - name: docker 1 - ok                     # Name from alternative content server 1
          host: http://localhost:2368             # Host from alternative content server 1
        - name: docker 2 - error                  # Name from alternative content server 2
          host: http://localhost:2369             # Host from alternative content server 2
        - name: docker 3 - error                  # Name from alternative content server 3
          host: http://localhost:2370             # Host from alternative content server 3

    # To the example below, run the docker command to enable ghost blog at port 2378
    # docker run -d --name ghost-blog-2 -p 2378:2368 ghost
    blog_2:                                       # Name from server 2 log file
      host: blog2.localhost:8080                  # Income host address
      server:                                     # Alternatives servers list
        - name: docker 4 - ok                     # Name from alternative content server 4
          host: http://localhost:2378             # Host from alternative content server 4
        - name: docker 5 - error                  # Name from alternative content server 5
          host: http://localhost:2379             # Host from alternative content server 5
        - name: docker 6 - error                  # Name from alternative content server 6
          host: http://localhost:2380             # Host from alternative content server 6

```
