import os
import sqlite3
import requests
import logging
from datetime import datetime
from utils import getWordsWithOcc
import json 

DB_FILE_PATH = "./db.sqlite3"
GUTENDEX_URL = "https://gutendex.com/books"
MIME_TYPE = "text"
MAX_PAGES = 10

## connect to the SQLite database
def connect_to_database(DB_FILE_PATH):
    try:
        conn = sqlite3.connect(DB_FILE_PATH)
        print("Successfully connected to the database.")
        return conn    
    except sqlite3.Error as e:
        print(f"Database error: {e}")
        return None
    

## construct the URL for the current page
def construct_url(api_endpoint, mime_type, page):
    return f"{api_endpoint}?mime_type={mime_type}&page={page}"


## fetch all book data from a page
def fetch_and_store_data(conn):
    page = 1
    while page < MAX_PAGES:
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
        url_text = book['formats'].get('text/plain; charset=us-ascii', None)
        if url_text:
            url_text = url_text.replace('.zip', '.txt')
        author = 'None'
        if book['authors']:
            author = book['authors'][0]['name']
        text_response = book['title'] + " " + " ".join(book['subjects'])
        current_time = datetime.now().strftime("%Y-%m-%d %H:%M:%S.%f")

        ## create a new book instance
        book_instance = {
            'book_id': book['id'],
            'created_at': current_time,
            'updated_at': current_time,
            'title': book['title'],
            'author': author,
            'language': book['languages'][0],
            'text': text_response,      # this is not the actual text
            'image_url': book['formats'].get('image/jpeg', None),
            'c_rank': 0.0,
            'occurrence': 0
        }

        ## save the book instance to database
        insert_book_query = """INSERT INTO books (created_at, updated_at, book_id, title, author, language, text, image_url, c_rank, occurrence)
                          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"""
        cursor.execute(insert_book_query, (book_instance['created_at'], book_instance['updated_at'], book_instance['book_id'],
                                        book_instance['title'], book_instance['author'], book_instance['language'], 
                                        book_instance['text'], book_instance['image_url'], book_instance['c_rank'], 
                                        book_instance['occurrence']))
        
        word_occurence = getWordsWithOcc(text_response, book['languages'][0])

        ## create a new indexed book instance
        indexed_book_instance = {
            'book_id': book['id'],
            'created_at': current_time,
            'updated_at': current_time,
            'title': book['title'],
            'word_occurrence': json.dumps(dict(word_occurence))
        }

        ## save the book instance to database
        insert_indexed_book_query = """INSERT INTO indexed_books (created_at, updated_at, title, word_occurrence_json, book_id)
                          VALUES (?, ?, ?, ?, ?)"""
        cursor.execute(insert_indexed_book_query, (indexed_book_instance['created_at'], indexed_book_instance['updated_at'],
                                                   indexed_book_instance['title'], indexed_book_instance['word_occurrence'],
                                                   indexed_book_instance['book_id']))

    conn.commit()


def main():

    abs_db_path = os.path.abspath(DB_FILE_PATH)

    if not os.path.exists(abs_db_path):
        print(f"Database file not found at {abs_db_path}")
    else:
        ## connect to the database
        conn = connect_to_database(DB_FILE_PATH)
        
        if conn:
            ## delete every book from the db
            delete_books_query = 'DELETE FROM books'
            delete_indexed_books_query = 'DELETE FROM indexed_books'

            try: 
                cursor = conn.cursor()
                cursor.execute(delete_books_query)
                cursor.execute(delete_indexed_books_query)
                conn.commit()
            except sqlite3.Error as e:
                print("An error occured while deleting.")

            ## fetch books from Gutenberg Project
            fetch_and_store_data(conn)
            print("All books successfully added to database.")
            
            ## close database connection
            conn.close()
            print("Database connection closed.")

if __name__ == "__main__":
    main()
