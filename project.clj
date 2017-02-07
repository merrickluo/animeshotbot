(defproject animeshotbot "0.1.0-SNAPSHOT"
  :description "FIXME: write description"
  :url "http://example.com/FIXME"
  :license {:name "Eclipse Public License"
            :url "http://www.eclipse.org/legal/epl-v10.html"}
  :dependencies [[org.clojure/clojure "1.8.0"]
                 [org.clojure/data.json "0.2.6"]
                 [morse "0.2.3-SNAPSHOT"]]
  :sub ["modules/morse"]
  :main ^:skip-aot animeshotbot.core
  :target-path "target/%s"
  :profiles {:uberjar {:aot :all}})
