(ns animeshotbot.core
  (:gen-class))

(require '[morse.api :as api])
(require '[morse.polling :as p])
(require '[morse.handlers :refer :all])
(require '[clojure.data.json :as json])
(require '[clj-http.client :as client])

(def token (System/getenv "ANIMESHOTBOT_TG_TOKEN"))

(defn search-shots
  [shot-keyword shot-offset]
  (try
    (let [respBody ((client/get (str "https://as.bitinn.net/api/shots?q="
                                     (java.net.URLEncoder/encode shot-keyword)
                                     "&page="
                                     (inc shot-offset))
                                {:insecure? true}) :body)]
      (json/read-str respBody :key-fn keyword))
    (catch Exception e (prn e))))

(defn build-inline-results
  [shots]
  (json/write-str (into [] (map (fn [shot]
                                  {:type "photo"
                                   :id (shot :sid)
                                   :photo_url (shot :image_large)
                                   :thumb_url (shot :image_thumbnail)
                                   :caption (shot :text)
                                   :photo_width 1200
                                   :photo_height 800 })
                                ;; {:type "article"
                                ;;  :id (shot :sid)
                                ;;  :title (shot :text)
                                ;;  :input_message_content {:message_text (str  (shot :text) "\n" (shot :image_large))}
                                ;;  :thumb_url (shot :image_thumbnail)}
                                shots))))

(defhandler bot-api
  (message msg (let* [shots (search-shots (:text msg) 0)
                      text (if (empty? shots)
                             "No Results"
                             (clojure.string/join "\n" (map :text shots)))]
;;                 (prn shots)
;;                 (println msg)
                 (api/send-text token (get-in msg [:chat :id]) text)))

  (inline query (let* [offset (if (clojure.string/blank? (:offset query)) 0 (read-string (:offset query)))
                       text (:query query)
                       shots (if (nil? text) [] (search-shots text offset))]
;;                  (println query)
;;                  (println offset)
                  (try (api/answer-inline token (:id query)
                                          (if (not (empty? shots))
                                            {:next_offset (inc offset)})
                                          (build-inline-results shots))
                       (catch Exception e (prn e))))))

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

;;for local test
;;(p/start token bot-api)
