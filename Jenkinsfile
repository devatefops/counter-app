pipeline {
    agent any

    environment {
        // --- Configuration ---
        // The name for your Docker image.
        DOCKER_IMAGE_NAME = "counter-app"
        // The URL of your Docker registry (e.g., 'yourdockerhubusername').
        DOCKER_REGISTRY = "devatefops"
        // The Jenkins credentials ID for your Docker registry.
        DOCKER_CREDENTIALS_ID = "docker-registry-credentials"
        // The Jenkins credentials ID for SSH access to your target server.
        TARGET_SERVER_SSH_CREDENTIALS_ID = "target-server-ssh-key"
        // The user and hostname for your target server (e.g., 'ubuntu@your-server.com').
        TARGET_SERVER = "ubuntu@16.16.159.60"
        // The name of the container on the target server.
        CONTAINER_NAME = "production-counter-app"
    }

    stages {
        stage('Checkout') {
            steps {
                // Get the source code from your Git repository.
                checkout scm
            }
        }

        stage('Build Docker Image') {
            steps {
                script {
                    // Build the Docker image using the Dockerfile in the repository.
                    def dockerImage = docker.build("${DOCKER_IMAGE_NAME}:${env.BUILD_NUMBER}", ".")
                }
            }
        }

        stage('Push to Docker Registry') {
            steps {
                script {
                    // Log in to the Docker registry, tag the image, and push it.
                    docker.withRegistry("https://${DOCKER_REGISTRY}", DOCKER_CREDENTIALS_ID) {
                        def imageNameWithRegistry = "${DOCKER_REGISTRY}/${DOCKER_IMAGE_NAME}"

                        // Push the image with the build number tag.
                        docker.image("${DOCKER_IMAGE_NAME}:${env.BUILD_NUMBER}").push()

                        // Also push the image with the 'latest' tag.
                        docker.image("${DOCKER_IMAGE_NAME}:${env.BUILD_NUMBER}").push("latest")
                    }
                }
            }
        }

        stage('Deploy to Target Server') {
            steps {
                // Use sshagent to securely connect to the remote server.
                sshagent(credentials: [TARGET_SERVER_SSH_CREDENTIALS_ID]) {
                    // The script to run on the remote server.
                    sh """
                        ssh -o StrictHostKeyChecking=no ${TARGET_SERVER} '''
                            echo "--- Deploying new version ---"
                            
                            # Pull the latest image from the registry.
                            docker pull ${DOCKER_REGISTRY}/${DOCKER_IMAGE_NAME}:latest
                            
                            # Stop and remove the old container if it exists.
                            docker stop ${CONTAINER_NAME} || true
                            docker rm ${CONTAINER_NAME} || true
                            
                            # Run the new container.
                            docker run -d --name ${CONTAINER_NAME} -p 8080:8080 ${DOCKER_REGISTRY}/${DOCKER_IMAGE_NAME}:latest
                            
                            echo "--- Deployment complete ---"
                        '''
                    """
                }
            }
        }
    }
}