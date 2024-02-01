# DAAR Choice A
Group members:
    - Ufuk BOMBAR
    - Caterina Leonelli
    - Elif OHRI

# Project Introduction
This report delves into the development of an efficient search engine based on the Gutenberg books
database, a repository containing a substantial volume of diverse literary works. In the context of the
ever-expanding digital landscape, search engines play a pivotal role in helping users access pertinent
information amidst the vastness of the Internet. As evidenced by the exponential growth in indexed
pages, with Google’s index surpassing 8 billion pages in 2004, the need for effective search mechanisms
becomes paramount.
Search engines serve as filters, allowing users to navigate through the plethora of available in-
formation quickly and effortlessly. The challenge lies in delivering relevant results to users, and to
achieve this, search engines employ sophisticated algorithms to assess and rank web pages for specific
search expressions [2]. This paper explores the intricate process of developing a search engine that
not only utilizes advanced technologies but also employs algorithms to ensure precision and relevance
in delivering search results.
In the subsequent sections, we will delve into the methodologies, technologies, and algorithms
employed in the creation of this search engine, emphasizing the significance of relevance and efficiency
in handling extensive datasets such as the Gutenberg books database. By leveraging automated
programs, known as ”spiders” or ”robots,” to collect information [1], and implementing complex
algorithms that are frequently updated and closely guarded [5], our objective is to present a search
engine capable of providing users with tailored and pertinent results for a diverse array of search
queries.

# Requirements
To run the whole system it is required to have _docker_ installed. 

For non-dockerized tests and runs:
- Backend: for tests you need the latest version of go (go1.21.3). 
- Frontend: you need the latest version of _node_ and _npm_ installed. 
- Script: The script cannot be run outside container, we designed it to be run as a container.

# How to Setup
The project uses _docker_ to run it's components in containers. The files to _build_, _start, and _stop_ the system can be found in the make file. You can use the following command to run the whole project. (Note that the preprocessing of the book data takes ~2 hours since the full text also downloaded. This is why we also provide a backup version of the postgres database. However we haven't included in the zip since it is very large ~100Mb. If requested this can be provided.)

To build the project
```bash
make build
```

To start the project
```bash
make start
```

Normally when CTRL+C is pressed the containers are shut down gracefully. But if for some reason they continue to run on the background, they can be destroyed by the following command.
```bash
make stop
```

## Non-Dockerized Backend
To run the backend unit tests:
```bash
cd ./backend # to run the tests you need to be on the backend folder
go test ./.. # run all of the unit tests
```

## Non-Dockerized Frontend
```bash
cd ./frontend
npm i
npm start
```

# Submission
The zip file should contain the README.md, report.pdf, video.mp4.

# Links
Github project: https://github.com/RootLeo00/book-web-app-daar
Youtube Video Pitch: https://youtu.be/NCjLZ_Y1YkQ

# Conclusion
The execution of this project has deepened our understanding of efficient result classification and,
more importantly, the significant impact that various classifications can have on a search. While it
might be more cost-effective and simpler for a search engine site to randomly display pages or results
based on word count, freshness, or other straightforward sorting systems, they refrain from doing so.
The primary reason for this choice is user-centric ˆa users would not find it useful. The users we aim to
satisfy will not be drawn to this search engine unless it provides precise results and optimal assistance.