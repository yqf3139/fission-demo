from flask import request
from flask import current_app
from time import sleep
import json
import re

def main():
    current_app.logger.info("Received request")
    body = request.get_data().decode("utf-8")
    category = 'cat'
    name = 'none'
    try:
        payload = json.loads(body)
        m = re.search(r'-(\w+)-', payload['Key'])
        name = payload['Key'].split('/')[1]
        current_app.logger.info(payload["Key"])
        if m:
            category = m.group(1)
    except Exception as e:
        current_app.logger.error(e)
        pass

    sleep(1)
    return json.dumps({'err':'', 'payload': {
        'category': category,
        'name': name,
    }})
