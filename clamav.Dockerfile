FROM uselagoon/commons as commons

FROM clamav/clamav:1.4.0

COPY --from=commons /lagoon /lagoon
COPY --from=commons /bin/fix-permissions /bin/ep /bin/docker-sleep /bin/wait-for /bin/

RUN apk add --no-cache tzdata

RUN sed -i "s/^LogFile /# LogFile /g" /etc/clamav/clamd.conf && \
    sed -i "s/^#LogSyslog /LogSyslog /g" /etc/clamav/clamd.conf && \
    sed -i "s/^UpdateLogFile /# UpdateLogFile /g" /etc/clamav/freshclam.conf && \
    sed -i "s/^#LogSyslog /LogSyslog /g" /etc/clamav/freshclam.conf

USER root

RUN fix-permissions /var/lib/clamav

ENTRYPOINT [ "/init-unprivileged" ]
