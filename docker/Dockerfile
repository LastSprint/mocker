FROM ubuntu

RUN apt-get update && \
    apt-get install -y inotify-tools

CMD inotifywait -e attrib /mocks && curl mocker:1111/update_models