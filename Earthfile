VERSION 0.6

# renovate: datasource=docker depName=aquasec/trivy
ARG TRIVY_VERSION=0.50.1
ARG CONTAINER_BASE=fedora

luet:
    FROM quay.io/luet/base:latest
    SAVE ARTIFACT /usr/bin/luet

govulncheck-scan:
    FROM golang
    RUN go install golang.org/x/vuln/cmd/govulncheck@latest
    COPY +luet/ /
    ARG CONTAINER_IMAGE
    RUN /luet util unpack ${CONTAINER_IMAGE} /tmp/image
    RUN apt-get update && apt-get install -y file
    COPY ./vulncheck.sh /vulncheck.sh
    RUN /vulncheck.sh /tmp/image

govulncheck-report:
    FROM golang
    RUN go install golang.org/x/vuln/cmd/govulncheck@latest
    COPY +luet/ /
    ARG CONTAINER_IMAGE
    RUN /luet util unpack ${CONTAINER_IMAGE} /tmp/image
    RUN apt-get update && apt-get install -y file
    COPY ./vulncheck.reports.sh /vulncheck.sh
    RUN mkdir /reports
    RUN /vulncheck.sh /tmp/image
    SAVE ARTIFACT /reports AS LOCAL build

jq-image:
    FROM ${CONTAINER_BASE}
    RUN dnf update -y
    RUN dnf install -y jq

update-images:
    FROM +jq-image
    COPY . .
    ARG ALL_RELEASES
    ENV ALL_RELEASES=$ALL_RELEASES
    RUN bash update.sh
    SAVE ARTIFACT images.json AS LOCAL images.json
###
### Tools dependencies
###
grype:
    FROM anchore/grype
    SAVE ARTIFACT /grype /grype


trivy:
    ARG TRIVY_VERSION
    FROM aquasec/trivy:$TRIVY_VERSION
    SAVE ARTIFACT /contrib contrib
    SAVE ARTIFACT /usr/local/bin/trivy /trivy

###
### Base container
###
security-container:
    FROM ${CONTAINER_BASE}

    COPY +trivy/trivy /trivy
    COPY +trivy/contrib /contrib
    COPY +grype/grype /grype

###
### Security target scan
###
security-scan:
    ARG CONTAINER_IMAGE
    ARG GOVULNCHECK=false
    ARG LEVEL=critical
    FROM +security-container

    RUN --no-cache /trivy image --scanners vuln ${CONTAINER_IMAGE}
    RUN --no-cache /grype ${CONTAINER_IMAGE} --fail-on ${LEVEL} --only-fixed --verbose
    IF [ $GOVULNCHECK = "true" ]
        BUILD +govulncheck-scan --CONTAINER-IMAGE=${CONTAINER_IMAGE}
    END
###
### Get a report
###
security-report:
    ARG CONTAINER_IMAGE
    ARG GO_VULNCHECK=false
    FROM +security-container

    WORKDIR /build

    RUN /grype ${CONTAINER_IMAGE} --output sarif --add-cpes-if-none --file report.sarif
    RUN /grype ${CONTAINER_IMAGE} --output json --add-cpes-if-none --file report.json
    SAVE ARTIFACT /build/report.sarif report.sarif AS LOCAL build/grype.sarif
    SAVE ARTIFACT /build/report.json report.json AS LOCAL build/grype.json
    RUN /trivy image --skip-dirs /tmp --timeout 30m --format sarif -o report.sarif --no-progress ${CONTAINER_IMAGE}
    RUN /trivy image --skip-dirs /tmp --timeout 30m --format template --template "@/contrib/html.tpl" -o report.html --no-progress ${CONTAINER_IMAGE}
    RUN /trivy image --skip-dirs /tmp --timeout 30m -f json -o results.json --no-progress ${CONTAINER_IMAGE}
    SAVE ARTIFACT /build/report.sarif report.sarif AS LOCAL build/trivy.sarif
    SAVE ARTIFACT /build/report.html report.html AS LOCAL build/trivy.html
    SAVE ARTIFACT /build/results.json results.json AS LOCAL build/trivy.json
    IF [ $GOVULNCHECK = "true" ]
        BUILD +govulncheck-report --CONTAINER-IMAGE=${CONTAINER_IMAGE}
    END
