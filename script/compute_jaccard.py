import os
# import sqlite3
import requests
import psycopg2
from datetime import datetime
from utils import jaccard_distance
import json 
import time

TRESHOLD = 0.65


## connect to the SQLite database


## computer the jaccard distance
def compute_jaccard_distance(conn):
    cursor = conn.cursor()
            
    ## fetch data from db
    query = "SELECT book_id, word_occurrence_json FROM indexed_books;"  
    cursor.execute(query)
    print("Query executed")
    all_books = cursor.fetchall()
    size_books = len(all_books)

    for book1 in all_books:
        book1_id, book1_word_occurrence = book1
        books_neighbor = []
        sum_distance = 0

        for book2 in all_books:
            book2_id, book2_word_occurrence = book2
            if book1_id != book2_id:
                d1 = json.loads(book1_word_occurrence)
                d2 = json.loads(book2_word_occurrence)
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
                        VALUES (%s, %s, %s, %s)"""
        cursor.execute(insert_neighbors_query, (jaccard_neighbors_instance['created_at'], 
                                                jaccard_neighbors_instance['updated_at'],
                                                jaccard_neighbors_instance['neighbors'], 
                                                jaccard_neighbors_instance['book_id']))
        
        crank = (size_books - 1) / sum_distance
        cursor.execute(f"UPDATE books SET c_rank = %s WHERE book_id = %s", (crank, book1_id))
        conn.commit()
        print(f"Jaccard distance added for book with id: {book1_id}, sum distance: {sum_distance}, crank: {crank}")


    conn.commit()

