(defproject animeshotbot "0.1.0-SNAPSHOT"
  :description "Telegram bot for animeshot"
  :url "https://github.com/merrickluo/animeshotbot"
  :license {:name "WTFPL"
            :url "http://www.wtfpl.net/"}
  :dependencies [[org.clojure/clojure "1.8.0"]
                 [org.clojure/data.json "0.2.6"]
                 [morse "0.2.5-SNAPSHOT"]]
  :plugins [[lein-sub "0.3.0"]]
  :sub ["modules/morse"]
  :main ^:skip-aot animeshotbot.core
  :target-path "target/%s"
  :profiles {:uberjar {:aot :all}})
