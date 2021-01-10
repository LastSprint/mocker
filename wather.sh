inotifywait -q -m -r -e modify /mocks |
while read -r filename event; do
  curl mocker:1111/update_models
done