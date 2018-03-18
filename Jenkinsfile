pipeline {
    agent any

    environment {
        REPO_SERVER = 'repo.youkebox.be'
        REPO_PATH   = "/var/vhosts/repo/${env.GIT_BRANCH}"
        NAME        = 'octaaf'
        VERSION     = '0.1.0'
        DESCRIPTION = 'A Go Telegram bot'
        ARCH        = 'x86_64'
    }

    stages {
        stage('Build') {
            steps {
                sh 'docker run --rm -v "$PWD":/go/src/octaaf -w /go/src/octaaf golang:1.10 /bin/bash -c "go get -v && go build -v"'
            }
        }

        stage('Package') {
            steps {
                sh "fpm -s dir -t rpm \\
                        --name ${NAME} \\
                        --description ${DESCRIPTION} \\
                        --version ${VERSION} \\
                        --architecture ${ARCH} \\
                        --chdir ${TMPDIR} \\
                        --iteration ${env.BUILD_NUMBER} \\
                        .; \\"
            }
        }

        stage('Deploy') {
            steps {
                sh "scp octaaf-*.rpm root@${REPO_SERVER}:${REPO_PATH}/packages/"
                sh "ssh root@${REPO_SERVER} 'createrepo --update ${REPO_PATH}'"
            }
        }
    }
}