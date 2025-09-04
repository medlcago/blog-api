# Blog API

## Quick Start

Use the provided `make` commands to manage the environment.

### Requirements

*   Docker
*   Docker Compose

### Available Commands

*   `make up` - Start the services in detached mode.
*   `make down` - Stop the running services.
*   `make logs` - Follow the logs from the services in real-time.
*   `make purge` - **Warning!** Full cleanup: stops services, removes images, volumes, and orphan containers.

### Environment Configuration

The project requires environment variables to run properly.

1.  Copy the example environment file:
    ```bash
    cp .env.example .env
    ```
2.  Open the newly created `.env` file and provide the actual values for all variables.

### Usage

1.  Clone this repository.
2.  Configure your environment variables as described above.
3.  Run the following command in your terminal:
    ```bash
    make up
    ```
4.  The application will start. To view the logs, run:
    ```bash
    make logs
    ```