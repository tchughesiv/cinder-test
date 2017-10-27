FROM centos:centos7

LABEL name="tchughesiv/cinder-test" \
      maintainer="tchughesiv@gmail.com"
ENV APP_ROOT=/opt/app-root
ENV PATH=${APP_ROOT}/bin:${PATH} HOME=${APP_ROOT}
COPY bin/ cinder-test ${APP_ROOT}/bin/
RUN chmod -R u+x ${APP_ROOT}/bin && \
    chgrp -R 0 ${APP_ROOT} && \
    chmod -R g=u ${APP_ROOT} /etc/passwd

USER 10001
WORKDIR ${APP_ROOT}
ENTRYPOINT [ "uid_entrypoint" ]
CMD run
