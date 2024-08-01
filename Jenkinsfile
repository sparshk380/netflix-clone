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
                    def encodedYaml = K8S_POD_YAML.bytes.encodeBase64().toString()
                    env.ENCODED_K8S_YAML = encodedYaml
                    wrap([$class: 'MaskPasswordsBuildWrapper', varPasswordPairs: [[password: encodedYaml, var: 'ENCODED_K8S_YAML']]]) {
                        echo "Kubernetes YAML prepared and masked"
                    }
                }
            }
        }

        stage('Run on Kubernetes') {
            agent {
                kubernetes {
                    label 'k8s-agent'
                    yaml """
                        ${new String(env.ENCODED_K8S_YAML.decodeBase64())}
                    """
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
}
