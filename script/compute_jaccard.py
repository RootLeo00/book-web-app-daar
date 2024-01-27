import os
import sqlite3
import requests
import psycopg2
from datetime import datetime
from utils import jaccard_distance
import json 

DB_FILE_PATH = "./db.sqlite3"
TRESHOLD = 0.65

POSTGRES_DB = os.environ.get("POSTGRES_DB")
POSTGRES_USER = os.environ.get("POSTGRES_USER")
POSTGRES_PASSWORD = os.environ.get("POSTGRES_PASSWORD")

## connect to the SQLite database
def connect_to_database(DB_FILE_PATH):
    ## try to connect to PostgreSQL first
    try:
        conn = psycopg2.connect(dbname=POSTGRES_DB,
                                user=POSTGRES_USER,
                                password=POSTGRES_PASSWORD)
        print("Successfully connected to the PostgreSQL database.")
        return conn
    except psycopg2.Error as e:
        print(f"PostgreSQL Database error: {e}")
        print("Trying to connect to SQLite database...")

        ## Fallback to SQLite3
        try:
            conn = sqlite3.connect(DB_FILE_PATH)
            print("Successfully connected to the SQLite database.")
            return conn
        except sqlite3.Error as e:
            print(f"SQLite Database error: {e}")
            return None


## computer the jaccard distance
def compute_jaccard_distance(conn):
    cursor = conn.cursor()
    sum_distance = 0
            
    ## fetch data from db
    query = "SELECT * FROM indexed_books;"  
    cursor.execute(query)
    print("Query executed")
    all_books = cursor.fetchall()
    size_books = len(all_books)

    for book1 in all_books:
        _, _, _, _, book1_id, _, book1_word_occurence = book1
        books_neighbor = []
        sum_distance = 0

        for book2 in all_books:
            _, _, _, _, book2_id, _, book2_word_occurence = book2
            if book1_id != book2_id:
                d1 = json.loads(book1_word_occurence)
                d2 = json.loads(book2_word_occurence)
                result_distance = jaccard_distance(d1, d2)
                
                if result_distance < TRESHOLD:
                    books_neighbor.append(book2_id)
                
                sum_distance += result_distance

        ## save the neighbors and crank value to the database
        current_time = datetime.now().strftime("%Y-%m-%d %H:%M:%S.%f")
        jaccard_neighbors_instance = {
            'book_id': book1_id,
            'created_at': current_time,
            'updated_at': current_time,
            "neighbors": json.dumps(books_neighbor)
        }

        ## save jaccard neighbors to database
        insert_neighbors_query = """INSERT INTO jaccard_neighbors (created_at, updated_at, neighbors_json, book_id)
                        VALUES (?, ?, ?, ?)"""
        cursor.execute(insert_neighbors_query, (jaccard_neighbors_instance['created_at'], jaccard_neighbors_instance['updated_at'],
                                            jaccard_neighbors_instance['neighbors'], jaccard_neighbors_instance['book_id']))
        
        crank = (size_books - 1) / sum_distance
        cursor.execute(f"UPDATE books SET c_rank = ? WHERE book_id = ?", (crank, book1_id))
        conn.commit()
        print(f"Jaccard distance added for book with id: {book1_id}, sum distance: {sum_distance}, crank: {crank}")


    conn.commit()


def main():

    abs_db_path = os.path.abspath(DB_FILE_PATH)

    if not os.path.exists(abs_db_path):
        print(f"Database file not found at {abs_db_path}")
    else:
        ## Connect to the database
        conn = connect_to_database(DB_FILE_PATH)
        
        if conn:
            delete_jaccard_distance_query = 'DELETE FROM jaccard_neighbors'

            try: 
                ## delete every book from db
                cursor = conn.cursor()
                cursor.execute(delete_jaccard_distance_query)
                conn.commit()
            except sqlite3.Error as e:
                print("An error occured while deleting.")

            ## fetch books from Gutenberg Project
            compute_jaccard_distance(conn)
            
            ## close database connection
            conn.close()
            print("Database connection closed.")

if __name__ == "__main__":
    main()
