def incrementDataSeedJobVersion(){
    echo "Incrementing the Application Version"
    def currentVersion = sh(script: "grep 'const version' app.go | awk '{print \$NF}' | tr -d '\"'", returnStdout: true).trim()
    // Incrementing the Version
    def newVersion = incrementVersion(currentVersion)
    // Updating the Version in the Source Code
    sh "sed -i 's/const version = \"$currentVersion\"/const version = \"$newVersion\"/' app.go"
    // Commit the Changes
    sh "git remote add oumayma git@github.com:ayadi-mohamed/videos-microservice.git"
    sh "git checkout main"
    sh "git commit -am 'Increment Version to $newVersion'"
    // Setting the New Version as an Environment Variable for Later Use
    env.IMAGE_VERSION = newVersion
}

def incrementVersion(currentVersion) {
    def versionParts = currentVersion.split("\\.")
    def newPatchVersion = versionParts[2].toInteger() + 1
    return "${versionParts[0]}.${versionParts[1]}.$newPatchVersion"
}

def buildGoBinary() {
    echo "Compiling and Building the Application..."
    sh "go build -o videos-microservice-${IMAGE_VERSION}"
}

def buildDockerImage() {
    echo "Building the Docker Image..."
    sh "docker build -t oumaymacharrad/videos-microservice:${IMAGE_VERSION} ."
}

def pushToDockerHub() {
    withCredentials([usernamePassword(credentialsId: "Docker-Hub-Credentials", passwordVariable: "PASS", usernameVariable: "USER")]) {
        echo "Pushing the Docker Image to Docker Hub..."
        sh "echo $PASS | docker login -u $USER --password-stdin"
        sh "docker push oumaymacharrad/videos-microservice:${IMAGE_VERSION}"
    }
}

def trivyScan(){
    echo "Running Trivy Security Scan..."
    sh "trivy image --format template --template '@/usr/local/share/trivy/templates/html.tpl' -o TrivyReport.html oumaymacharrad/videos-microservice:${IMAGE_VERSION} --scanners vuln"
}

def pushToDeploymentGitHub() {
    echo "Pushing to Deployment GitHub..."
}

def gitpush(){
    // Push the Changes to GitHub
    sshagent (credentials: ["Private-Key"]) {
        sh "git push oumayma main"
        sh "git remote remove oumayma"
    }
}

return this