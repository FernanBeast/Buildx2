#!/bin/bash

# Definimos el archivo de log
LOG_FILE="/workspace/process.log"

# Función para registrar la salida de los comandos en el log
log_result() {
    echo "$(date) - $1" | tee -a $LOG_FILE
}

# Paso 1: Verificar Docker y Buildx
log_result "Iniciando la construcción de la imagen..."
docker --version &>> $LOG_FILE
docker buildx version &>> $LOG_FILE

# Paso 2: Construir la imagen Docker
log_result "Construyendo la imagen Docker..."
docker buildx build --platform linux/amd64,linux/arm64 -t $DOCKER_REGISTRY/$IMAGE_NAME:$DOCKER_TAG . &>> $LOG_FILE
if [ $? -eq 0 ]; then
    log_result "Imagen Docker construida con éxito."
else
    log_result "Error al construir la imagen Docker."
    exit 1
fi

# Paso 3: Empujar la imagen al registro
log_result "Iniciando sesión en Docker Registry..."
docker login -u $DOCKER_USERNAME -p $DOCKER_PASSWORD &>> $LOG_FILE
if [ $? -eq 0 ]; then
    log_result "Sesión iniciada correctamente en Docker Registry."
else
    log_result "Error al iniciar sesión en Docker Registry."
    exit 1
fi

log_result "Empujando imagen al registro..."
docker push $DOCKER_REGISTRY/$IMAGE_NAME:$DOCKER_TAG &>> $LOG_FILE
if [ $? -eq 0 ]; then
    log_result "Imagen empujada correctamente al registro."
else
    log_result "Error al empujar la imagen al registro."
    exit 1
fi

# Finalizar el informe
log_result "Proceso completado."

# Analizar el archivo de log para determinar el resultado final
if grep -q "Error" $LOG_FILE; then
    log_result "Proceso fallido. Verifique los errores en el log."
    exit 1
else
    log_result "Proceso completado con éxito."
fi
