FROM bitnami/minideb:latest
COPY opensds-broker /opensds-broker
CMD ["/opensds-broker", "-logtostderr"]
