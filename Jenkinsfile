pipeline {
    agent any

    environment {
        DOCKERHUB_CREDENTIALS = credentials('dockerhub') // Replace with your actual credentials ID
        DOCKER_IMAGE = 'sparshk848/netflix-clone'
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build Docker Image') {
            steps {
                script {
                    // Build the Docker image and tag it with the build number
                    def dockerImage = docker.build("${DOCKER_IMAGE}:${env.BUILD_NUMBER}")
                }
            }
        }

        stage('Push Docker Image') {
            steps {
                script {
                    // Push the Docker image to Docker Hub
                    docker.withRegistry('https://index.docker.io/v1/', DOCKERHUB_CREDENTIALS) {
                        docker.image("${DOCKER_IMAGE}:${env.BUILD_NUMBER}").push()
                    }
                }
            }
        }
    }

    post {
        success {
            script {
                currentBuild.result = 'SUCCESS'
            }
        }
        failure {
            script {
                currentBuild.result = 'FAILURE'
                error('Build failed, merge to main branch is not allowed.')
            }
        }
    }
}
