# MQTT Panel Broker

This repository is part of a two-repository project that provides an easy-to-use panel for managing information 
processed by an MQTT broker. This repository, **MQTT Panel Broker**, focuses on setting up the broker functionality, 
with MQTT business logic handled via the [Mochi MQTT library](https://github.com/mochi-mqtt/server). The **MQTT
Panel** in the companion repository serves as the control panel, allowing users to view, manage, and control the flow 
of MQTT information in real time.

- This repository: **MQTT Panel Broker** - [mqtt-panel-broker](https://github.com/qreidt/mqtt-panel-broker)
- Control Panel repository: **MQTT Panel** - [mqtt-panel](https://github.com/qreidt/mqtt-panel)

## Table of Contents
- [Project Overview](#project-overview)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [License](#license)

## Project Overview

The **MQTT Panel Broker** is written in Golang and uses the [Mochi MQTT](https://github.com/mochi-mqtt/server) library
to handle all core MQTT operations. It leverages Mochi’s hooks for data collection on each action performed by the 
broker, such as connection, subscription, publishing, and disconnection events.

An authentication layer, implemented through hooks on top of Mochi, is also integrated into this broker. This layer 
authenticates incoming client connections via an API from the **MQTT Panel** repository, ensuring that only authorized
clients can connect and interact with the broker.

## Features

- **Easy Integration with Control Panel**: Works in tandem with the **MQTT Panel** control interface to provide real-time data management and user-friendly controls.
- **Customizable Hook-Based System**: Uses Mochi's hook-based architecture to allow fine-grained data collection and action handling.
- **Authentication Layer**: Ensures secure client access by authenticating new connections through the API provided by the **MQTT Panel**.
- **Efficient MQTT Operations**: Mochi MQTT library provides efficient and robust MQTT business logic.

## Installation

To set up and run the **MQTT Panel Broker**:

1. **Clone this repository**:
   ```bash
   git clone https://github.com/qreidt/mqtt-panel-broker.git
   cd mqtt-panel-broker
   ```
   
2. **Install Dependencies:** Ensure you have Go installed. Then install any required packages:
   ```bash
   go mod tidy
   ```
3. **Build the Project:**
   ```bash
   go build -o mqtt-panel-broker
   ```
   
## Usage
Run the broker executable to start the MQTT broker:
```bash
./mqtt-panel-broker
```

## Configuration
You may configure the broker through environment variables or configuration files. Key configurations include:

- **MQTT Settings**: Configure MQTT topics, QoS levels, and other parameters as needed.
- **Authentication API URL:** Set the URL for the MQTT Panel API used to authenticate client connections.
- **Logging Options:** Enable or disable detailed logs for debugging and performance monitoring. 
For full configuration options, refer to the code comments and configuration files provided in this repository.

## License
This project is licensed under the [MIT License](https://opensource.org/license/mit).