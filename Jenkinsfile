pipeline {
    agent {
        kubernetes {
            label 'k8s-agent'
            yaml '''
            apiVersion: v1
            kind: Pod
            metadata:
              labels:
                some-label: some-label-value
            spec:
              containers:
              - name: jnlp
                image: jenkins/inbound-agent
                args: ['$(JENKINS_SECRET)', '$(JENKINS_NAME)']
              - name: kaniko
                image: gcr.io/kaniko-project/executor:debug
                command:
                - /busybox/sh
                tty: true
                volumeMounts:
                - name: kaniko-secret
                  mountPath: /kaniko/.docker
              - name: golang
                image: golang:1.21
                command:
                - cat
                tty: true
              volumes:
              - name: kaniko-secret
                secret:
                  secretName: kaniko-secret
                  items:
                    - key: .dockerconfigjson
                      path: config.json
            '''
        }
    }
    environment {
        DOCKERHUB_CREDENTIALS = credentials('dockerhub') // Replace 'dockerhub' with your Jenkins credentials ID
        DOCKERHUB_REPO = 'gaganr31/jenkins' // Your Docker Hub repository
        IMAGE_TAG = 'my-app' // Image tag, can be changed if needed
        BUILD_TAG = "${env.BUILD_ID}" // Unique tag for each build
        IMAGE_NAME = "${DOCKERHUB_REPO}:${IMAGE_TAG}-${BUILD_TAG}" // Full image name
    }
    stages {
        stage('Build Docker Image with Kaniko') {
            steps {
                container('kaniko') {
                    script {
                        sh '''
                        /kaniko/executor --dockerfile=Dockerfile \
                                         --context=${WORKSPACE} \
                                         --destination=${IMAGE_NAME} \
                                         --tarPath=/kaniko/image.tar \
                                         --cleanup
                        '''
                    }
                }
            }
        }
        stage('Load and Test Docker Image') {
            steps {
                container('golang') {
                    script {
                        sh '''
                        # Load the Docker image from the tar file
                        docker load -i /kaniko/image.tar
                        # Run the Docker container
                        docker run -d --name test-container ${IMAGE_NAME}
                        # Run tests on the container
                        docker exec test-container go test -v ./...
                        # Stop and remove the container
                        docker stop test-container
                        docker rm test-container
                        '''
                    }
                }
            }
        }
        stage('Push Docker Image to Docker Hub') {
            steps {
                container('kaniko') {
                    script {
                        sh '''
                        /kaniko/executor --dockerfile=Dockerfile \
                                         --context=${WORKSPACE} \
                                         --destination=${IMAGE_NAME} \
                                         --cleanup
                        '''
                    }
                }
            }
        }
    }
    post {
        always {
            script {
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
        }
    }
}
