pipeline {
    agent any

    environment {
        GITHUB_TOKEN = credentials('github-token1')
        IMAGE_TAG = 'unode-onboard-api'
        SOURCE_BRANCH = "${env.CHANGE_BRANCH ?: env.GIT_BRANCH}"
        DOCKERHUB_REPO = 'gaganr31/jenkins'
        K8S_POD_YAML = credentials('k8s-pod-yaml')
    }

    stages {
        stage('Create Kubernetes Agent') {
            steps {
                script {
                    podTemplate = readYaml text: env.K8S_POD_YAML
                    podTemplate.spec.containers[0].image = 'jenkins/inbound-agent'
                    podTemplate.spec.containers[0].args = ['$(JENKINS_SECRET)', '$(JENKINS_NAME)']
                    
                    def encodedYaml = podTemplate.toString().bytes.encodeBase64().toString()
                    
                    withCredentials([string(credentialsId: 'k8s-pod-yaml', variable: 'YAML_CRED')]) {
                        env.ENCODED_YAML = encodedYaml
                    }
                }
            }
        }

        stage('Run on Kubernetes') {
            agent {
                kubernetes {
                    label 'k8s-agent'
                    yaml """
                        ${new String(env.ENCODED_YAML.decodeBase64())}
                    """
                }
            }
            stages {
                stage('Clone Repository and Get Commit SHA') {
                    steps {
                        script {
                            sh """
                            echo "Cloning branch: ${env.SOURCE_BRANCH}"
                            git clone -b ${env.SOURCE_BRANCH} https://${GITHUB_TOKEN}@github.com/Gagan-R31/netflix-clone.git
                            cd netflix-clone
                            """
                            env.COMMIT_SHA = sh(script: "git rev-parse --short HEAD", returnStdout: true).trim()
                            echo "Commit SHA: ${env.COMMIT_SHA}"
                        }
                    }
                }
                stage('Check Go Installation') {
                    steps {
                        container('kaniko') {
                            script {
                                sh '''
                                cd netflix-clone
                                which go
                                go version
                                '''
                            }
                        }
                    }
                }
                stage('Build Docker Image with Kaniko') {
                    steps {
                        container('kaniko') {
                            script {
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
}
