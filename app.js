const http = require('http');
const fs = require('fs');
const mysql = require('mysql');

const server = http.createServer((req, res) => {
    if (req.url === '/run') {
        console.log("Be Happy! Script is running...");
        const paragraphs = openJson();
        const connection = createConnection();
        createTable(connection);
        insertParagraphs(connection, paragraphs);
        connection.end();
        console.log("All done! get some rest...");

        // Send a response to indicate the script has been executed
        res.writeHead(200, { 'Content-Type': 'text/plain' });
        res.end('JavaScript script executed successfully');
    }
});

function openJson() {
    const fileContent = fs.readFileSync('paragraphs.json', 'utf8');
    return JSON.parse(fileContent);
}

function createConnection() {
    return mysql.createConnection({
        host: 'tai.liara.cloud',
        port: 30983,
        user: 'root',
        password: 'seh1iWk2MvRySPWhUHp01m1N',
        database: 'trusting_merkle'
    });
}

function createTable(connection) {
    connection.connect();
    connection.query(`
        CREATE TABLE IF NOT EXISTS random_words (
            id INT AUTO_INCREMENT PRIMARY KEY,
            paragraph TEXT
        )`);
}

function insertParagraphs(connection, paragraphs) {
    paragraphs.forEach(paragraph => {
        connection.query(`
            INSERT INTO random_words (paragraph) VALUES (?)`, [paragraph.paragraph]);
    });
}

const PORT = 3000;
const HOST = 'nodejs-paragraph';

server.listen(PORT, HOST, () => {
    console.log(`Server running at http://${HOST}:${PORT}/`);
});
