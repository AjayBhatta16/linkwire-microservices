# initialize function directory
cd ..\functions
mkdir $args[0]

# setup new project
cd $args[0]
go mod init myfunction
Copy-Item ..\..\generator\defaults\* .\ -Recurse

# add shared utilities and models
Copy-Item ..\..\shared\utilities\* .\ -Recurse
Copy-Item ..\..\shared\models\* .\ -Recurse

# install default dependencies
go get cloud.google.com/go/firestore
go get github.com/golang-jwt/jwt/v5
go get cloud.google.com/go/functions
go get github.com/GoogleCloudPlatform/functions-framework-go

# get latest version of shared utilities and models
go get github.com/AjayBhatta16/linkwire-golang-shared@latest

# create deploy script
$template = Get-Content ..\..\generator\templates\deploy.sh
$template = $template -replace "{{FUNCTION_NAME}}", $args[0]
$template | Out-File deploy.sh -Encoding utf8

# create CI/CD Pipeline
$template = Get-Content ..\..\generator\templates\ci-cd.yaml
$template = $template -replace "{{FUNCTION_NAME}}", $args[0]
$template | Out-File ..\..\..\..\.github\workflows\deploy-$($args[0]).yaml -Encoding utf8