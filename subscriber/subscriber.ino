#define ESP32

#include <M5Atom.h>
#include <WiFiClient.h>
#include <WiFiClientSecure.h>
#include <PubSubClient.h>
#include <ArduinoJson.h>
#include "Secret.h"

//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
static WiFiClientSecure httpsClient;
static PubSubClient mqttClient(httpsClient);

//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// WiFi接続
void setupWifi() {
  int count = 0;

  Serial.println("Wifi Connecting...");
  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);
  while (WiFi.status() != WL_CONNECTED) {
    count++;
    int status = WiFi.status();
    delay(1000);
    if (count % 10 == 0 && (status == WL_DISCONNECTED || status == WL_CONNECT_FAILED || status == WL_CONNECTION_LOST || status == WL_NO_SSID_AVAIL)) {
      Serial.println("Reconnect...");
      WiFi.disconnect();
      WiFi.begin(WIFI_SSID, WIFI_PASSWORD);
    }
    if (count >= 30) {
      Serial.println("Restart");
      ESP.restart();
    }
    Serial.print(status);
  }
  Serial.println("");
  Serial.println("Wifi Connected!");
}

//-----------------------------------------------------------------------------
// MQTT接続のセットアップ
void setupMqtt() {
  httpsClient.setCACert(ROOT_CA);
  httpsClient.setCertificate(CERTIFICATE);
  httpsClient.setPrivateKey(PRIVATE_KEY);
  mqttClient.setServer(CLOUD_ENDPOINT, CLOUD_PORT);
  mqttClient.setCallback(mqttCallback);
}

//-----------------------------------------------------------------------------
// MQTT接続
void connectMqtt() {
  Serial.println("MQTT Connecting...");
  while (!mqttClient.connected()) {
    if (mqttClient.connect(CLIENT_ID)) {
      mqttClient.subscribe(CLOUD_TOPIC, 0);
    } else {
      delay(1000);
      Serial.print(".");
    }
  }
  Serial.println("MQTT Connected!");
}

//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// MQTTサブスクライブ
void mqttCallback (char* topic, byte* payload, unsigned int length) {
    Serial.print("Received. topic=");
    Serial.println(topic);

    char message[length];
    for (int i = 0; i < length; i++) {
      message[i] = (char)payload[i];
    }

    DynamicJsonDocument doc(200);
    deserializeJson(doc, message);
    String status = doc["status"];
    Serial.print("status:");
    Serial.println(status);

    if (status == "true") {
       M5.dis.fillpix(0xff0000);
    } else {
       M5.dis.fillpix(0x00ff00);
    }
}

//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
void setup()
{
    M5.begin(true, false, true);
    delay(10);

    Serial.begin(9600);
    delay(10);

    M5.dis.fillpix(0xffffff);
    M5.dis.setBrightness(50);

    setupWifi();
    setupMqtt();
}

//-----------------------------------------------------------------------------
void loop()
{
    if (!mqttClient.connected()) {
        connectMqtt();
        M5.dis.fillpix(0x0000ff);
    }
    mqttClient.loop();

    delay(50);
    M5.update();
}
