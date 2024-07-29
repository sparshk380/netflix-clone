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
              - name: kubectl
                image: joshendriks/alpine-k8s
                command;
                - /bin/cat
              - name: kaniko
                image: gcr.io/kaniko-project/executor:debug
                command:
                - /busybox/sh
                tty: true
                volumeMounts:
                - name: kaniko-secret
                  mountPath: /kaniko/.docker
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
                                         --cleanup
                        '''
                    }
                }
            }
        }
        stage('Run Tests on Docker Image') {
            steps {
                container('kubectl')
                script {
                    // Define Kubernetes Pod for testing
                    def testPodYaml = """
                    apiVersion: v1
                    kind: Pod
                    metadata:
                      name: test-pod
                    spec:
                      containers:
                      - name: test-container
                        image: ${IMAGE_NAME}
                        command: ['go', 'test', '-v', './...']
                      restartPolicy: Never
                    """
                    
                    // Create the Pod
                    sh "echo '${testPodYaml}' | kubectl apply -f -"
                    
                    // Wait for the Pod to complete
                    sh 'kubectl wait --for=condition=complete pod/test-pod --timeout=300s'
                    
                    // Get the test results
                    sh 'kubectl logs pod/test-pod'
                    
                    // Clean up the test Pod
                    sh 'kubectl delete pod/test-pod'
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
