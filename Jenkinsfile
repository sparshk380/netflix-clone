pipeline {
    agent {
        // This assumes Jenkins has access to Docker
        docker { image 'docker:stable' }
    }
    environment {
        // Set Docker Hub credentials and other necessary variables
        DOCKERHUB_CREDENTIALS = credentials('dockerhub') // Replace 'dockerhub' with your Jenkins credentials ID
        DOCKERHUB_REPO = 'gaganr31/jenkins' // Replace with your Docker Hub repository
        IMAGE_NAME = 'netflix-clone' // Replace with your desired image name
    }
    stages {
        stage('Checkout') {
            steps {
                // Checkout the repository
                checkout scm
            }
        }

        stage('Build Docker Image') {
            steps {
                script {
                    // Build the Docker image
                    sh """
                    docker build -t ${env.DOCKERHUB_REPO}/${env.IMAGE_NAME}:${env.BUILD_ID} .
                    """
                }
            }
        }
        stage('Push to Docker Hub') {
            steps {
                script {
                    // Log in to Docker Hub
                    sh """
                    echo ${env.DOCKERHUB_CREDENTIALS_PSW} | docker login -u ${env.DOCKERHUB_CREDENTIALS_USR} --password-stdin
                    """
                    // Push the Docker image
                    sh """
                    docker push ${env.DOCKERHUB_REPO}/${env.IMAGE_NAME}:${env.BUILD_ID}
                    """
                }
            }
        }
    }
    post {
        always {
            // Clean up Docker images to save space
            sh 'docker rmi ${env.DOCKERHUB_REPO}/${env.IMAGE_NAME}:${env.BUILD_ID} || true'
        }
    }
}
