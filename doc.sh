#!/bin/bash

echo "Getting ready"

h () {
  sleep 0.5

  echo "Opening in browser..."

  open http://localhost:6060/pkg/github.com/KeyizeBiometry/keyize/
}

h & godoc
