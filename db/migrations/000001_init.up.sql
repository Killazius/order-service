CREATE TABLE IF NOT EXISTS orders (
                        id TEXT PRIMARY KEY,
                        item TEXT NOT NULL,
                        quantity INTEGER NOT NULL
);