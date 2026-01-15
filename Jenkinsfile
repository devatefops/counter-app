// Defines the entire CI/CD process as a Declarative Pipeline.
pipeline {
    // Specifies that the pipeline can run on any available Jenkins agent.
    agent any

    // Defines parameters that can be set by the user when triggering the pipeline.
    // This makes the pipeline flexible without needing to change the code.
    parameters {
        string(name: 'GITHUB_REPO', defaultValue: 'https://github.com/devatefops/counter-app.git', description: 'GitHub repository URL')
        string(name: 'BRANCH', defaultValue: 'master', description: 'Branch to clone')
        string(name: 'TARGET_HOST', defaultValue: '13.63.34.25', description: 'App server IP') // optional for deploy later
    }

    // Sets environment variables that are available throughout the pipeline.
    environment {
        DEPLOY_USER = "deploy"       // The user on the app server
        WORKSPACE_DIR = "workspace"  // Optional workspace subdirectory
    }

    // Contains the sequential stages that make up the workflow.
    stages {

        stage('Clone Repo on App Server') {
            steps {
                echo "Cloning repo ${params.GITHUB_REPO} on host ${params.TARGET_HOST}"
                // Use the sshagent wrapper to securely provide SSH credentials for the connection.
                // 'server-app-ssh' is the ID of the credential stored in Jenkins.
                sshagent(credentials: ['server-app-ssh']) {
                    sh """
                        # Execute a git clone command on the remote TARGET_HOST.
                        ssh -o StrictHostKeyChecking=no ${DEPLOY_USER}@${params.TARGET_HOST} \
                        "git clone --branch ${params.BRANCH} ${params.GITHUB_REPO} ${WORKSPACE_DIR}/counter-app"
                    """
                }
            }
        }

        stage('Build Docker Image') {
            steps {
                echo "Building Docker image on ${params.TARGET_HOST}"
                sshagent(credentials: ['server-app-ssh']) {
                    // A multi-line script is sent to the remote host using a 'heredoc' (<<'EOF').
                    sh """
                        ssh -o StrictHostKeyChecking=no ${DEPLOY_USER}@${params.TARGET_HOST} /bin/bash -s <<'EOF'
                            # Navigate to the repository directory cloned in the previous stage
                            cd ${WORKSPACE_DIR}/counter-app

                            # Build the Docker image, tagging it with the Jenkins build number
                            # Using ${BUILD_NUMBER} is a great practice for versioning artifacts.
                            echo "Building Docker image counter-app:${BUILD_NUMBER}"
                            docker build -t counter-app:${BUILD_NUMBER} .
EOF
                    """
                }
            }
        }

        stage('Deploy') {
            steps {
                echo "Deploying container on ${params.TARGET_HOST}"
                sshagent(credentials: ['server-app-ssh']) {
                    // Execute a deployment script on the remote host.
                    sh """
                        ssh -o StrictHostKeyChecking=no ${DEPLOY_USER}@${params.TARGET_HOST} /bin/bash -s <<'EOF'
                            # Stop and remove the old container to prevent port conflicts
                            # '|| true' ensures the command doesn't fail if the container doesn't exist.
                            echo "Stopping and removing existing container..."
                            docker stop counter-app-container || true
                            docker rm counter-app-container || true

                            # Run the new container from the image built in the previous stage
                            # -d runs the container in detached mode, --name gives it a predictable name.
                            echo "Running new container from image counter-app:${BUILD_NUMBER}"
                            docker run -d --name counter-app-container -p 8080:8080 counter-app:${BUILD_NUMBER}
EOF
                    """
                }
            }
        }
    }

    // The post block defines actions that run at the end of the pipeline's execution.
    post {
        success {
            echo "Pipeline completed successfully!"
        }
        failure {
            echo "Pipeline failed. Check console log for details."
        }
    }
}
