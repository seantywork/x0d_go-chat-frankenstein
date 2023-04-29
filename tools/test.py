import asyncio
from websockets.sync.client import connect

def hello():
    with connect("ws://localhost:8889/serverside-message") as websocket:
        websocket.send("imtheserver")
        resp_message = websocket.recv()

        print(resp_message)

        if resp_message == 'READY':

            websocket.send("b03b606740340d6f50128c0a81c40390")
            resp_message = websocket.recv()


        else :
            print("failed")
            return
        

        print(resp_message)


hello()