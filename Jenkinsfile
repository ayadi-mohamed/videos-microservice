#!/usr/bin/env groovy

def gv

pipeline {
    agent any

    tools {
        go "Go-1.19"
    }

    stages {
        stage("Initialize") {
            steps {
                script {
                    gv = load "script.groovy"
                }
            }
        }

        stage("Increment Version") {
            steps {
                script {
                    gv.incrementDataSeedJobVersion()
                }
            }
        }

        stage("Build Go Binary") {
            steps {
                script {
                    gv.buildGoBinary()
                }
            }
        }

        stage("Build Docker Image") {
            steps {
                script {
                    gv.buildDockerImage()
                }
            }
        }

        stage("Publish Docker Image to Docker Hub") {
            steps {
                script {
                    gv.pushToDockerHub()
                }
            }
        }

        stage("Trivy Scan") {
            steps {
                script {
                    gv.trivyScan()
                }
            }
        }

        stage("Push to Deployment GitHub") {
            steps {
                script {
                    gv.pushToDeploymentGitHub()
                }
            }
        }
    }

    post {
        always {
            emailext body: '''$PROJECT_NAME - Build # $BUILD_NUMBER - $BUILD_STATUS: Check console output at $BUILD_URL to view the results.''',
                    subject: '$PROJECT_NAME - Build # $BUILD_NUMBER - $BUILD_STATUS!', to: 'charradoumayma1@gmail.com',
                    attachmentsPattern: "TrivyReport.html"
        }
        success{
            script {
                gv.gitpush()
            }
        }
    }

}