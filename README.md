
# RIoTPoT

<!-- markdownlint-disable MD033 -->
<div align="center" style="box-shadow: 0px 4px 4px rgba(0, 0, 0, 0.25); background-color: #EDF2F4; border-radius: 4px; margin: 2em 0;">
  <img src="docs/assets/aau_logo.png" height="100px;" style="margin: 1em 0; padding: 1em;">
  <div>
    <!-- Workflow status -->
    <a href="https://github.com/aau-network-security/RIoTPot/actions">
        <img alt="GitHub Actions status" src="https://github.com/aau-network-security/RIoTPot/workflows/cyber/badge.svg">
    </a>
    <a href="https://pkg.go.dev/riopot">
        <img src="https://pkg.go.dev/badge/riopot.svg">
    </a>
    <a href="https://goreportcard.com/badge/github.com/aau-network-security/RIoTPot">
        <img src="https://goreportcard.com/badge/github.com/aau-network-security/RIoTPot?style=flat-square">
    </a>
    <a href="">
        <img src="https://img.shields.io/github/release/RIoTPot/project-layout.svg?style=flat-square">
    </a>
  </div>
</div>

RIoTPot is a hybrid interaction honeypot, primarily focused on the emulation IoT and OT protocols, although, it is also capable of emulating other services.
In essence, RIoTPot acts as a proxy service for other honeypots included in the system.
Therefore, you can run any honeypot and other services alongside RIoTPot.
In addition, there is an UI web-application that you can use to manage your routing.

Moreover, RIoTPot comes with multiple low-interaction services ready to use.
Since these services are written as [plugins](https://pkg.go.dev/plugin), they are only supported on Linux; however, you can start RIoTPot without them.
The following table contains the list of services included in RIoTPot by defaul, their internal port, and proxy port.

<div align="center">

| Service | Internal Port | Proxy Port |
| ------- | ------------- | ---------- |
| Coap    | 25683         | 5683       |
| Echo    | 20007         | 7          |
| HTTP    | 28080         | 8080       |
| Modbus  | 20502         | 502        |
| MQTT    | 21883         | 1883       |
| SSH     | 20022         | 22         |
| Telnet  | 20023         | 23         |

</div>

> ## Table of Contents
>
> - [RIoTPoT](#riotpot)
>   - [1. Architecture](#1-architecture)
>   - [2. How to use RIoTPot](#2-how-to-use-riotpot)

## 1. Architecture

The RIoTPot architecture is based on proxy connections to internal and surrounding (or external) services (e.g., honeypots, full-services, containers, remote hosts, etc.).
For this, the honeypot manages a number of user-defined `proxies` that relays connections between services and RIoTPot [^proxies].
This way, RIoTPot can decide how and where to route incomming attacks.
The logic used to determine how to handle the incomming attack is implemented in the form of `middlewares` [^middlewares].
To manage services, middlewares and proxies, RIoTPot ships with a REST API [^api].
The API endpoints can be accessed through your browser at the location `localhost:2022/api/swagger`, showing a [Swagger](https://swagger.io/) interface that allows you to manage the honeypot in real-time.

[^proxies]: **_NOTE:_**
    Internal and surrounding services are not accessible through the Internet.
    Internal services are integrated and only accessible to RIoTPot.
    These services are loaded on-start and can not be deleted, but they can be stopped.
    Surrounding services **must** be in the same network as RIoTPot.
    External services **must** whitelist RIoTPot **only**.

[^middlewares]: **_NOTE:_** Middlewares are currently under development.

[^api]: **_NOTE:_** The RIoTPot API **must not** be exposed to the Internet.
    Regardless, the API currently only accepts connections from the localhost.
    This may be changed in the future, providing a whitelist of hosts and standard authentication.

**Figure 1** shows the RIoTPot architecture, including the two main applications that constitute RIoTPot (RIoTPot itself, and RIoTPot UI) and their components, and a section to enclose external (or adjacent) services.

<div align="center" style="margin: 2em 0">
    <div style="max-width: 60%; text-align: justify; display: flex; flex-direction: column;">
        <img src="docs/assets/architecture.png" style="background-color: #EDF2F4; border-radius: 4px; margin: 1em 0; box-shadow: 0px 4px 4px rgba(0, 0, 0, 0.25);">
        <div>
        <b>Figure 1.</b> RIoTPot Architecture, including the main application, external services and the webapp UI to manage RIoTPot instances.
        </div>
    </div>
</div>

RIoTPot is written in [Golang](https://go.dev/)[^os].
Each RIoTPot instance exposes registered proxies (based on their port) on demand.
To serve a proxy, it **must** have a binded service and the proxy port **must** be available (currently, RIoTPot does not accept multiple services running in the same port).
When a proxy has been binded and served, attackers will be able to send messages to RIoTPot on that port, relying the messages to the binded service and back to the attacker[^reversed].

[^os]: **_NOTE:_** While the base application is interoperable, internal services (plugins) can only be used in [Linux, FreeBSD and macOS environments](https://pkg.go.dev/plugin).
    We plan to overcome this limitation by replacing plugins with micro-services communicating through [gRPC](https://grpc.io/).

[^reversed]: **_NOTE:_** For ethical and security reasons, RIoTPot does not allow unsolicited requests to the outside, i.e., reversed shells and the like are not allowed.

For ease of access, multiple instances of RIoTPot can be managed from the RIoTPot UI webapp.
In addition to managing the proxies registered in each instance, the UI allows you to create, use and edit `profiles`.
Each profile contains a number of proxies named after protocols or other services making a RIoTPot instance resemble, for example, a real-life device.
In few words, profiles speed up the process of setting up and provision a RIoTPot instance with specific configurations.
The UI is written using the React fonrt-end JavaScript library (we use Typescript for this project) and [Recoil](https://recoiljs.org/) state management library.

## 2. How to use RIoTPot

Running RIoTPot is relatively simple.
Overall, you have three options.
**The first** is to download a RIoTPot release; you can either choose to download the latest release, or previous one.
**The second option** is to build the project yourself.
**The last option** is to use the source code to create a Docker container with RIoTPot and some additional applications to enhance the honeypot.

<details open>
    <summary><b>Using a Release Version</b></summary>

> **_Info_:** This guide is meant for users with no special needs, who want a simple out-of-the-box experience.

Each release comes in a folder named `bin` with an executable (also) named `riotpot`, a `plugins` folder filled with multiple services (or low-interaction honeypots), and a folder named `ui` containing the UI server files.
It is important to keep this folder structure for RIoTPot to work as intended.

---

1. First, download the release of your choice from the [releases](https://github.com/aau-network-security/riotpot/releases) page. Choose the one you need for your Operative System (OS).
2. Extact the `riotpot` folder.
3. Run the `riotpot` binary. This will start RIoTPot with the API enabled, all the plugins ready to use, and the UI server.
    - The UI is accessible through the address `localhost:3000` or `local.riotpot.ui`
    - The API is accessible through the address `localhost:2022/api/swagger` or `local.riotpot.hp/api/swagger`

</details>

<details>
    <summary><b>Build it yourself</b></summary>

> **_Info_:** This guide is meant for advanced users confortable in development environments.

<blockquote>
<details>
<summary><b>Requirements</b></summary>

- Golang - Required to build the project
- Node - Required to build the UI

**Optional**:

- Git - Used to download the source code
- Make - To run already-prepared commands

</details>
</blockquote>

---

1. Download the RIoTPot source code from GitHub. Open a console and introduce the following line.

    ```bash
    git clone git@github.com:aau-network-security/riotpot.git
    ```

2. Navigate to the folder in where you have downloaded the RIoTPot source.
3. If you have installed [Make](https://www.gnu.org/software/make/), we have included multiple command helpers to assist you building the project. To put it simple, you can run two simple commands that will build the RIoTPot binary, the plugins (and place them in the right folder), and then serve the UI.

    ```bash
    # Builds RIoTPot and the plugins
    make build-all
    
    # Starts the UI server
    make ui
    ```

</details>

<details>
    <summary><b>Docker (Virtualisation)</b></summary>

> **_Info_:** This guide is meant for advanced users who prefer to use RIoTPot in a virtual environment.

<blockquote>
<details>
<summary><b>Requirements</b></summary>

- Docker - Used to build an image of a RIoTPot instance and UI server.
- Docker-compose - Used to create a single container with a RIoTPot instance, the UI and other applications and services.

</details>
</blockquote>

Some of main the advantages of using this setup are the additional security features with minimal changes to the container configuration and the containers themselves.
For example, we can define separated virtual private networks and overlay networks to hide, sandbox and encapsule RIoTPot and other adjacent services.
In addition, containers allow us to bind services using their docker addres name rather than their IP, which is very convenient.
Lastly, we can spawn and stop separated containers on demand without affecting other services.

On the other hand, virtualisation is arguably more demanding than usign applications on bare-metal.
While a single instance of RIoTPot is relatively lightweight, it is important to consider the shortcomings introduced by virtualisation and hardware emulation (e.g., reponse delays).

> **_Warning_:** Technically speaking, a dedicated attacker may realize that RIoTPot is in fact a honeypot by analysing and comparing the response-time delays introduced by virtualisation to real servers (!!). While this type of honeypot fingerprinting has been studied before, the results for common Internet services are still inconclussive (e.g., HTTP, Telnet and SSH), due to the commoditization of cloud hosting services using virtual machines and detailed server configurations.

The `docker-compose` file includes additional services to enhance the RIoTPot experience.
The following table summarises the list of services and applications packed in this container.

<blockquote>
<details>
<summary><b>Services</b></summary>
<div align="center">

| Service | Image                  | Port | Details                                    |
| ------- | ---------------------- | ---- | ------------------------------------------ |
| MQTT    | eclipse-mosquitto      | 1883 | Mosquito  MQTT Server                      |
| HTTP    | httpd                  | 80   | Regular HTTP Server                        |
| Modbus  | oitc/modbus-server     | 502  | Modbus Server                              |
| OCPP    | ocpp1.6-central-system | 443  | OCPP v1.6 (used in cars charging stations) |

</div>
</details>

<details>
<summary><b>Applications</b></summary>
<div align="center">

| Application | Image           | Details                                                   |
| ----------- | --------------- | --------------------------------------------------------- |
| TCPDump     | kaazing/tcpdump | Packet recorder. It stores network traffic in .pcap files |

</div>
</details>
</blockquote>

---

The container can be setup in three simple steps:

1. Download the RIoTPot source code from GitHub. Open a console and introduce the following line.

    ```bash
    git clone git@github.com:aau-network-security/riotpot.git
    ```

2. Navigate to the folder in where you have downloaded the RIoTPot source.
3. With Docker running: if you have Make installed, run the following command. Otherwise run a docker-compose command using the docker-compose file included in the `build/docker` folder.
    - With make
  
    ```bash
    # With make
    make up
    ```

   - With Docker-compose

    ```bash
    # With docker-compose
    docker-compose -p riotpot -f build/docker/docker-compose.yaml up -d --build
    ```

</details>
