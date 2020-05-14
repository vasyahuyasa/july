#!/bin/bash

../kubectl create secret generic july-credentials --from-file=credentials.json=./credentials.json --from-file=token.json=./token.json