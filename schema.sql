CREATE TABLE graph (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

CREATE TABLE node (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    graph_id VARCHAR(36) NOT NULL,
    FOREIGN KEY (graph_id) REFERENCES graph(id)
);

CREATE TABLE edge (
    id VARCHAR(36) PRIMARY KEY,
    from_node VARCHAR(36) NOT NULL,
    to_node VARCHAR(36) NOT NULL,
    cost DECIMAL(10, 2),
    FOREIGN KEY (from_node) REFERENCES node(id),
    FOREIGN KEY (to_node) REFERENCES node(id)
);