#!/bin/bash

HOME="`pwd`"
$HOME/minio server --config-dir=$HOME/config $@ $HOME/data