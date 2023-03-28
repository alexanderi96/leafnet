# Leafnet - A Virtual Genealogy Tree Application

Leafnet is a virtual genealogy tree application written in Go, allowing registered users to create and manage family trees by adding family members or individuals and specifying their parent relationships. The application connects to a Neo4j database to store the information and leverages the 3d-force-graph library to create a visual representation of the family relationships.

## Features

- User registration and authentication
- Session management with Gorilla sessions
- Add family members or individuals and specify parent relationships
- Display past and future anniversaries
- Visual representation of family relationships using [vasturiano/3D-force-graph](https://github.com/vasturiano/3d-force-graph)

## Configuration

Leafnet uses a YAML configuration file with the following structure:

```
neo4j_endpoint: "neo4j-db-endpoint"
neo4j_port: "neo4j-db-port"
neo4j_schema: "neo4j-db-schema"
neo4j_username: "neo4j-db-username"
neo4j_password: ".neo4j-db-password"
leafnet_port: "leafnet-port"
```

## Tech Stack

- Go (Golang)
- Gorilla sessions
- Neo4j Database
- HTML templates (Go)
- JavaScript
- Bootstrap 5
- DataTables
- Select2
- 3d-force-graph library

## Getting Started

Clone the repository and navigate to the project directory.

Install the required dependencies:

```
go get
```

Compile the application:

```
go build
```

Fill in the configuration details in the YAML configuration file and put it under `~/.config/leafnet/config.yaml` or specify the path to it while starting the application.

Run it:

```
./leafnet -s my-session-key -c config/file/path
```

Open your browser and navigate to http://localhost:8080 to access the application.

## Contributing

Please feel free to submit issues, fork the repository, and send pull requests for any improvements or bug fixes.
License

This project is licensed under the MIT License.