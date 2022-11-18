pipeline {
    agent any

    stages {
        stage('socks 5 project pull') {
            steps {
                checkout([$class: 'GitSCM', branches: [[name: '*/master']], extensions: [], userRemoteConfigs: [[credentialsId: '4b0b0294-702f-4370-8afe-3bc462ae66d4', url: 'https://gitee.com/lai-zongji/sockes5.git']]])
            }
        }
        stage('build') {
            steps {
                sh 'go build cmd/main.go'
            }
        }
    }
}
