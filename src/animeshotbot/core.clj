(ns animeshotbot.core
  (:gen-class))

(require '[morse.api :as api])
(require '[morse.polling :as p])
(require '[morse.handlers :refer :all])
(require '[clojure.data.json :as json])
(require '[clj-http.client :as client])

(def token (System/getenv "ANIMESHOTBOT_TG_TOKEN"))

(defn search-shots
  [shot-keyword]
  (let [respBody ((client/get (str "https://as.bitinn.net/api/shots?q=" shot-keyword)
                              {:insecure? true}) :body)]
    (json/read-str respBody :key-fn keyword)))

(defn build-inline-results
  [shots]
  (json/write-str (reduce conj [] (map (fn [shot] {:type "photo"
                                                   :id (shot :sid)
                                                   :photo_url (shot :image_large)
                                                   :thumb_url (shot :image_thumbnail)
                                                   :caption (shot :text)})
                                       shots))))

(defhandler bot-api
  (message msg (let [shots (search-shots (:text msg))]
                 (api/send-text token (get-in msg [:chat :id])
                                (clojure.string/join "\n" (map :text shots)))))
  (inline query (let [shots (search-shots (:query query))]
                  (api/answer-inline token (:id query)
                                     (build-inline-results shots)))))

(defn run-bot []
  (let [channel (p/start token bot-api)]
    (println "start polling updates...")
    (clojure.core.async/<!! channel)))

(defn -main
  [& args]
  (if (clojure.string/blank? token)
    (do (println "Telegram bot token not set please set environment variable ANIMESHOTBOT_TG_TOKEN")
        (System/exit 1))
    (run-bot)))

;;(p/start token bot-api) //for local test
