pipeline {
    agent any

    // Parameters allow you to customize the build at runtime
    parameters {
        string(name: 'GITHUB_REPO', defaultValue: 'git@github.com:devatefops/counter-app.git', description: 'GitHub repository URL')
        string(name: 'BRANCH', defaultValue: 'master', description: 'Branch to clone')
        string(name: 'TARGET_HOST', defaultValue: '13.63.34.25', description: 'App server IP') // optional for deploy later
    }

    environment {
        DEPLOY_USER = "deploy"       // The user on the app server
        WORKSPACE_DIR = "workspace"  // Optional workspace subdirectory
    }

    stages {

        stage('Clone Repository') {
            steps {
                echo "Cloning repository ${params.GITHUB_REPO} branch ${params.BRANCH}"
                git branch: "${params.BRANCH}",
                    url: "${params.GITHUB_REPO}"
                   // credentialsId: 'github-ssh' // SSH credentials added in Jenkins
            }
        }

        // Optional: Stage to create a file in repo (can be removed later)
        stage('Add File (Optional)') {
            steps {
                script {
                    // Create a file in the workspace
                    writeFile file: 'newfile.txt', text: 'Hello from Jenkins pipeline!'

                    echo "File newfile.txt created in workspace."
                }
            }
        }

        // Placeholder for future build stage
        stage('Build') {
            steps {
                echo "Build stage will go here (e.g., Go binary, Docker build, etc.)"
            }
        }

        // Placeholder for future deploy stage
        stage('Deploy') {
            steps {
                echo "Deploy stage will go here (ssh to ${DEPLOY_USER}@${params.TARGET_HOST})"
            }
        }
    }

    post {
        success {
            echo "Pipeline completed successfully!"
        }
        failure {
            echo "Pipeline failed. Check console log for details."
        }
    }
}
