pipeline {
    agent any
    environment {
        DOCKERHUB_CREDENTIALS = credentials('dockerhub') // Replace 'dockerhub' with your Jenkins credentials ID
        DOCKERHUB_REPO = 'gaganr31/jenkins' // Your Docker Hub repository
        IMAGE_TAG = 'netflix-clone' // Image tag, can be changed if needed
        BUILD_TAG = "${env.BUILD_ID}" // Unique tag for each build
    }
    stages {
        stage('Install Docker') {
            steps {
                script {
                    // Install Docker
                    sh '''
                    if ! [ -x "$(command -v docker)" ]; then
                        echo "Docker not found, installing..."
                        curl -fsSL https://get.docker.com -o get-docker.sh
                        sh get-docker.sh
                        sudo usermod -aG docker $USER
                        sudo systemctl start docker
                        sudo chmod 666 /var/run/docker.sock
                    else
                        echo "Docker is already installed"
                    fi
                    '''
                }
            }
        }
        stage('Install Python and Pip') {
            steps {
                script {
                    // Install Python and pip if not installed
                    sh '''
                    if ! [ -x "$(command -v python3)" ]; then
                        echo "Python not found, installing..."
                        sudo apt-get update
                        sudo apt-get install -y python3 python3-pip
                    else
                        echo "Python is already installed"
                    fi
                    if ! [ -x "$(command -v pip3)" ]; then
                        echo "pip not found, installing..."
                        sudo apt-get install -y python3-pip
                    else
                        echo "pip is already installed"
                    fi
                    '''
                }
            }
        }
        stage('Install TruffleHog') {
            steps {
                script {
                    // Install TruffleHog if not already installed
                    sh '''
                    if ! [ -x "$(command -v trufflehog)" ]; then
                        echo "TruffleHog not found, installing..."
                        pip3 install truffleHog
                    else
                        echo "TruffleHog is already installed"
                    fi
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
                    // Run TruffleHog
                    sh '''
                    trufflehog git https://github.com/Gagan-R31/Jenkins.git --branch Dev
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
