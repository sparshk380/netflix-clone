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
                - /busybox/cat
                tty: true
                volumeMounts:
                - name: kaniko-secret
                  mountPath: /kaniko/.docker
              - name: golang
                image: golang:1.16
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
    }
    stages {
        stage('Install Go') {
            steps {
                container('golang') {
                    script {
                        sh '''
                        # Go should already be installed in golang:1.16
                        go version
                        '''
                    }
                }
            }
        }
        stage('Fetch Dockerfile') {
            steps {
                container('kaniko') {
                    script {
                        sh '''
                        curl -L -o /workspace/Dockerfile https://github.com/Gagan-R31/Jenkins/raw/feat-1/Dockerfile
                        '''
                    }
                }
            }
        }
        stage('Build and Push Docker Image with Kaniko') {
            steps {
                container('kaniko') {
                    script {
                        sh '''
                        /kaniko/executor --dockerfile=/workspace/Dockerfile \
                                         --context=/workspace \
                                         --destination=$DOCKERHUB_REPO:$IMAGE_TAG-$BUILD_TAG \
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
