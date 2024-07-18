pipeline {
    agent any

    environment {
        // Set the Docker Hub credentials ID stored in Jenkins
        DOCKERHUB_CREDENTIALS = credentials('docker-hub-credentials-id') // Replace 'docker-hub-credentials-id' with the ID of your Docker Hub credentials in Jenkins
        // Define the Docker image name to build and push
        DOCKER_IMAGE = 'sparshk380/netflix-clone' // Replace 'sparshk380' with your Docker Hub username
    }

    stages {
        stage('Checkout') {
            steps {
                // Checkout the code from the specified Git repository and branch
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
                // Perform any additional actions on successful build
                currentBuild.result = 'SUCCESS'
            }
        }
        failure {
            script {
                // Set the build result to FAILURE and prevent merge to main branch
                currentBuild.result = 'FAILURE'
                error('Build failed, merge to main branch is not allowed.')
            }
        }
    }
}
