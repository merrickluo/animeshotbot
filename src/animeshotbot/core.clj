(ns animeshotbot.core
  (:gen-class))

(require '[morse.api :as api])
(require '[morse.polling :as p])
(require '[morse.handlers :refer :all])
(require '[clojure.data.json :as json])
(require '[clj-http.client :as client])

(def token (System/getenv ))

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

(defn -main
  "I don't do a whole lot ... yet."
  [& args]
  (p/start token bot-api))
