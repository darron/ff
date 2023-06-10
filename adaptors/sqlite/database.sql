CREATE TABLE IF NOT EXISTS records (
    id VARCHAR(36) NOT NULL PRIMARY KEY,
    [date] DATE,
    [name] VARCHAR(255),
    city VARCHAR(255),
    province VARCHAR(2),
    licensed BOOLEAN,
    victims INT,
    deaths INT,
    injuries INT,
    suicide BOOLEAN,
    devices_used VARCHAR(255),
    firearms BOOLEAN,
    possessed_legally BOOLEAN,
    warnings TEXT,
    oic_impact BOOLEAN,
    ai_summary TEXT
);

CREATE TABLE IF NOT EXISTS news_stories (
    id VARCHAR(36) NOT NULL PRIMARY KEY,
    record_id VARCHAR(36) NOT NULL,
    [url] TEXT,
    body_text TEXT,
    ai_summary TEXT,
    FOREIGN KEY (record_id) REFERENCES records (id)
);
