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
        DOCKERHUB_CREDENTIALS = credentials('dockerhub')
        DOCKERHUB_REPO = 'gaganr31/jenkins'
        IMAGE_TAG = 'my-app'
        BUILD_TAG = "${env.BUILD_ID}"
    }
    
    stages {
        stage('Build and Test') {
            steps {
                container('kaniko') {
                    script {
                        sh '''
                        /kaniko/executor --dockerfile=Dockerfile \
                        --context=${WORKSPACE} \
                        --destination=${DOCKERHUB_REPO}:${IMAGE_TAG}-${BUILD_TAG} \
                        --no-push \
                        --tarPath=/workspace/image.tar \
                        --cleanup
                        '''
                    }
                }
                container('golang') {
                    script {
                        sh '''
                        # Extract the image contents
                        mkdir -p /tmp/app
                        tar -xf /workspace/image.tar -C /tmp/app

                        # Find the app directory and run tests
                        APP_DIR=$(find /tmp/app -type d -name "app" | head -n 1)
                        if [ -z "$APP_DIR" ]; then
                            echo "App directory not found"
                            exit 1
                        fi
                        cd $APP_DIR
                        go test -v ./...
                        '''
                    }
                }
            }
        }
        
        stage('Push Docker Image') {
            steps {
                container('kaniko') {
                    script {
                        sh '''
                        /kaniko/executor --dockerfile=Dockerfile \
                        --context=${WORKSPACE} \
                        --destination=${DOCKERHUB_REPO}:${IMAGE_TAG}-${BUILD_TAG} \
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
                    -d '{ "state": "${status}", "target_url": "${env.BUILD_URL}", "description": "Jenkins Build ${status}", "context": "jenkins-ci" }' \
                    ${repoUrl}
                    """
                }
            }
        }
    }
}
