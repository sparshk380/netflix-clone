pipeline {
    agent any
    environment {
        DOCKERHUB_CREDENTIALS = credentials('dockerhub') // Replace 'dockerhub' with your Jenkins credentials ID
        DOCKERHUB_REPO = 'gaganr31/jenkins' // Your Docker Hub repository
        IMAGE_TAG = 'netflix-clone' // Image tag, can be changed if needed
        BUILD_TAG = "${env.BUILD_ID}" // Unique tag for each build
    }
    stages {
        stage('Install TruffleHog') {
            steps {
                script {
                    // Clean previous TruffleHog installation if any
                    sh '''
                    if [ -x "$(command -v trufflehog)" ]; then
                        sudo rm -f /usr/local/bin/trufflehog
                    fi
                    '''
                    // Download and install TruffleHog
                    sh '''
                    curl -sSfL https://github.com/trufflesecurity/trufflehog/releases/download/v3.8.0/trufflehog_Linux_x86_64.tar.gz -o trufflehog.tar.gz
                    tar -xzf trufflehog.tar.gz
                    sudo mv trufflehog /usr/local/bin/
                    rm trufflehog.tar.gz
                    '''
                }
            }
        }
        stage('Checkout') {
            steps {
                // Checkout the repository
                checkout scm
            }
        }
        stage('Run TruffleHog') {
            steps {
                script {
                    // Run TruffleHog directly
                    sh '''
                    trufflehog git https://github.com/Gagan-R31/Jenkins.git --only-verified || {
                        echo "TruffleHog scan failed"
                        exit 1
                    }
                    '''
                }
            }
        }
        stage('Build Docker Image') {
            steps {
                script {
                    // Build the Docker image
                    sh """
                    docker build -t ${DOCKERHUB_REPO}:${IMAGE_TAG}-${BUILD_TAG} .
                    """
                }
            }
        }
        stage('Push to Docker Hub') {
            steps {
                script {
                    // Log in to Docker Hub
                    sh """
                    echo ${DOCKERHUB_CREDENTIALS_PSW} | docker login -u ${DOCKERHUB_CREDENTIALS_USR} --password-stdin
                    """
                    // Push the Docker image
                    sh """
                    docker push ${DOCKERHUB_REPO}:${IMAGE_TAG}-${BUILD_TAG}
                    """
                }
            }
        }
    }
    post {
        always {
            // Clean up Docker images to save space
            sh 'docker rmi ${DOCKERHUB_REPO}:${IMAGE_TAG}-${BUILD_TAG} || true'
        }
    }
}
