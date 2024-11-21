pipeline {
    agent any  // Utiliza un agente genérico; cámbialo si es necesario

    environment {
        // Configuración del entorno
        DOCKER_CLI_EXPERIMENTAL = 'enabled'
        IMAGE_NAME = 'mi-imagen'  // Nombre de la imagen Docker
        DOCKER_REGISTRY = 'mi-registro'  // Nombre del registro Docker (puede ser Docker Hub u otro)
        DOCKER_TAG = "latest1.0"  // El tag de la imagen (puedes usar otro si lo prefieres)
    }

    stages {
        stage('Preparar Docker Buildx') {
            steps {
                script {
                    // Verifica que Docker y Buildx están instalados
                    sh 'docker --version'
                    sh 'docker buildx version'

                    // Crea un nuevo builder usando Buildx
                    sh 'docker buildx create --use'
                }
            }
        }

        stage('Construir Imagen') {
            steps {
                script {
                    // Construcción de la imagen Docker con soporte multi-plataforma
                    sh """
                        docker buildx build --platform linux/amd64,linux/arm64 -t $DOCKER_REGISTRY/$IMAGE_NAME:$DOCKER_TAG .
                    """
                }
            }
        }

        stage('Iniciar Sesión en Docker Registry') {
            steps {
                script {
                    // Inicia sesión en Docker (solo si es necesario empujar la imagen)
                    // Se recomienda usar credenciales de Jenkins para esto.
                    withCredentials([usernamePassword(credentialsId: 'docker-hub-credentials', usernameVariable: 'DOCKER_USERNAME', passwordVariable: 'DOCKER_PASSWORD')]) {
                        sh """
                            docker login -u \$DOCKER_USERNAME -p \$DOCKER_PASSWORD
                        """
                    }
                }
            }
        }

        stage('Empujar Imagen a Registro') {
            when {
                branch 'main'  // Solo empuja en la rama 'main'
            }
            steps {
                script {
                    // Empuja la imagen construida al registro Docker
                    sh """
                        docker push $DOCKER_REGISTRY/$IMAGE_NAME:$DOCKER_TAG
                    """
                }
            }
        }
    }

    post {
        always {
            // Realiza cualquier acción adicional al finalizar el pipeline
            echo 'Pipeline completado.'
        }
        success {
            echo 'Imagen construida y subida correctamente.'
        }
        failure {
            echo 'Hubo un error en el pipeline.'
        }
    }
}
