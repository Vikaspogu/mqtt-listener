# MQTT Listener for frigate snapshots

[![GitHub last commit](https://img.shields.io/github/last-commit/vikaspogu/mqtt-listener?color=purple&style=flat-square)](https://github.com/vikaspogu/mqtt-listener/commits/master) [![Docker Build Status](https://github.com/vikaspogu/mqtt-listener/workflows/push_latest/badge.svg)](https://github.com/vikaspogu/mqtt-listener/actions)

[frigate](https://github.com/blakeblackshear/frigate) is a realtime object detection for cameras, which publishes images to mqtt topic. This repo subscribes to topic snapshot topic and saves it to disk.

