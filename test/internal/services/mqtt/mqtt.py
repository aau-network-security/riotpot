"""
This file is meant to be placed in the `attacker` host.
It simply connects to the mqtt broker honeypot in `riotpot`
and publishes messages to the topic `hello/there`
"""

from paho.mqtt import client as mqtt_client
import random
import time

broker = 'riotpot'
port = 1883
topic = "/hello/there"
client_id = f'python-mqtt-{random.randint(0, 1000)}'
username = "username"
password = "password"

def connect_mqtt():
    def on_connect(client, userdata, flags, rc):
        if rc == 0:
            print("Connected to MQTT Broker!")
        else:
            print("Failed to connect, return code %d\n", rc)
    # Set Connecting Client ID
    client = mqtt_client.Client(client_id)
    client.username_pw_set(username, password)
    client.on_connect = on_connect
    client.connect(broker, port)
    client.loop_start()
    return client

def publish(client):
    msg_count = 0
    while True:
        time.sleep(1)
        msg = f"messages: {msg_count}"
        result = client.publish(topic, msg, qos=2)
        # result: [0, 1]
        status = result[0]
        if status == 0:
            print(f"Send `{msg}` to topic `{topic}`")
        else:
            print(f"Failed to send message to topic {topic}")
        msg_count += 1

connection = connect_mqtt()
connection.loop_start()
publish(connection)
