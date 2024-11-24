# Usamos Alpine como base
FROM alpine:latest

# Instalamos las herramientas necesarias: bash, git, curl, docker-cli, tee para el registro del log
RUN apk update && apk add --no-cache \
    bash \
    git \
    curl \
    docker-cli \
    make \
    tee \
    && rm -rf /var/cache/apk/*

# Verificar las versiones de las herramientas instaladas
RUN docker --version && docker buildx version

# Directorio de trabajo donde se guardará el log
WORKDIR /workspace

# Crear un archivo de log para almacenar los resultados
RUN touch process.log

# Script para ejecutar el proceso de construcción e informes
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

# Establecer ENTRYPOINT para ejecutar el script de generación de informe
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
