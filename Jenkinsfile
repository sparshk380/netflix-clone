pipeline {
    agent any

    environment {
        GITHUB_TOKEN = credentials('github-token1')
        IMAGE_TAG = 'unode-onboard-api'
        SOURCE_BRANCH = "${env.CHANGE_BRANCH ?: env.GIT_BRANCH}"
        DOCKERHUB_REPO = 'gaganr31/jenkins'
        K8S_POD_YAML = credentials('k8s-pod-yaml')
    }

    options {
        maskPasswords()
    }

    stages {
        stage('Prepare Kubernetes Agent') {
            steps {
                script {
                    // Write the YAML to a temporary file without echoing content
                    writeFile file: 'k8s-pod.yaml', text: K8S_POD_YAML
                    echo "Kubernetes YAML prepared"
                }
            }
        }

        stage('Run on Kubernetes') {
            agent {
                kubernetes {
                    inheritFrom 'k8s-agent'
                    yamlFile 'k8s-pod.yaml'
                }
            }
            stages {
                stage('Clone Repository and Get Commit SHA') {
                    steps {
                        script {
                            wrap([$class: 'MaskPasswordsBuildWrapper', varPasswordPairs: [[password: env.GITHUB_TOKEN, var: 'GITHUB_TOKEN']]]) {
                                sh """
                                echo "Cloning branch: ${env.SOURCE_BRANCH}"
                                git clone -b ${env.SOURCE_BRANCH} https://${GITHUB_TOKEN}@github.com/Gagan-R31/netflix-clone.git
                                cd netflix-clone
                                """
                            }
                            env.COMMIT_SHA = sh(script: "git rev-parse --short HEAD", returnStdout: true).trim()
                            echo "Commit SHA: ${env.COMMIT_SHA}"
                        }
                    }
                }
                stage('Check Go Installation') {
                    steps {
                        container('kaniko') {
                            sh '''
                            cd netflix-clone
                            which go
                            go version
                            '''
                        }
                    }
                }
                stage('Build Docker Image with Kaniko') {
                    steps {
                        container('kaniko') {
                            sh """
                                cd netflix-clone
                                /kaniko/executor --dockerfile=./Dockerfile \
                                                 --context=. \
                                                 --destination=${DOCKERHUB_REPO}:${IMAGE_TAG}-${env.COMMIT_SHA}
                            """
                        }
                    }
                }
            }
        }
    }
    post {
        always {
            script {
                // Clean up the temporary YAML file
                sh 'rm -f k8s-pod.yaml'
            }
        }
    }
}
