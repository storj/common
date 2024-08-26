pipeline {
    agent {
        docker {
            label 'main'
            image 'storjlabs/ci:latest'
            alwaysPull true
            args '-u root:root --cap-add SYS_PTRACE -v "/tmp/gomod":/go/pkg/mod'
        }
    }
    options {
          timeout(time: 26, unit: 'MINUTES')
    }
    environment {
        NPM_CONFIG_CACHE = '/tmp/npm/cache'
    }
    stages {
        stage('Build') {
            steps {
                checkout scm
                sh 'git restore-mtime'

                sh 'mkdir -p .build'
                sh 'cp go.mod .build/go.mod.orig'

                sh 'service postgresql start'

                dir(".build") {
                    sh 'cockroach start-single-node --insecure --store=type=mem,size=2GiB --listen-addr=localhost:26256 --http-addr=localhost:8086 --cache 512MiB --max-sql-memory 512MiB --background'
                }
            }
        }

        stage('Verification') {
            parallel {
                stage('Lint') {
                    steps {
                        sh 'check-copyright'
                        sh 'check-large-files'
                        sh 'check-imports ./...'
                        sh 'check-peer-constraints'
                        sh 'storj-protobuf --protoc=$HOME/protoc/bin/protoc lint'
                        sh 'storj-protobuf --protoc=$HOME/protoc/bin/protoc check-lock'
                        sh 'check-atomic-align ./...'
                        sh 'check-errs ./...'
                        sh 'check-monkit ./...'
                        sh './scripts/check-dependencies.sh'
                        sh 'staticcheck ./...'
                        sh './scripts/arm-staticcheck.sh'
                        sh 'golangci-lint --config /go/ci/.golangci.yml -j=2 run'
                        sh 'check-mod-tidy -mod .build/go.mod.orig'
                        sh 'go-licenses check ./...'
                        sh 'check-downgrades'
                        sh 'go.min vet ./...'
                    }
                }

                stage('Tests') {
                    environment {
                        STORJ_TEST_POSTGRES = 'postgres://postgres@localhost/teststorj?sslmode=disable'
                        STORJ_TEST_COCKROACH = 'cockroach://root@localhost:26256/testcockroach?sslmode=disable'
                        COVERFLAGS = "${ env.BRANCH_NAME != 'main' ? '' : '-coverprofile=.build/coverprofile -coverpkg=./...'}"
                    }
                    steps {
                        sh 'use-ports -from 1024 -to 10000 &'
                        sh 'cockroach sql --insecure --host=localhost:26256 -e \'create database testcockroach;\''
                        sh 'psql -U postgres -c \'create database teststorj;\''
                        sh 'go test -parallel 4 -p 6 -vet=off $COVERFLAGS -timeout 20m -json -race ./... | tee .build/tests.json | xunit -out .build/tests.xml'
                        sh 'check-clean-directory'
                    }

                    post {
                        always {
                            sh script: 'cat .build/tests.json | tparse -all -slow 100', returnStatus: true
                            archiveArtifacts artifacts: '.build/tests.json'
                            junit '.build/tests.xml'

                            script {
                                if(fileExists(".build/coverprofile")){
                                    sh script: 'filter-cover-profile < .build/coverprofile > .build/clean.coverprofile', returnStatus: true
                                    sh script: 'gocov convert .build/clean.coverprofile > .build/cover.json', returnStatus: true
                                    sh script: 'gocov-xml  < .build/cover.json > .build/cobertura.xml', returnStatus: true
                                    cobertura coberturaReportFile: '.build/cobertura.xml'
                                }
                            }
                        }
                    }
                }

                stage('Go Compatibility') {
                    steps {
                        sh './scripts/check-cross-compile.sh'
                    }
                }
            }
        }
    }

    post {
        always {
            sh "chmod -R 777 ." // ensure Jenkins agent can delete the working directory
            deleteDir()
        }
    }
}
