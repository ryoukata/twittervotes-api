FROM scratch

ARG MONGO_PASS

ENV MONGO_HOST=twitter-votes-mongodb MONGO_PORT=27017 MONGO_DB=ballots MONGO_USER=mongo MONGO_PASS=${MONGO_PASS} MONGO_SOURCE=ballots

COPY twittervotes-api .

ENTRYPOINT ["./twittervotes-api"]