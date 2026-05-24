# initialize function directory
cd ..\functions
mkdir $args[0]

# setup new project
cd $args[0]
npm init
Copy-Item ..\..\generator\defaults\* .\ -Recurse

# install default dependencies
npm install @google-cloud/functions-framework

# create deploy script
$template = Get-Content ..\..\generator\templates\deploy.sh
$template = $template -replace "{{FUNCTION_NAME}}", $args[0]
$template | Out-File deploy.sh -Encoding utf8

# create CI/CD Pipeline
$template = Get-Content ..\..\generator\templates\ci-cd.yaml
$template = $template -replace "{{FUNCTION_NAME}}", $args[0]
$template | Out-File ..\..\..\..\.github\workflows\deploy-$($args[0]).yaml -Encoding utf8