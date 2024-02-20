import json
import base64
import time 
from selenium import webdriver
import mysql.connector
import logging
from http.server import BaseHTTPRequestHandler, HTTPServer

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

class RequestHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path == "/run":

            logging.info(f"Be Happy! Script is running...")

            driver, websites         = pre_process()
            db_cursor, db_connection = db_conf()
            take_and_save_screenshots(websites, driver, db_cursor, db_connection)
            # db_cursor.close()
            # db_connection.close()
            # driver.quit()
            logging.info(f"All done! Time to rest.")

            # Send a response to indicate the script has been executed
            self.send_response(200)
            self.send_header('Content-type', 'text/plain')
            self.end_headers()
            self.wfile.write(b'Python script executed successfully')

def pre_process():

    # Load website URLs from JSON file
    with open('websites.json') as f:
        websites = json.load(f)

    # Setup driver
    chrome_options = webdriver.ChromeOptions()
    chrome_options.set_capability('browserless:token', '0iUMlmDN_NB6YWDQ1RF')
    chrome_options.add_argument("--no-sandbox")
    chrome_options.add_argument("--headless")
    try:
        driver = webdriver.Remote(
            command_executor='https://headless.liara.run/webdriver',
            options=chrome_options) 
        
        logging.info("WebDriver connected successfully.")
        return driver, websites
    
    except Exception as e:
        logging.error(f"Error connecting to WebDriver: {str(e)}")
        return None, None

def db_conf():
    db_connection = mysql.connector.connect(
        host="tai.liara.cloud",
        port=30983,
        user="root",
        password="seh1iWk2MvRySPWhUHp01m1N",
        database="trusting_merkle"
    )
    db_cursor = db_connection.cursor()
    db_cursor.execute("SHOW TABLES LIKE 'screenshots'")
    if not db_cursor.fetchone():
        db_cursor.execute("CREATE TABLE screenshots (id INT AUTO_INCREMENT PRIMARY KEY, url VARCHAR(255), screenshot MEDIUMTEXT)")
    
    return db_cursor, db_connection

# Function to take screenshots and save to database
def take_and_save_screenshots(websites, driver, db_cursor, db_connection):
    try:
        # Iterate over websites
        for website in websites:
            url = website['url']
            driver.get(url)

            # Take screenshot
            screenshot = driver.get_screenshot_as_png()

            # Convert screenshot to base64
            screenshot_base64 = base64.b64encode(screenshot).decode('utf-8')

            # Save base64 screenshot to MariaDB
            query = "INSERT INTO screenshots (url, screenshot) VALUES (%s, %s)"
            values = (url, screenshot_base64)
            db_cursor.execute(query, values)
            db_connection.commit()
            logging.info(f"Screenshot saved for {url}")
    except Exception as e:
        logging.error(f"Error occurred: {str(e)}")

def run(server_class=HTTPServer, handler_class=RequestHandler, host='python-script', port=80):
    server_address = (host, port)
    httpd = server_class(server_address, handler_class)
    logging.info(f'Starting server on {host}:{port}')
    httpd.serve_forever()

if __name__ == "__main__":
    run()