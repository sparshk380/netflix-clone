pipeline {
    agent any
    environment {
        DOCKERHUB_CREDENTIALS = credentials('dockerhub') // Replace 'dockerhub' with your Jenkins credentials ID
        DOCKERHUB_REPO = 'sparshk848/netflix-clone' // Your Docker Hub repository
        IMAGE_TAG = 'netflix-clone' // Image tag, can be changed if needed
        BUILD_TAG = "${env.BUILD_ID}" // Unique tag for each build
        SONARQUBE_SCANNER = tool 'SonarQube Scanner' // Name of your SonarQube Scanner tool installation
        SONARQUBE_SERVER = 'SonarQube' // Name of your SonarQube server configuration
    }
    stages {
        stage('Install Docker') {
            steps {
                script {
                    // Install Docker if not already installed
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
        stage('SonarQube Analysis') {
            steps {
                withSonarQubeEnv(SONARQUBE_SERVER) {
                    script {
                        // Run SonarQube analysis
                        withCredentials([string(credentialsId: 'sonarqube-token', variable: 'SONAR_AUTH_TOKEN')]) {
                            sh """
                            ${SONARQUBE_SCANNER}/bin/sonar-scanner \
                                -Dsonar.projectKey=netflix-clone \
                                -Dsonar.sources=. \
                                -Dsonar.host.url=http://3.83.237.11:9000 \
                                -Dsonar.login=${SONAR_AUTH_TOKEN}
                            """
                        }
                    }
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
                // Post build GitHub status update
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
        }
    }
}
