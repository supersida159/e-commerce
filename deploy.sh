APP_NAME=e-commerce
DEPLOY_CONNECTION=root@167.172.75.249

echo "Downloading packages..."
npm install

echo "Compiling..."
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app

echo "Building Docker image..."
docker build -t ${APP_NAME} -f ./Dockerfile .

echo "Saving Docker image to tar file..."
docker save -o ${APP_NAME}.tar ${APP_NAME}

echo "Deploying application..."
scp -o StrictHostKeyChecking=no ./${APP_NAME}.tar ${DEPLOY_CONNECTION}:~
ssh -o StrictHostKeyChecking=no ${DEPLOY_CONNECTION} 'bash -s' < ./deploy/stg.sh

echo "Cleaning up local tar file..."
rm -f ./${APP_NAME}.tar

echo "Deployment complete."
