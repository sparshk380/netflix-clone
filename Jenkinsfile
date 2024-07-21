pipeline {
    agent any
    environment {
        IMAGE_TAG = 'my-app' // Image tag, can be changed if needed
        BUILD_TAG = "${env.BUILD_ID}" // Unique tag for each build
    }
    stages {
        stage('Checkout') {
            steps {
                // Checkout the repository
                checkout scm
            }
        }
        stage('Install Cosign') {
            steps {
                script {
                    // Install Cosign if not already installed
                    sh '''
                        if ! [ -x "$(command -v cosign)" ]; then
                            echo "Cosign not found, installing..."
                            curl -O -L "https://github.com/sigstore/cosign/releases/latest/download/cosign-linux-amd64"
                            mv cosign-linux-amd64 /usr/local/bin/cosign
                            chmod +x /usr/local/bin/cosign
                        else
                            echo "Cosign is already installed"
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
                            curl -sSfL https://raw.githubusercontent.com/trufflesecurity/trufflehog/main/scripts/install.sh | sh -s -- -v -b /usr/local/bin
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
                    // Run TruffleHog to scan the repository
                    sh 'trufflehog git https://github.com/Gagan-R31/Jenkins --debug'
                }
            }
        }
        stage('Install Go') {
            steps {
                script {
                    // Install Go if not already installed
                    sh '''
                    if ! [ -x "$(command -v go)" ]; then
                        echo "Go not found, installing..."
                        curl -LO https://golang.org/dl/go1.21.1.linux-amd64.tar.gz
                        tar -C /usr/local -xzf go1.21.1.linux-amd64.tar.gz
                        echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.profile
                    else
                        echo "Go is already installed"
                    fi
                    '''
                    // Ensure the new Go binary is in the PATH
                    sh 'export PATH=$PATH:/usr/local/go/bin'
                }
            }
        }
        stage('Build Docker Image') {
            steps {
                script {
                    // Build the Docker image
                    sh "docker build -t ${DOCKERHUB_REPO}:${IMAGE_TAG}-${BUILD_TAG} ."
                }
            }
        }
        stage('Test') {
            steps {
                // Run Go tests
                sh 'go test -v ./...'
            }
        }
        // stage('Push to Docker Hub') {
        //     steps {
        //         script {
        //             // Log in to Docker Hub and push the Docker image
        //             sh "echo ${DOCKERHUB_CREDENTIALS_PSW} | docker login -u ${DOCKERHUB_CREDENTIALS_USR} --password-stdin"
        //             sh "docker push ${DOCKERHUB_REPO}:${IMAGE_TAG}-${BUILD_TAG}"
        //         }
        //     }
        // }
    }
    post {
        always {
            node('kubernetes') {
                script {
                    // Report build status to GitHub
                    def repoUrl = "https://api.github.com/repos/Gagan-R31/Jenkins/statuses/${env.GIT_COMMIT}"
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
}
