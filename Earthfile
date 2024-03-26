VERSION 0.6

# renovate: datasource=docker depName=aquasec/trivy
ARG TRIVY_VERSION=0.49.1

security-container:
    FROM fedora

    COPY +trivy/trivy /trivy
    COPY +trivy/contrib /contrib
    COPY +grype/grype /grype


grype:
    FROM anchore/grype
    SAVE ARTIFACT /grype /grype

###
### Security target scan
###
trivy:
    ARG TRIVY_VERSION
    FROM aquasec/trivy:$TRIVY_VERSION
    SAVE ARTIFACT /contrib contrib
    SAVE ARTIFACT /usr/local/bin/trivy /trivy

security-scan:
    ARG CONTAINER_IMAGE
    FROM +security-container

    RUN /grype ${CONTAINER_IMAGE} --fail-on critical
    RUN /trivy image --scanners vuln ${CONTAINER_IMAGE}

security-report:
    ARG CONTAINER_IMAGE
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