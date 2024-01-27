class BookService {
  urlApi;

  constructor() {
    this.urlApi = process.env["BACKEND_URL"] || "http://127.0.0.1:8080"
    // this.urlApi = "http://192.168.188.134:8080"
  }

  searchBook(word:string) {
    return fetch(`${this.urlApi}/Search/${word}`)
      .then(response => response.json())
      .catch(error => {
        console.error('Error fetching data:', error);
        throw error;
      });
  }

  searchBookRegex(regex:string) {
    return fetch(`${this.urlApi}/RegexSearch/${regex}`)
      .then(response => response.json())
      .catch(error => {
        console.error('Error fetching data:', error);
        throw error;
      });
  }
}

export default BookService;
