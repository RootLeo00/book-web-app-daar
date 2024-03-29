import os
import requests
import logging
import psycopg2
from datetime import datetime
from utils import get_word_occurrence
import json 
import compute_jaccard
import sys 

DB_FILE_PATH = "../backend/db.sqlite3"
GUTENDEX_URL = "https://gutendex.com/books"
MIME_TYPE = "text"
MAX_PAGES = os.environ.get("MAX_PAGES")
MAX_WORDS = os.environ.get("MAX_WORDS")

if MAX_PAGES == "":
    MAX_PAGES = "60"

if MAX_WORDS == "":
    MAX_WORDS = "10000"

DATABASE_URL = os.environ.get("DATABASE_URL")

## connect to the SQLite database
def connect_to_database():
     # try to connect to PostgreSQL first
    try:
        conn = psycopg2.connect(dsn=DATABASE_URL)
        print("Successfully connected to the PostgreSQL database.")
        return conn
    except psycopg2.Error as e:
        print(f"PostgreSQL Database error: {e}")
        print("Trying to connect to SQLite database...")

    # Fallback to SQLite3
    # try:
    #     conn = sqlite3.connect(DB_FILE_PATH)
    #     print("Successfully connected to the SQLite database.")
    #     return conn
    # except sqlite3.Error as e:
    #     print(f"SQLite Database error: {e}")
    #     return None


## construct the URL for the current page
def construct_url(api_endpoint, mime_type, page):
    return f"{api_endpoint}?mime_type={mime_type}&page={page}"


## fetch all book data from a page
def fetch_and_store_data(conn):
    page = 1
    while page < int(MAX_PAGES):
        books_url = construct_url(GUTENDEX_URL, MIME_TYPE, page)
        response = requests.get(books_url)
        print(f"Fetching books from: {books_url} ...")
        if response.status_code != 200:
            break
        json_data = response.json()
        books = json_data['results']
        process_and_store_books(books, conn)
        page += 1


## get book information
def process_and_store_books(books, conn):
    cursor = conn.cursor()

    ## manipulation of JSON data
    for book in books:
        url_text = None 
        for key in book['formats'].keys():
            if "text/plain" in key:
                url_text = book['formats'][key]
    
        if url_text:
            url_text = url_text.replace('.zip', '.txt')
        else:
            print("cannot get the context of the book, skipping...")
            continue

        author = 'None'
        if book['authors']:
            author = book['authors'][0]['name']
        
        text_response = requests.get(url_text)

        if text_response.status_code == 200:
            text_in_words = text_response.text.split()
            book_text= ' '.join(text_in_words[:int(MAX_WORDS)])
        else: 
            book_text = book['title'] + " " + " ".join(book['subjects'])

        current_time = datetime.now().strftime("%Y-%m-%d %H:%M:%S.%f")

        ## create a new book instance
        book_instance = {
            'book_id': book['id'],
            'created_at': current_time,
            'updated_at': current_time,
            'title': book['title'],
            'author': author,
            'language': book['languages'][0],
            'text': book_text,      # this is not the actual text
            'image_url': book['formats'].get('image/jpeg', None),
            'c_rank': 0.0,
            'occurrence': 0
        }

        ## save the book instance to database
        insert_book_query = """INSERT INTO books (created_at, updated_at, book_id, title, author, language, text, image_url, c_rank, occurrence)
                          VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s)"""
        cursor.execute(insert_book_query, (book_instance['created_at'], book_instance['updated_at'], book_instance['book_id'],
                                        book_instance['title'], book_instance['author'], book_instance['language'], 
                                        book_instance['text'], book_instance['image_url'], book_instance['c_rank'], 
                                        book_instance['occurrence']))
        
        word_occurrence = get_word_occurrence(book_text, book['languages'][0])

        ## create a new indexed book instance
        indexed_book_instance = {
            'book_id': book['id'],
            'created_at': current_time,
            'updated_at': current_time,
            'title': book['title'],
            'word_occurrence': json.dumps(dict(word_occurrence))
        }

        ## save the book instance to database
        insert_indexed_book_query = """INSERT INTO indexed_books (created_at, updated_at, title, word_occurrence_json, book_id)
                          VALUES (%s, %s, %s, %s, %s)"""
        cursor.execute(insert_indexed_book_query, (indexed_book_instance['created_at'], indexed_book_instance['updated_at'],
                                                   indexed_book_instance['title'], indexed_book_instance['word_occurrence'],
                                                   indexed_book_instance['book_id']))
        
        print("Book added...")
        sys.stdout.flush()

    conn.commit()


def main():
    ## connect to the database
    conn = connect_to_database()
    
    if conn:
        ## delete every book from the db
        delete_books_query = 'DELETE FROM books'
        delete_indexed_books_query = 'DELETE FROM indexed_books'
        delete_jaccard_neighbors_query = 'DELETE FROM jaccard_neighbors'

        try: 
            # cursor = conn.cursor()
            # cursor.execute(delete_books_query)
            # cursor.execute(delete_indexed_books_query)
            # cursor.execute(delete_jaccard_neighbors_query)
            # conn.commit()
            ...
        except Exception as e:
            print("An error occured while deleting.")

        ## fetch books from Gutenberg Project
        fetch_and_store_data(conn)
        print("All books successfully added to database.")

        compute_jaccard.compute_jaccard_distance(conn)
        
        ## close database connection
        conn.close()
        print("Database connection closed.")
    else:
        print("Hatme jeeb")

if __name__ == "__main__":
    main()
    # import time 
    # while True: time.sleep(99999)
