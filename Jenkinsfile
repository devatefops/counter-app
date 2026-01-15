pipeline {
    agent any

    parameters {
        string(name: 'GITHUB_REPO', defaultValue: 'git@github.com:YOUR_ORG/YOUR_REPO.git', description: 'GitHub repo')
        string(name: 'BRANCH', defaultValue: 'main', description: 'Branch to clone')
    }

    stages {
        stage('Clone Repository') {
            steps {
                git branch: "${params.BRANCH}",
                    url: "${params.GITHUB_REPO}",
                    credentialsId: 'github-ssh'
            }
        }
    }
}
