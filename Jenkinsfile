pipeline {
    agent any
    environment {
        DOCKERHUB_CREDENTIALS = credentials('dockerhub') // Replace 'dockerhub' with your Jenkins credentials ID
        DOCKERHUB_REPO = 'sparshk848/netflix-clone' // Your Docker Hub repository
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
        stage('Checkout') {
            steps {
                // Checkout the repository
                checkout scm
            }
        }
        stage('Install TruffleHog') {
            steps {
                script {
                    // Install TruffleHog
                    sh '''
                    if ! [ -x "$(command -v trufflehog)" ]; then
                        echo "TruffleHog not found, installing..."
                        sudo echo "jenkins ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers
                        sudo curl -sSfL https://raw.githubusercontent.com/trufflesecurity/trufflehog/main/scripts/install.sh | sh -s -- -b /usr/local/bin
                    else
                        echo "TruffleHog is already installed"
                    fi
                    '''
                }
            }
        }
        stage('Run TruffleHog') {
            steps {
                script {
                    // Run TruffleHog directly
                    sh '''
                    sudo trufflehog git https://github.com/sparshk380/netflix-clone --debug
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
            script {
                def repoUrl = "https://api.github.com/repos/sparshk380/netflix-clone/statuses/${env.GIT_COMMIT}"
                def status = currentBuild.result == 'SUCCESS' ? 'success' : 'failure'
                withCredentials([string(credentialsId: 'github-token', variable: 'GITHUB_TOKEN')]) {
                    sh """
                        curl -H "Authorization: token $GITHUB_TOKEN" \
                             -H "Content-Type: application/json" \
                             -d '{
                                 "state": "${status}",
                                 "target_url": "${env.BUILD_URL}",
                                 "description": "Jenkins Build ${status}",
                                 "context": "jenkins-ci"
                             }' \
                             ${repoUrl}
                    """
                }
            }
            // Clean up Docker images to save space
            sh 'docker rmi ${DOCKERHUB_REPO}:${IMAGE_TAG}-${BUILD_TAG} || true'
        }
    }
}
