all: test
build:
	go build -o ./bin/eventmaker -v

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./bin/eventmaker

clean:
	go clean
	rm -f ./bin/eventmaker

run:
	go run *.go --project=${GCP_PROJECT} --region=us-central1 --registry=demo-reg \
		--device=demo-device-1 --ca=root-ca.pem --key=device1-private.pem \
		--src=model-client --freq=2s --metric=friction --range=0.01-2.00

certs:
	openssl req -x509 -nodes -newkey rsa:2048 \
			-keyout device1-private.pem \
			-out device1-public.pem \
			-days 365 \
			-subj "/CN=demo"
	curl https://pki.google.com/roots.pem > ./root-ca.pem


setup:
	gcloud pubsub topics create demo-iot-events

	gcloud iot registries create demo-reg \
		--project=${GCP_PROJECT} \
		--region=us-central1 \
		--event-notification-config=topic=demo-iot-events

	gcloud iot devices create demo-device-1 \
		--project=${GCP_PROJECT} \
		--region=us-central1 \
		--registry=demo-reg \
		--public-key path=device1-public.pem,type=rsa-x509-pem

cleanup:
	gcloud iot devices delete demo-device-1 \
		--project=${GCP_PROJECT} \
		--registry=demo-reg \
		--region=us-central1

	gcloud iot registries delete demo-reg \
		--project=${GCP_PROJECT} \
		--region=us-central1

	gcloud beta pubsub topics delete demo-iot-events