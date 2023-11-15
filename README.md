# Alexa Clone Project

"Welcome to the Alexa Clone project! This project replicates some of the functionalities of Amazon's Alexa using a microservices architecture, with services implemented in GoLang."

## Table of Contents

- [Overview](#overview)
- [Folder Structure](#folder-structure)
- [Getting Started](#getting-started)
- [Microservices](#microservices)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

## Overview

"This project is an Alexa Clone developed in GoLang, employing a microservices architecture. The goal is to provide a modular and scalable system for implementing voice-controlled functionalities similar to Amazon Alexa."

## Folder Structure

"The project is organized as follows:

- code/
  - alexa.go
  - alpha.go
  - go.mod
  - go.sum
  - stt.go
  - text.xml
  - tts.go
- testing/
  - alexatest.sh
  - alphatest.sh
  - stttest.sh
  - ttstest.sh
  - ttstest.sh
 
-code: Contains the source code for each microservice.

-testing: Houses the testing files corresponding to each microservice.

## Getting Started
To get started with the Alexa Clone project, follow these steps:

1. Clone the repository:

```
bash
git clone https://github.com/your-username/alexa-clone.git
```

## Microservices
The Alexa Clone project comprises the following microservices:

### Service 1: Speech to Text Service
A microservice that listens for a request containing a base64-encoded audio file of the question. It utilizes the Microsoft Speech to Text API to convert the audio to a text question and responds with the text.

### Service 2: Answering Service
Listens for a request containing a question and uses the Wolfram Alpha API to answer the question. It responds with the text answer.

### Service 3: Text to Speech Service
A microservice that does the opposite of Service 1. It converts text to speech, creating an audio response.

### Service 4: Alexa Service
Listens for a base64-encoded question, pipelines the other three services together, and responds with the text answer.

## Testing
Tests for each microservice can be found in the testing/ directory.

## License
This project is licensed under the GNU General Public License v3.0 - see the LICENSE file for details.


 
