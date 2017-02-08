FROM clojure:latest

ADD . /usr/src/app
WORKDIR /usr/src/app

RUN git submodule update --init && \
		lein sub install && \
		lein uberjar

CMD java -jar target/uberjar/animeshotbot-0.1.0-SNAPSHOT-standalone.jar
